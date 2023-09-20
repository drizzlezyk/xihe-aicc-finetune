package app

import (
	"errors"

	"github.com/opensourceways/xihe-aicc-finetune/domain"
	"github.com/opensourceways/xihe-aicc-finetune/domain/aiccfinetune"
	"github.com/opensourceways/xihe-aicc-finetune/domain/watch"
	"github.com/sirupsen/logrus"
)

type AICCFinetuneCreateCmd struct {
	FinetuneId string

	domain.AICCFinetune
}

func (cmd *AICCFinetuneCreateCmd) Validate() error {
	err := errors.New("invalid cmd of creating aicc finetune")

	b := cmd.User != nil &&
		cmd.Name != nil &&
		cmd.FinetuneId != ""

	if !b {
		return err
	}

	f := func(kv []domain.KeyValue) error {
		for i := range kv {
			if kv[i].Key == nil {
				return err
			}
		}

		return nil
	}

	if f(cmd.Hyperparameters) != nil {
		return err
	}

	if f(cmd.Env) != nil {
		return err
	}

	return nil
}

type JobInfoDTO = domain.JobInfo

type FinetuneService interface {
	Create(cmd *AICCFinetuneCreateCmd) (JobInfoDTO, error)
	Delete(jobId string) error
	Terminate(jobId string) error
	GetLogDownloadURL(jobId string) (string, error)
	GenFileDownloadURL(obsfile string) (string, error)
}

func NewAICCFinetuneService(
	ts aiccfinetune.AICCFinetune,
	ws watch.WatchService,
	log *logrus.Entry,
) FinetuneService {
	return &aiccFinetuneService{
		ts:  ts,
		ws:  ws,
		log: log,
	}
}

type aiccFinetuneService struct {
	log *logrus.Entry
	ts  aiccfinetune.AICCFinetune
	ws  watch.WatchService
}

func (s *aiccFinetuneService) Create(cmd *AICCFinetuneCreateCmd) (JobInfoDTO, error) {
	dto := JobInfoDTO{}

	f := func(info *watch.FinetuneInfo) error {
		v, err := s.create(cmd)
		if err != nil {
			return err
		}

		dto = v
		*info = watch.FinetuneInfo{
			User:       cmd.User,
			FinetuneId: cmd.FinetuneId,
			JobInfo:    v,
		}

		return nil
	}

	err := s.ws.ApplyWatch(f)

	return dto, err
}

func (s *aiccFinetuneService) create(cmd *AICCFinetuneCreateCmd) (info domain.JobInfo, err error) {
	if cmd.Task == "finetune" {
		return s.ts.Create(&cmd.AICCFinetune)
	}
	return s.ts.CreateInference(&cmd.AICCFinetune)
}

func (s *aiccFinetuneService) Terminate(jobId string) error {
	return s.ts.Terminate(jobId)
}

func (s *aiccFinetuneService) GetLogDownloadURL(jobId string) (string, error) {
	return s.ts.GetLogDownloadURL(jobId)
}

func (s *aiccFinetuneService) GenFileDownloadURL(obsfile string) (string, error) {
	return s.ts.GenFileDownloadURL(obsfile)
}

func (s *aiccFinetuneService) Delete(jobId string) error {
	return s.ts.Delete(jobId)
}
