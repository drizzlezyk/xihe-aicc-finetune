package watch

import "github.com/opensourceways/xihe-aicc-finetune/domain"

type FinetuneInfo struct {
	User       domain.Account
	FinetuneId string
	Model      string

	domain.JobInfo
}

type WatchService interface {
	ApplyWatch(f func(*FinetuneInfo) error) (err error)
}
