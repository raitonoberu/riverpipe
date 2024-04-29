package event

type FocusedOutput struct {
	Output uint32 `json:"output"`
}

func (e FocusedOutput) Event() string {
	return "focused_output"
}

type UnfocusedOutput struct {
	Output uint32 `json:"output"`
}

func (e UnfocusedOutput) Event() string {
	return "unfocused_output"
}

type FocusedView struct {
	Title string `json:"title"`
}

func (e FocusedView) Event() string {
	return "focused_view"
}

type Mode struct {
	Name string `json:"name"`
}

func (e Mode) Event() string {
	return "mode"
}
