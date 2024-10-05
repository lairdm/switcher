package monitor

type Command struct {
	Input   Input  `json:"input,omitempty"`
	Monitor string `json:"monitor,omitempty"`
}

func New(input Input, monitor string) *Command {
	return &Command{
		Input:   input,
		Monitor: monitor,
	}
}
