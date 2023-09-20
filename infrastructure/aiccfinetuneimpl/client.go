package aiccfinetuneimpl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/opensourceways/xihe-aicc-finetune/config"
	"github.com/opensourceways/xihe-aicc-finetune/infrastructure/aicc"

	"github.com/opensourceways/community-robot-lib/utils"
)

func newClient(cfg *config.AICCConfig) (aiccClient, error) {
	return aiccClient{
		Domain:       cfg.Domain,
		User:         cfg.User,
		Password:     cfg.Password,
		Project:      cfg.Project,
		AuthEndpoint: cfg.AuthEndpoint,
		Endpoint:     cfg.Endpoint,
	}, nil
}

func (s *aiccClient) token() (string, error) {
	str := `
{
    "auth":{
       "identity":{
          "methods":[
             "password"
          ],
          "password":{
             "user":{
                "name":"%v",
                "password":"%v",
                "domain":{
                   "name":"%v"
                }
             }
          }
       },
       "scope":{
          "project":{
             "name":"%s"
          }
       }
    }
}
	`

	body := fmt.Sprintf(
		str, s.User, s.Password, s.Domain, s.Project,
	)
	resp, err := http.Post(
		s.AuthEndpoint, "application/json",
		strings.NewReader(body),
	)
	if err != nil {
		return "", err
	}

	t := resp.Header.Get("x-subject-token")
	if err = resp.Body.Close(); err != nil {
		return "", err
	}

	return t, nil
}

type aiccClient struct {
	Domain       string
	User         string
	Password     string
	Project      string
	AuthEndpoint string
	Endpoint     string
}

func (cli *aiccClient) createURL() string {
	return cli.Endpoint + "/v2/e0412da2cb3b4ebfb70c117343b8992a/training-jobs"
}

func (cli *aiccClient) jobURL(jobId string) string {
	return cli.createURL() + "/" + jobId
}

func (cli *aiccClient) terminateURL(jobId string) string {
	return cli.createURL() + "/" + jobId + "/actions"
}

func (cli *aiccClient) logURL(jobId string) string {
	return cli.jobURL(jobId) + "/tasks/worker-0/logs/url"
}

func (cli *aiccClient) createJob(options aicc.JobCreateOption) (jobId string, err error) {
	payload, err := utils.JsonMarshal(options)
	if err != nil {
		return
	}
	req, err := http.NewRequest(
		http.MethodPost, cli.createURL(), bytes.NewBuffer(payload),
	)
	if err != nil {
		return
	}

	token, err := cli.token()
	if err != nil {
		return
	}

	resp, err := cli.forwardTo(req, token)
	if err != nil {
		return
	}

	if resp.StatusCode != 201 {
		return "", err
	}

	cr := new(aicc.Job)
	err = ParseResponse(resp, cr)
	jobId = cr.Metadata.Id

	return
}

func ParseResponse(resp *http.Response, result interface{}) (err error) {
	body, err := ioutil.ReadAll(resp.Body)
	if err == nil {
		err = json.Unmarshal(body, &result)
	}
	return
}

func (cli *aiccClient) forwardTo(req *http.Request, token string) (resp *http.Response, err error) {
	if token != "" {
		req.Header.Set("X-Auth-Token", token)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err = http.DefaultClient.Do(req)

	return
}

func (cli *aiccClient) getJob(jobId string) (info aicc.Job, err error) {
	req, err := http.NewRequest(http.MethodGet, cli.jobURL(jobId), nil)
	if err != nil {
		return
	}

	token, err := cli.token()
	if err != nil {
		return
	}

	resp, err := cli.forwardTo(req, token)
	if err != nil {
		return
	}

	if resp.StatusCode == 200 {
		res := new(aicc.Job)
		err = ParseResponse(resp, res)
		info = *res
	} else {
		err = errors.New(resp.Status)
	}

	return
}

func (cli *aiccClient) deleteJob(jobId string) (err error) {
	req, err := http.NewRequest(http.MethodDelete, cli.jobURL(jobId), nil)
	if err != nil {
		return
	}

	token, err := cli.token()
	if err != nil {
		return
	}

	resp, err := cli.forwardTo(req, token)
	if err != nil {
		return
	}

	if resp.StatusCode != 202 {
		err = errors.New(resp.Status)
		return
	}

	return
}

func (cli *aiccClient) terminateJob(jobId string) (err error) {
	options := new(aicc.TerminateBody)
	options.ActionType = "terminate"

	payload, err := utils.JsonMarshal(options)
	if err != nil {
		return
	}

	req, err := http.NewRequest(http.MethodPost, cli.terminateURL(jobId), bytes.NewBuffer(payload))
	if err != nil {
		return
	}

	token, err := cli.token()
	if err != nil {
		return
	}

	resp, err := cli.forwardTo(req, token)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = errors.New(resp.Status)
	}

	return
}

func (cli *aiccClient) getLogURL(jobId string) (log string, err error) {
	req, err := http.NewRequest(http.MethodGet, cli.logURL(jobId), nil)
	if err != nil {
		return
	}

	token, err := cli.token()
	if err != nil {
		return
	}

	resp, err := cli.forwardTo(req, token)
	if err != nil {
		return
	}
	res := new(aicc.LogResp)
	err = ParseResponse(resp, &res)
	if resp.StatusCode == 200 {
		log = res.ObsUrl
	} else {
		err = errors.New(resp.Status)
	}

	return
}
