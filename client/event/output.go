package event

type FocusedTags struct {
	Tags uint32 `json:"tags"`
}

func (e FocusedTags) Event() string {
	return "focused_tags"
}

type ViewTags struct {
	Tags []uint32 `json:"tags"`
}

func (e ViewTags) Event() string {
	return "view_tags"
}

type UrgentTags struct {
	Tags uint32 `json:"tags"`
}

func (e UrgentTags) Event() string {
	return "urgent_tags"
}

type LayoutName struct {
	Name string `json:"name"`
}

func (e LayoutName) Event() string {
	return "layout_name"
}

type LayoutNameClear struct{}

func (e LayoutNameClear) Event() string {
	return "layout_name_clear"
}
