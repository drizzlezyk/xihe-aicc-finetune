package aicc

type JobCreateOption struct {
	Kind      string          `json:"kind" required:"true"`
	Metadata  MetadataOption  `json:"metadata" required:"true"`
	Algorithm AlgorithmOption `json:"algorithm"`
	Spec      SpecOption      `json:"spec"`
}

type MetadataOption struct {
	Name string `json:"name" required:"true"`
	Desc string `json:"description"`
}

type AlgorithmOption struct {
	WorkingDir   string              `json:"working_dir"`
	CodeDir      string              `json:"code_dir"`
	Command      string              `json:"command"`
	Engine       EngineOption        `json:"engine"`
	Parameters   []ParameterOption   `json:"parameters"`
	Environments map[string]string   `json:"environments"`
	Inputs       []InputOutputOption `json:"inputs"`
	Outputs      []InputOutputOption `json:"outputs"`
}

type EngineOption struct {
	ImageURL string `json:"image_url"`
}

type ParameterOption struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type InputOutputOption struct {
	Name   string       `json:"name" required:"true"`
	Remote RemoteOption `json:"remote" required:"true"`
}

type RemoteOption struct {
	OBS OBSOption `json:"obs" required:"true"`
}

type OBSOption struct {
	OBSURL string `json:"obs_url" required:"true"`
}

type SpecOption struct {
	Resource      ResourceOption      `json:"resource"`
	LogExportPath LogExportPathOption `json:"log_export_path"`
}

type ResourceOption struct {
	FlavorId  string `json:"flavor_id" Required:"true"`
	PoolId    string `json:"pool_id" Required:"true"`
	PoolName  string `json:"pool_name" Required:"true"`
	NodeCount int    `json:"node_count,omitempty"`
}

type LogExportPathOption struct {
	OBSURL string `json:"obs_url,omitempty"`
}
