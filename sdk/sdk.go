package sdk

import (
	"bytes"
	"fmt"
	"net/http"
	"strings"

	"github.com/opensourceways/community-robot-lib/utils"
	"github.com/opensourceways/xihe-aicc-finetune/app"
	"github.com/opensourceways/xihe-aicc-finetune/controller"
)

type AICCFinetuneCreateOption = controller.AICCFinetuneCreateRequest
type DownloadURL = controller.AICCFinetuneResultResp
type KeyValue = controller.AICCKeyValue
type JobInfo = app.JobInfoDTO

func NewAICCFinetuneCenter(endpoint string) AICCFinetuneCenter {
	s := strings.TrimSuffix(endpoint, "/")
	if p := "/api/v1/aiccfinetune"; !strings.HasSuffix(s, p) {
		s += p
	}

	return AICCFinetuneCenter{
		endpoint: s,
		cli:      utils.NewHttpClient(3),
	}
}

type AICCFinetuneCenter struct {
	endpoint string
	cli      utils.HttpClient
}

func (t AICCFinetuneCenter) jobURL(jobId string) string {
	return fmt.Sprintf("%s/%s", t.endpoint, jobId)
}

func (t AICCFinetuneCenter) CreateAICCFinetune(opt *AICCFinetuneCreateOption) (
	dto JobInfo, err error,
) {
	payload, err := utils.JsonMarshal(&opt)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, t.endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	v := new(app.JobInfoDTO)
	if err = t.forwardTo(req, v); err != nil {
		return
	}

	return *v, nil
}

func (t AICCFinetuneCenter) DeleteAICCFinetune(jobId string) error {
	req, err := http.NewRequest(http.MethodDelete, t.jobURL(jobId), nil)
	if err != nil {
		return err
	}

	return t.forwardTo(req, nil)
}

func (t AICCFinetuneCenter) TerminateeAICCFinetune(jobId string) error {
	req, err := http.NewRequest(http.MethodPut, t.jobURL(jobId), nil)
	if err != nil {
		return err
	}

	return t.forwardTo(req, nil)
}

func (t AICCFinetuneCenter) GetLogDownloadURL(jobId string) (r DownloadURL, err error) {
	req, err := http.NewRequest(http.MethodGet, t.jobURL(jobId)+"/log", nil)
	if err != nil {
		return
	}

	if err = t.forwardTo(req, &r); err != nil {
		return
	}

	return
}

func (t AICCFinetuneCenter) GetResultDownloadURL(jobId, file string) (r DownloadURL, err error) {
	req, err := http.NewRequest(http.MethodGet, t.jobURL(jobId)+"/result/"+file, nil)
	if err != nil {
		return
	}

	if err = t.forwardTo(req, &r); err != nil {
		return
	}

	return
}

func (t AICCFinetuneCenter) forwardTo(req *http.Request, jsonResp interface{}) (err error) {
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", "xihe-aicc-finetune")

	if jsonResp != nil {
		v := struct {
			Data interface{} `json:"data"`
		}{jsonResp}

		_, err = t.cli.ForwardTo(req, &v)
	} else {
		_, err = t.cli.ForwardTo(req, jsonResp)
	}

	return
}
