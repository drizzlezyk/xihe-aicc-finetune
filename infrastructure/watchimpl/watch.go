package watchimpl

import (
	"errors"
	"fmt"
	"sync"
	"time"

	pt "github.com/opensourceways/xihe-grpc-protocol/grpc/aiccfinetune"

	"github.com/opensourceways/xihe-grpc-protocol/grpc/client"
	"github.com/sirupsen/logrus"

	"github.com/opensourceways/xihe-aicc-finetune/domain/aiccfinetune"
	"github.com/opensourceways/xihe-aicc-finetune/domain/watch"
)

type aiccFinetuneData = pt.AICCFinetuneInfo

func NewWatcher(
	cfg *Config, as aiccfinetune.AICCFinetune,
) (*Watcher, error) {
	cli, err := client.NewAICCFinetuneClient(cfg.Endpoint)
	if err != nil {
		return nil, err
	}

	return &Watcher{
		cli:         cli,
		as:          as,
		timeout:     cfg.Timeout,
		interval:    time.Duration(cfg.Interval) * time.Second,
		stop:        make(chan struct{}),
		stopped:     make(chan struct{}),
		finetunes:   make(chan finetuneInfo, cfg.MaxWatchNum+1),
		maxWatchNum: cfg.MaxWatchNum,
	}, nil
}

type finetuneInfo struct {
	watch.FinetuneInfo

	result aiccFinetuneData

	done       bool
	success    bool
	logDone    bool
	outputDone bool
}

func (t *finetuneInfo) toIndex() pt.AICCFinetuneIndex {
	return pt.AICCFinetuneIndex{
		Id:    t.FinetuneId,
		User:  t.User.Account(),
		Model: "wukong",
	}
}

func (t *finetuneInfo) isDone() bool {
	done := t.done && t.logDone

	if done && t.success {
		done = t.outputDone
	}

	return done
}

// Watcher
type Watcher struct {
	log *logrus.Entry
	cli *client.AICCFinetuneClient
	as  aiccfinetune.AICCFinetune

	timeout  int
	interval time.Duration

	stop      chan struct{}
	stopped   chan struct{}
	finetunes chan finetuneInfo

	lock        sync.RWMutex
	currentNum  int
	maxWatchNum int
}

func (w *Watcher) ApplyWatch(f func(*watch.FinetuneInfo) error) (err error) {
	if !w.increase() {
		return errors.New("exceed max watch num")
	}

	info := new(watch.FinetuneInfo)

	if err = f(info); err != nil {
		w.decrease()
	} else {
		w.addFinetune(info)
	}

	return
}

func (w *Watcher) addFinetune(t *watch.FinetuneInfo) {
	info := finetuneInfo{FinetuneInfo: *t}
	w.finetunes <- info
}

func (w *Watcher) increase() (b bool) {
	w.lock.Lock()
	if w.currentNum+1 <= w.maxWatchNum {
		w.currentNum++
		b = true
	}
	w.lock.Unlock()

	return
}

func (w *Watcher) decrease() {
	w.lock.Lock()
	w.currentNum--
	w.lock.Unlock()
}

func (w *Watcher) Run() {
	start := time.Now()

	// add the tag
	w.finetunes <- finetuneInfo{}

	for {
		select {
		case info := <-w.finetunes:
			// use =="" stands for the case that the loop is done
			if info.User == nil {
				w.log.Debug("finish a loop")

				t := start.Add(w.interval)

				if n := time.Now(); t.After(n) {
					time.Sleep(t.Sub(n))
				}

				w.finetunes <- finetuneInfo{}

				start = time.Now()

			} else {
				changed := w.check(&info)
				fmt.Printf("aicc info: %+v\n", info)
				w.log.Debugf("check aicc finetune %s/%s", info.FinetuneId, info.JobId)
				if info.isDone() {
					index := info.toIndex()

					if err := w.cli.SetAICCFinetuneInfo(&index, &info.result); err == nil {
						w.decrease()
					} else {
						w.log.Errorf("set aicc finetune info failed, err:%s", err.Error())
						w.finetunes <- info
					}

				} else {
					if changed {
						fmt.Printf("aicc info: %+v\n", info)
						index := info.toIndex()
						if err := w.cli.SetAICCFinetuneInfo(&index, &info.result); err != nil {
							w.log.Errorf("set aicc finetune info failed, err:%s", err.Error())
						}
					}

					w.finetunes <- info
				}
			}

		case <-w.stop:
			close(w.stopped)

			return
		}
	}
}

func (w *Watcher) Exit() {
	close(w.stop)

	<-w.stopped

	w.cli.Disconnect()
}

func (w *Watcher) check(info *finetuneInfo) (changed bool) {
	result := &info.result

	if !info.done {
		detail, err := w.as.GetDetail(info.JobId)
		if err != nil {
			return
		}
		if detail.Duration != result.Duration {
			result.Duration = detail.Duration
			changed = true
		}

		if s := detail.Status.TrainingStatus(); s != result.Status {
			result.Status = s
			changed = true
		}

		if !detail.Status.IsDone() {
			if detail.Duration < w.timeout {
				return
			}

			if err := w.as.Terminate(info.JobId); err != nil {
				w.log.Errorf(
					"terminate the job(%s) failed, err:%s",
					info.JobId, err.Error(),
				)

				return
			}

			result.Status = "Timeout"
			changed = true
		} else {
			info.success = detail.Status.IsSuccess()
		}
		info.done = true
	}

	if !info.logDone {
		if v, err := w.as.GetLogFilePath(info.LogDir); err != nil {
			w.log.Errorf("generate log failed, err:%s", err.Error())
		} else {
			result.LogPath = v
			info.logDone = true
			changed = true
		}
	}

	if !info.success {
		return
	}

	if !info.outputDone {
		if v, err := w.as.GenOutput(info.OutputDir); err != nil {
			w.log.Errorf("generate output failed, err:%s", err.Error())
		} else {
			info.outputDone = true

			if v != "" {
				result.OutputZipPath = v
				changed = true
			}
		}
	}

	return
}
