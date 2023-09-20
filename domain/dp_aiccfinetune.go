package domain

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

var (
	reName      = regexp.MustCompile("^[a-zA-Z0-9_-]+$")
	reDirectory = regexp.MustCompile("^[a-zA-Z0-9_/-]+$")
	reFilePath  = regexp.MustCompile("^[a-zA-Z0-9_/.-]+$")

	TrainingStatusFailed      = trainingStatus("Failed")
	TrainingStatusPending     = trainingStatus("Pending")
	TrainingStatusRunning     = trainingStatus("Running")
	TrainingStatusCreating    = trainingStatus("Creating")
	TrainingStatusAbnormal    = trainingStatus("Abnormal")
	TrainingStatusCompleted   = trainingStatus("Completed")
	TrainingStatusTerminated  = trainingStatus("Terminated")
	TrainingStatusTerminating = trainingStatus("Terminating")

	trainingDoneStatus = map[string]bool{
		"Failed":     true,
		"Abnormal":   true,
		"Completed":  true,
		"Terminated": true,
	}
)

// Account
type Account interface {
	Account() string
}

func NewAccount(v string) (Account, error) {
	if v == "" || strings.ToLower(v) == "root" || !reName.MatchString(v) {
		return nil, errors.New("invalid user name")
	}

	return dpAccount(v), nil
}

type dpAccount string

func (r dpAccount) Account() string {
	return string(r)
}

// FinetuneName
type FinetuneName interface {
	FinetuneName() string
}

func NewFinetuneName(v string) (FinetuneName, error) {
	max := 30
	min := 3

	if n := len(v); n > max || n < min {
		return nil, fmt.Errorf("name's length should be between %d to %d", min, max)
	}

	if !reName.MatchString(v) {
		return nil, errors.New("invalid name")
	}

	return finetuneName(v), nil
}

type finetuneName string

func (r finetuneName) FinetuneName() string {
	return string(r)
}

// FinetuneDesc
type FinetuneDesc interface {
	FinetuneDesc() string
}

func NewFinetuneDesc(v string) (FinetuneDesc, error) {
	if v == "" {
		return finetuneDesc(v), nil
	}

	max := 100
	if len([]rune(v)) > max {
		return nil, fmt.Errorf("the length of desc should be less than %d", max)
	}

	return finetuneDesc(v), nil
}

type finetuneDesc string

func (r finetuneDesc) FinetuneDesc() string {
	return string(r)
}

// Directory
type Directory interface {
	Directory() string
	LastDirectory() string
}

func NewDirectory(v string) (Directory, error) {
	if v == "" {
		return directory(""), nil
	}

	if !reDirectory.MatchString(v) {
		return nil, errors.New("invalid directory")
	}

	return directory(v), nil
}

type directory string

func (r directory) Directory() string {
	return string(r)
}

func (r directory) LastDirectory() string {
	s := strings.TrimRight(string(r), "/")
	splitDir := strings.Split(s, "/")
	return splitDir[len(splitDir)-1]
}

// FilePath
type FilePath interface {
	FilePath() string
}

func NewFilePath(v string) (FilePath, error) {
	if v == "" {
		return nil, errors.New("empty file path")
	}

	if !reFilePath.MatchString(v) {
		return nil, errors.New("invalid filePath")
	}

	return filePath(v), nil
}

type filePath string

func (r filePath) FilePath() string {
	return string(r)
}

// CustomizedKey
type CustomizedKey interface {
	CustomizedKey() string
}

func NewCustomizedKey(v string) (CustomizedKey, error) {
	if v == "" {
		return nil, errors.New("empty key")
	}

	return customizedKey(v), nil
}

type customizedKey string

func (r customizedKey) CustomizedKey() string {
	return string(r)
}

// CustomizedValue
type CustomizedValue interface {
	CustomizedValue() string
}

func NewCustomizedValue(v string) (CustomizedValue, error) {
	if v == "" {
		return nil, nil
	}

	return customizedValue(v), nil
}

type customizedValue string

func (r customizedValue) CustomizedValue() string {
	return string(r)
}

// TrainingStatus
type TrainingStatus interface {
	TrainingStatus() string
	IsDone() bool
	IsSuccess() bool
}

type trainingStatus string

func (s trainingStatus) TrainingStatus() string {
	return string(s)
}

func (s trainingStatus) IsDone() bool {
	return trainingDoneStatus[string(s)]
}

func (s trainingStatus) IsSuccess() bool {
	return string(s) == TrainingStatusCompleted.TrainingStatus()
}

// ModelName
type ModelName interface {
	ModelName() string
}

func NewModelName(v string) (ModelName, error) {
	if v == "" {
		return nil, nil
	}

	return modelName(v), nil
}

type modelName string

func (m modelName) ModelName() string {
	return string(m)
}
