package aicc

type Job struct {
	Metadata JobMetadata `json:"metadata"`
	Status   JobStatus   `json:"status"`
}

type JobMetadata struct {
	Id string `json:"id"`
}

type JobStatus struct {
	Phase     string `json:"phase"`
	Duration  int    `json:"duration"`
	StartTime int    `json:"start_time"`
}

type GetResp struct {
	Response

	Job Job `json:"job"`
}

type CreateResp struct {
	Status string `json:"status"`
	Body   Job    `json:"body"`
}

type CreateBody struct {
	Metadata JobMetadata `json:"metadata"`
}

type Response struct {
	Status int    `json:"status"`
	Msg    string `json:"msg"`
}

type LogResp struct {
	ObsUrl string `json:"obs_url"`
}

type TerminateBody struct {
	ActionType string `json:"action_type"`
}
