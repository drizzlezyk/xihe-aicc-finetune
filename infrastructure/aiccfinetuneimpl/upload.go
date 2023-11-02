package aiccfinetuneimpl

import (
	"fmt"
	"strings"

	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
	"github.com/opensourceways/xihe-aicc-finetune/config"
)

func newHelper(cfg *config.Config) (*helper, error) {
	obsCfg := &cfg.OBS
	cli, err := obs.New(obsCfg.AccessKey, obsCfg.SecretKey, obsCfg.Endpoint)
	if err != nil {
		return nil, fmt.Errorf("new obs client failed, err:%s", err.Error())
	}

	suc := &cfg.Upload

	return &helper{
		obsClient: cli,
		bucket:    obsCfg.Bucket,
		suc:       *suc,
	}, nil
}

type helper struct {
	obsClient *obs.ObsClient
	bucket    string
	suc       config.UploadConfig
}

func (s *helper) GetLogFilePath(logDir string) (p string, err error) {
	if !strings.HasSuffix(logDir, "/") {
		logDir += "/"
	}

	input := &obs.ListObjectsInput{}
	input.Bucket = s.bucket
	input.Prefix = logDir // "src0/"

	output, err := s.obsClient.ListObjects(input)
	if err != nil {
		return
	}

	v := output.Contents
	for i := range v {
		if p = v[i].Key; p != logDir {
			break
		}
	}

	return
}

func (s *helper) GenFileDownloadURL(p string) (string, error) {
	input := &obs.CreateSignedUrlInput{}
	input.Method = obs.HttpMethodGet
	input.Bucket = s.bucket
	input.Key = p
	input.Expires = s.suc.DownloadExpiry

	output, err := s.obsClient.CreateSignedUrl(input)
	if err != nil {
		return "", err
	}

	return output.SignedUrl, nil
}
