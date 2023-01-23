package common

import (
	"encoding/json"
	"github.com/flosch/pongo2/v6"
)

func ToJSONFilter(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
	output, err := json.Marshal(in.Interface())
	if err != nil {
		return nil, &pongo2.Error{
			Template:  nil,
			Filename:  "",
			Line:      0,
			Column:    0,
			Token:     nil,
			Sender:    "ToJSONFilter",
			OrigError: err,
		}
	}
	return pongo2.AsSafeValue("<script type=\"application/json\" id=\"" + param.String() + "\">" + string(output) + "</script>"), nil
}
