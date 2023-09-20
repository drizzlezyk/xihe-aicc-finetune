package aiccfinetuneimpl

import (
	"strconv"
	"strings"
	"time"

	"github.com/opensourceways/xihe-aicc-finetune/config"
	"github.com/opensourceways/xihe-aicc-finetune/domain"
	"github.com/opensourceways/xihe-aicc-finetune/domain/aiccfinetune"

	"github.com/opensourceways/xihe-aicc-finetune/infrastructure/aicc"
)

const (
	obsDelimiter = "/"
	modelPathKey = "model_path"
	dataPathKey  = "finetune_data_path"
)

var statusMap = map[string]domain.TrainingStatus{
	"failed":      domain.TrainingStatusFailed,
	"pending":     domain.TrainingStatusPending,
	"running":     domain.TrainingStatusRunning,
	"creating":    domain.TrainingStatusCreating,
	"abnormal":    domain.TrainingStatusAbnormal,
	"completed":   domain.TrainingStatusCompleted,
	"terminated":  domain.TrainingStatusTerminated,
	"terminating": domain.TrainingStatusTerminating,
}

func NewAiccFinetune(cfg *config.Config) (aiccfinetune.AICCFinetune, error) {
	cli, err := newClient(&cfg.AICC)

	if err != nil {
		return nil, err
	}

	h, err := newHelper(cfg)
	if err != nil {
		return nil, err
	}

	return aiccFinetuneImpl{
		cli:    cli,
		config: cfg.Finetune,
		helper: h,
	}, nil
}

type aiccFinetuneImpl struct {
	cli    aiccClient
	config config.FinetuneConfig

	*helper
}

func (impl aiccFinetuneImpl) genJobParameter(t *domain.AICCFinetune, opt *aicc.JobCreateOption) {
	if n := len(t.Hyperparameters); n > 0 {
		p := make([]aicc.ParameterOption, n)

		for i, v := range t.Hyperparameters {
			s := ""
			if v.Value != nil {
				s = v.Value.CustomizedValue()
			}

			p[i] = aicc.ParameterOption{
				Name:  v.Key.CustomizedKey(),
				Value: s,
			}
		}

		opt.Algorithm.Parameters = p
	}

	if n := len(t.Env); n > 0 {
		m := make(map[string]string)

		for _, v := range t.Env {
			s := ""
			if v.Value != nil {
				s = v.Value.CustomizedValue()
			}

			m[v.Key.CustomizedKey()] = s
		}

		opt.Algorithm.Environments = m
	}
}

