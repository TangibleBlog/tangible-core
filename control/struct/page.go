package _struct

type PostStruct struct {
	MetaData map[string]interface{}
	Content  string
}
type PostBody struct {
	MetaData map[string]interface{} `json:"MetaData"`
	Content  string                 `json:"Content"`
}
