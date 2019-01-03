package function

type Function struct {
	Handler string
	//Events  []Event `yaml:",omitempty"`
}

type Event struct {
}

type TemplateFunction struct {
	FunctionSnake string
}
