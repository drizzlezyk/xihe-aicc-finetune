package domain

type AICCFinetune struct {
	Id    string
	User  Account
	Model ModelName
	Task  string

	AICCFinetuneConfig

	Job       JobInfo
	JobDetail JobDetail
}

type AICCFinetuneConfig struct {
	Name FinetuneName
	Desc FinetuneDesc

	Hyperparameters []KeyValue
	Env             []KeyValue
}

type KeyValue struct {
	Key   CustomizedKey
	Value CustomizedValue
}

type Input struct {
	Key  CustomizedKey
	User Account
	File string
}

type JobInfo struct {
	Endpoint  string
	JobId     string
	LogDir    string
	OutputDir string
}

type JobDetail struct {
	Status     TrainingStatus
	Error      string
	LogPath    string
	OutputPath string
	Duration   int
}

type AICCFinetuneIndex struct {
	User       Account
	Model      string
	FinetuneId string
}

func (t *AICCFinetune) DefaultCommand() string {
	m := map[string]string{
		"wukong": "python train-lora.py",
	}
	command := m[t.Model.ModelName()]

	return command
}

func (t *AICCFinetune) DefaultInferenceCommand() string {
	m := map[string]string{
		"wukong": "python txt2img-lora.py",
	}
	command := m[t.Model.ModelName()]

	return command
}
