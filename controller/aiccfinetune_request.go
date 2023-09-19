package controller

import (
	"errors"

	"github.com/opensourceways/xihe-aicc-finetune/app"
	"github.com/opensourceways/xihe-aicc-finetune/domain"
)

type AICCFinetuneResultResp struct {
	URL string `json:"url"`
}

type AICCFinetuneCreateRequest struct {
	User       string `json:"user"`
	Model      string `json:"model"`
	FinetuneId string `json:"finetune_id"`
	Task       string `json:"task"`

	Name string `json:"name"`
	Desc string `json:"desc"`

	Hyperparameters []AICCKeyValue `json:"hyperparameter"`
	Env             []AICCKeyValue `json:"env"`
}

func (req *AICCFinetuneCreateRequest) toCmd(cmd *app.AICCFinetuneCreateCmd) (err error) {
	if cmd.User, err = domain.NewAccount(req.User); err != nil {
		return
	}

	if cmd.Model, err = domain.NewModelName(req.Model); err != nil {
		return
	}

	if cmd.Name, err = domain.NewFinetuneName(req.Name); err != nil {
		return
	}

	if cmd.Desc, err = domain.NewFinetuneDesc(req.Desc); err != nil {
		return
	}

	if cmd.Env, err = req.toKeyValue(req.Env); err != nil {
		return
	}

	if cmd.Hyperparameters, err = req.toKeyValue(req.Hyperparameters); err != nil {
		return
	}

	cmd.Task = req.Task

	return
}

func (req *AICCFinetuneCreateRequest) toKeyValue(kv []AICCKeyValue) (r []domain.KeyValue, err error) {
	n := len(kv)
	if n == 0 {
		return nil, nil
	}

	r = make([]domain.KeyValue, n)
	for i := range kv {
		if r[i], err = kv[i].toKeyValue(); err != nil {
			return
		}
	}

	return
}

type AICCKeyValue struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

func (kv *AICCKeyValue) toKeyValue() (r domain.KeyValue, err error) {
	if kv.Key == "" {
		err = errors.New("invalid key value")

		return
	}

	if r.Key, err = domain.NewCustomizedKey(kv.Key); err != nil {
		return
	}

	r.Value, err = domain.NewCustomizedValue(kv.Value)

	return
}
