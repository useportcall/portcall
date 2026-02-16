package downscale

type NodePool struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Size  string `json:"size"`
	Count int    `json:"count"`
}

type DBRow struct {
	ID     string
	Name   string
	Engine string
	Size   string
	Nodes  int
	Status string
}

type RunCmdOutWithEnv func(env map[string]string, name string, args ...string) (string, error)
