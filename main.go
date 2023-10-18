package main

import (
	"flag"
	"os"

	"github.com/opensourceways/community-robot-lib/logrusutil"
	liboptions "github.com/opensourceways/community-robot-lib/options"
	"github.com/opensourceways/xihe-aicc-finetune/app"
	"github.com/opensourceways/xihe-aicc-finetune/config"
	"github.com/opensourceways/xihe-aicc-finetune/infrastructure/aiccfinetuneimpl"
	"github.com/opensourceways/xihe-aicc-finetune/infrastructure/watchimpl"
	"github.com/opensourceways/xihe-aicc-finetune/server"
	"github.com/sirupsen/logrus"
)

type options struct {
	service     liboptions.ServiceOptions
	enableDebug bool
}

func (o *options) Validate() error {
	return o.service.Validate()
}

func gatherOptions(fs *flag.FlagSet, args ...string) (options, error) {
	var o options

	o.service.AddFlags(fs)

	fs.BoolVar(
		&o.enableDebug, "enable_debug", false,
		"whether to enable debug model.",
	)

	err := fs.Parse(args)

	return o, err
}

func main() {
	logrusutil.ComponentInit("xihe-aicc-finetune")
	log := logrus.NewEntry(logrus.StandardLogger())

	o, err := gatherOptions(
		flag.NewFlagSet(os.Args[0], flag.ExitOnError),
		os.Args[1:]...,
	)
	if err != nil {
		logrus.Fatalf("new options failed, err:%s", err.Error())
	}

	if err := o.Validate(); err != nil {
		logrus.Fatalf("Invalid options, err:%s", err.Error())
	}

	if o.enableDebug {
		logrus.SetLevel(logrus.DebugLevel)
		logrus.Debug("debug enabled.")
	}

	// cfg
	cfg, err := config.LoadConfig(o.service.ConfigFile)
	if err != nil {
		logrus.Fatalf("load config, err:%s", err.Error())
	}

	// finetune
	as, err := aiccfinetuneimpl.NewAiccFinetune(cfg)
	if err != nil {
		logrus.Errorf("new finetune client failed, err:%s", err.Error())
	}

	// watch
	ws, err := watchimpl.NewWatcher(&cfg.Watch, as)
	if err != nil {
		logrus.Errorf("new watch service failed, err:%s", err.Error())
	}

	service := app.NewAICCFinetuneService(as, ws, log)
	go ws.Run()

	defer ws.Exit()

	server.StartWebServer(&server.Service{
		Log:      log,
		Port:     o.service.Port,
		Timeout:  o.service.GracePeriod,
		Finetune: service,
	})
}
