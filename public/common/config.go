package common

import (
	"log"
	"tangible-core/public/globalstruct"
	"tangible-core/utils"

	jsoniter "github.com/json-iterator/go"
)

func LoadSystemConfig() globalstruct.SystemConfigStruct {
	systemConfig := globalstruct.DefaultSystemConfig()
	err := jsoniter.Unmarshal(utils.OpenFile("./config/system.json"), &systemConfig)
	HandleError(err)
	if utils.CheckFileExist("./config/renderer.json") {
		err = jsoniter.Unmarshal(utils.OpenFile("./config/renderer.json"), &systemConfig.Renderer)
		HandleError(err)
	}
	log.Println("The system configuration was loaded successfully.")
	return systemConfig
}

func LoadDocuments() []map[string]interface{} {
	var pageIndex []map[string]interface{}
	err := jsoniter.Unmarshal(utils.OpenFile("./documents/index.json"), &pageIndex)
	HandleError(err)
	return pageIndex
}

func LoadMenu() []map[string]interface{} {
	var menuIndex []map[string]interface{}
	err := jsoniter.Unmarshal(utils.OpenFile("./config/menu.json"), &menuIndex)
	HandleError(err)
	return menuIndex
}
