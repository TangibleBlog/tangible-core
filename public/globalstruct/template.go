package globalstruct

import "github.com/flosch/pongo2/v6"

type IncludePointStruct struct {
	Header  string
	Footer  string
	Comment string
}

type IndexPageInfoStruct struct {
	NowPage int
	PageRow int
	Title   string
}

type TemplateConfigStruct struct {
	Config  map[string]interface{}
	Include IncludePointStruct
}

type Template struct {
	IndexTemplate  *pongo2.Template
	PostTemplate   *pongo2.Template
	PageTemplate   *pongo2.Template
	TemplateConfig TemplateConfigStruct
}
