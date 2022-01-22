package local_exec

type LocalExec struct {
	Label      string `hcl:",key"`
	WorkingDir string `hcl:"working_dir"`
	Command    string `hcl:"command"`
}

func New(other LocalExec) LocalExec {
	other.Label = "local-exec"
	return other
}
