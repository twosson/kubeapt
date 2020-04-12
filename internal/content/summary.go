package content

var _ Content = (*Summary)(nil)

type Summary struct {
	Type     string    `json:"type"`
	Title    string    `json:"title"`
	Sections []Section `json:"sections"`
}

func (s *Summary) IsEmpty() bool {
	return len(s.Sections) == 0
}

type Section struct {
	Title string `json:"title"`
	Items []Item `json:"items"`
}

type Item struct {
	Type  string      `json:"type"`
	Label string      `json:"label"`
	Data  interface{} `json:"data"`
}

func TextItem(label, text string) Item {
	return Item{
		Type:  "text",
		Label: label,
		Data: map[string]interface{}{
			"value": text,
		},
	}
}

func LinkItem(label, value, link string) Item {
	return Item{
		Type:  "link",
		Label: label,
		Data: map[string]interface{}{
			"value": value,
			"ref":   link,
		},
	}
}

func JSONItem(label string, blob interface{}) Item {
	return Item{
		Type:  "json",
		Label: label,
		Data:  blob,
	}
}

func NewSummary(title string, sections []Section) Summary {
	return Summary{
		Type:     "summary",
		Title:    title,
		Sections: sections,
	}
}
