package aiccfinetune

import (
	"github.com/opensourceways/xihe-aicc-finetune/domain"
)

type AICCFinetune interface {
	Create(*domain.AICCFinetune) (domain.JobInfo, error)

	Delete(string) error

	CreateInference(*domain.AICCFinetune) (domain.JobInfo, error)

	// GetLogDownloadURL returns the log url which can be used
	// to download the log of running finetune.
	GetLogDownloadURL(string) (string, error)

	GetDetail(string) (domain.JobDetail, error)

	// GetLogFilePath return the obs path of log
	GetLogFilePath(logDir string) (string, error)

	// GenOutput generates the zip file of output dir and
	// return the obs path of that file.
	GenOutput(outputDir string) (string, error)

	// GenFileDownloadURL generate the temprary
	// download url of obs file.
	GenFileDownloadURL(p string) (string, error)

	Terminate(string) error
}
