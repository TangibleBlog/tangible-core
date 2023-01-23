package template

import (
	"encoding/json"
	"tangible-core/public/net"
)

type postBody struct {
	MetaData map[string]interface{} `json:"MetaData"`
	Content  string                 `json:"Content"`
}

func ExtensionRenderer(renderer string, metaData map[string]interface{}, markdownFileByte []byte) (map[string]interface{}, string) {
	var postMeta postBody
	postMeta.MetaData = metaData
	postMeta.Content = string(markdownFileByte)
	marshal, err := json.Marshal(postMeta)
	if err != nil {
		return nil, ""
	}
	result, err := net.Post(renderer, "application/json;charset=utf-8", marshal)
	err = json.Unmarshal(result, &postMeta)
	if err != nil {
		return nil, ""
	}
	return postMeta.MetaData, postMeta.Content

}