func (impl aiccFinetuneImpl) Create(t *domain.AICCFinetune) (info domain.JobInfo, err error) {
	var cfg *config.ModelConfig
	if t.Model.ModelName() == "wukong" {
		cfg = &impl.config.WukongConfig
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	logDir := cfg.LogDir + t.Task + obsDelimiter + t.User.Account() + obsDelimiter
	outputDir := cfg.OutputDir + t.Task + obsDelimiter + t.User.Account() + obsDelimiter

	outputs := []aicc.InputOutputOption{}
	outputs = append(outputs, aicc.InputOutputOption{
		Name: cfg.OutputKey,
		Remote: aicc.RemoteOption{
			OBS: aicc.OBSOption{
				OBSURL: outputDir,
			},
		},
	})

	inputs := []aicc.InputOutputOption{}
	inputs = append(inputs, aicc.InputOutputOption{
		Name: modelPathKey,
		Remote: aicc.RemoteOption{
			OBS: aicc.OBSOption{
				OBSURL: cfg.ModelDir,
			},
		},
	})
	inputs = append(inputs, aicc.InputOutputOption{
		Name: dataPathKey,
		Remote: aicc.RemoteOption{
			OBS: aicc.OBSOption{
				OBSURL: cfg.InputDir + t.Task + obsDelimiter + t.User.Account() + obsDelimiter,
			},
		},
	})

	opt := aicc.JobCreateOption{
		Kind: "job",
		Metadata: aicc.MetadataOption{
			Name: t.Name.FinetuneName() + t.User.Account() + "-" + timestamp + "-" + t.Task,
			Desc: t.Desc.FinetuneDesc(),
		},
		Algorithm: aicc.AlgorithmOption{
			CodeDir:    cfg.CodeDir,
			WorkingDir: cfg.WorkingDir,
			Command:    t.DefaultCommand(),
			Engine: aicc.EngineOption{
				ImageURL: cfg.ImageURL,
			},
			Inputs:  inputs,
			Outputs: outputs,
		},
		Spec: aicc.SpecOption{
			Resource: aicc.ResourceOption{
				FlavorId:  cfg.FlavorId,
				PoolId:    cfg.PoolId,
				PoolName:  cfg.PoolName,
				NodeCount: 1,
			},
			LogExportPath: aicc.LogExportPathOption{
				OBSURL: logDir,
			},
		},
	}

	impl.genJobParameter(t, &opt)

	info.JobId, err = impl.cli.createJob(opt)

	if err == nil {
		p := ""
		info.LogDir = strings.TrimPrefix(logDir, p)
		info.OutputDir = strings.TrimPrefix(outputDir, p)
	}

	return
}

func (impl aiccFinetuneImpl) CreateInference(t *domain.AICCFinetune) (info domain.JobInfo, err error) {
	var cfg *config.ModelConfig
	if t.Model.ModelName() == "wukong" {
		cfg = &impl.config.WukongConfig
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	logDir := cfg.LogDir + t.Task + t.User.Account() + obsDelimiter
	outputDir := cfg.OutputDir + t.Task + t.User.Account() + obsDelimiter

	outputs := []aicc.InputOutputOption{}
	outputs = append(outputs, aicc.InputOutputOption{
		Name: cfg.OutputKey,
		Remote: aicc.RemoteOption{
			OBS: aicc.OBSOption{
				OBSURL: outputDir,
			},
		},
	})

	inputs := []aicc.InputOutputOption{}
	inputs = append(inputs, aicc.InputOutputOption{
		Name: modelPathKey,
		Remote: aicc.RemoteOption{
			OBS: aicc.OBSOption{
				OBSURL: cfg.ModelDir,
			},
		},
	})

	inputs = append(inputs, aicc.InputOutputOption{
		Name: dataPathKey,
		Remote: aicc.RemoteOption{
			OBS: aicc.OBSOption{
				OBSURL: cfg.InputDir + t.Task + t.User.Account() + obsDelimiter,
			},
		},
	})

	opt := aicc.JobCreateOption{
		Kind: "job",
		Metadata: aicc.MetadataOption{
			Name: t.Name.FinetuneName() + t.User.Account() + "-" + timestamp + "-" + t.Task,
			Desc: t.Desc.FinetuneDesc(),
		},
		Algorithm: aicc.AlgorithmOption{
			CodeDir:    cfg.CodeDir,
			WorkingDir: cfg.WorkingDir,
			Command:    t.DefaultInferenceCommand(),
			Engine: aicc.EngineOption{
				ImageURL: cfg.ImageURL,
			},
			Inputs:  inputs,
			Outputs: outputs,
		},
		Spec: aicc.SpecOption{
			Resource: aicc.ResourceOption{
				FlavorId:  cfg.FlavorId,
				PoolId:    cfg.PoolId,
				PoolName:  cfg.PoolName,
				NodeCount: 1,
			},
			LogExportPath: aicc.LogExportPathOption{
				OBSURL: logDir,
			},
		},
	}

	impl.genJobParameter(t, &opt)

	info.JobId, err = impl.cli.createJob(opt)

	if err == nil {
		p := ""
		info.LogDir = strings.TrimPrefix(logDir, p)
		info.OutputDir = strings.TrimPrefix(outputDir, p)
	}

	return
}

func (impl aiccFinetuneImpl) GetDetail(jobId string) (r domain.JobDetail, err error) {
	v, err := impl.cli.getJob(jobId)
	if err != nil {
		return
	}

	if status, ok := statusMap[strings.ToLower(v.Status.Phase)]; ok {
		r.Status = status
	} else {
		r.Status = domain.TrainingStatusFailed
	}

	// convert millisecond to second
	r.Duration = v.Status.Duration / 1000

	return
}
func (impl aiccFinetuneImpl) Delete(jobId string) error {
	err := impl.cli.deleteJob(jobId)
	return err
}

func (impl aiccFinetuneImpl) GetLogFilePath(jobId string) (string, error) {
	return impl.cli.getLogURL(jobId)
}

func (impl aiccFinetuneImpl) GenOutput(jobId string) (string, error) {
	return impl.uploadFolder(jobId)
}

func (impl aiccFinetuneImpl) GetLogDownloadURL(outputDir string) (string, error) {
	return impl.cli.getLogURL(outputDir)
}

func (impl aiccFinetuneImpl) Terminate(jobId string) error {
	return impl.cli.terminateJob(jobId)
}
