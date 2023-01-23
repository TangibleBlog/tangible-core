package globalstruct

type IndexStruct struct {
	PageList []map[string]interface{}
	MenuList []map[string]interface{}
	PageInfo IndexPageInfoStruct
}

type PostStruct struct {
	MetaData   map[string]interface{}
	Content    string
	MenuList   []map[string]interface{}
	PageConfig map[string]interface{}
}
