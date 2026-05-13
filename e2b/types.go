package e2b

type JSONMap map[string]any

type Sandbox struct {
	TemplateID string  `json:"templateID,omitempty"`
	SandboxID  string  `json:"sandboxID,omitempty"`
	State      string  `json:"state,omitempty"`
	Alias      string  `json:"alias,omitempty"`
	StartedAt  string  `json:"startedAt,omitempty"`
	EndAt      string  `json:"endAt,omitempty"`
	Domain     string  `json:"domain,omitempty"`
	Metadata   JSONMap `json:"metadata,omitempty"`
}

type Template struct {
	TemplateID string   `json:"templateID,omitempty"`
	BuildID    string   `json:"buildID,omitempty"`
	Public     bool     `json:"public,omitempty"`
	Names      []string `json:"names,omitempty"`
	Aliases    []string `json:"aliases,omitempty"`
}
