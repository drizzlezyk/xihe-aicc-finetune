package watchimpl

type Config struct {
	// Interval specifies the interval of second between two loops
	// that check all finetunes in a loop.
	Interval int `json:"interval"`

	// Timeout specifies the time that a finetune can live
	// The unit is second.
	Timeout int `json:"timeout"`

	// MaxWatchNum specifies the max num of finetune
	// which the aicc finetune center can support
	MaxWatchNum int `json:"max_watch_num"`

	Endpoint string `json:"endpoint" required:"true"`
}

func (cfg *Config) SetDefault() {
	if cfg.Interval <= 0 {
		cfg.Interval = 10
	}

	if cfg.Timeout <= 0 {
		cfg.Timeout = 864000
	}

	if cfg.MaxWatchNum <= 0 {
		cfg.MaxWatchNum = 100
	}
}
