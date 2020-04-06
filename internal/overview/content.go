package overview

type contentResponse struct {
	Contents []content `json:"contents,omitempty"`
}

type content interface {
}

type table struct {
	Type    string        `json:"type,omitempty"`
	Title   string        `json:"title,omitempty"`
	Columns []tableColumn `json:"columns,omitempty"`
	Rows    []tableRow    `json:"rows,omitempty"`
}

type tableColumn struct {
	Name     string `json:"name,omitempty"`
	Accessor string `json:"accessor,omitempty"`
}

type tableRow map[string]string

func newTable(title string) table {
	return table{
		Type:  "table",
		Title: title,
	}
}
