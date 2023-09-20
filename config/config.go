package config

import (
	"github.com/opensourceways/community-robot-lib/utils"
	"github.com/opensourceways/xihe-aicc-finetune/infrastructure/watchimpl"
)

type configSetDefault interface {
	SetDefault()
}

type configValidate interface {
	Validate() error
}

type Config struct {
	Watch    watchimpl.Config `json:"watch"        required:"true"`
	Finetune FinetuneConfig   `json:"finetune"     required:"true"`
	AICC     AICCConfig       `json:"aicc"         required:"true"`
	Upload   UploadConfig     `json:"upload"         required:"true"`
	OBS      OBSConfig        `json:"obs"         required:"true"`
}

func (cfg *Config) configItems() []interface{} {
	return []interface{}{
		&cfg.Watch,
		&cfg.Finetune,
		&cfg.AICC,
	}
}

func (cfg *Config) validate() error {
	if _, err := utils.BuildRequestBody(cfg, ""); err != nil {
		return err
	}

	items := cfg.configItems()

	for _, i := range items {
		if v, ok := i.(configValidate); ok {
			if err := v.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cfg *Config) setDefault() {
	items := cfg.configItems()

	for _, i := range items {
		if v, ok := i.(configSetDefault); ok {
			v.SetDefault()
		}
	}
}

func LoadConfig(path string) (*Config, error) {
	v := new(Config)

	if err := utils.LoadFromYaml(path, v); err != nil {
		return nil, err
	}

	v.setDefault()

	if err := v.validate(); err != nil {
		return nil, err
	}

	return v, nil
}

func (cfg *Config) Validate() error {
	items := cfg.configItems()

	for _, i := range items {
		if v, ok := i.(configValidate); ok {
			if err := v.Validate(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (cfg *Config) SetDefault() {
	items := cfg.configItems()

	for _, i := range items {
		if v, ok := i.(configSetDefault); ok {
			v.SetDefault()
		}
	}
}

type AICCConfig struct {
	Domain       string `json:"domain" required:"true"`
	User         string `json:"user" required:"true"`
	Password     string `json:"password" required:"true"`
	Project      string `json:"project" required:"true"`
	ProjectId    string `json:"project_id" required:"true"`
	AuthEndpoint string `json:"auth_endpoint" required:"true"`

	// modelarts endpoint
	Endpoint string `json:"endpoint" required:"true"`
}

type FinetuneConfig struct {
	WukongConfig ModelConfig `json:"wukong"`
}

type ModelConfig struct {
	PoolId     string `json:"pool_id"`
	PoolName   string `json:"pool_name"`
	FlavorId   string `json:"flavor_id"`
	OutputKey  string `json:"output_key"`
	OutputDir  string `json:"output_dir"`
	WorkingDir string `json:"working_dir"`
	CodeDir    string `json:"code_dir"`
	InputDir   string `json:"input_dir"`
	LogDir     string `json:"log_dir"`
	ModelDir   string `json:"ckpt_file"`
	ImageURL   string `json:"image_url"`
}

type UploadConfig struct {
	UploadWorkDir     string `json:"upload_work_dir"      required:"true"`
	UploadFolderShell string `json:"upload_folder_shell"  required:"true"`

	// DownloadExpiry specifies the timeout to download a obs file.
	// The unit is second.
	DownloadExpiry int    `json:"download_expiry"`
	OBSUtilPath    string `json:"obsutil_path"             required:"true"`
}

func (c *UploadConfig) SetDefault() {
	if c.DownloadExpiry <= 0 {
		c.DownloadExpiry = 3600
	}
}

type OBSConfig struct {
	AccessKey string `json:"access_key"    required:"true"`
	SecretKey string `json:"secret_key"    required:"true"`
	Endpoint  string `json:"endpoint"      required:"true"`
	Bucket    string `json:"bucket"        required:"true"`
}
