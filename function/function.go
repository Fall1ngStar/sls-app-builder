package function

type Function struct {
	Handler string
	Layers  []string
	//Events  []Event `yaml:",omitempty"`
}

type Event struct {
}

type TemplateFunction struct {
	FunctionSnake string
}
