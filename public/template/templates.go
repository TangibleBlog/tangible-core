package template

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"tangible-core/public/common"
	"tangible-core/public/globalstruct"
	"tangible-core/utils"

	"github.com/flosch/pongo2/v6"
)

func LoadTemplates(systemConfig globalstruct.SystemConfigStruct) globalstruct.Template {
	location := fmt.Sprintf("./templates/%s/", systemConfig.Theme)
	if _, err := os.Stat(location); os.IsNotExist(err) {
		log.Fatalln(err)
	}
	var metadata map[string]interface{}
	metadataFile := location + "metadata.json"
	if utils.CheckFileExist(metadataFile) {
		log.Println("Load template configuration file.")
		err := json.Unmarshal(utils.OpenFile(metadataFile), &metadata)
		common.HandleError(err)
	}
	var config map[string]interface{}
	configFile := fmt.Sprintf("./config/template/%s.json", systemConfig.Theme)
	if utils.CheckFileExist(configFile) {
		log.Println("Load template configuration file.")
		err := json.Unmarshal(utils.OpenFile(configFile), &config)
		common.HandleError(err)
	}
	var templateConfig globalstruct.TemplateConfigStruct
	templateConfig.Include = LoadInclude()
	templateConfig.Config = config
	index := pongo2.Must(pongo2.FromFile(location + "index.html"))
	post := pongo2.Must(pongo2.FromFile(location + "post.html"))
	page := post
	if utils.CheckFileExist(location + "page.html") {
		page = pongo2.Must(pongo2.FromFile(location + "page.html"))
	}
	log.Println("The template configuration file is loaded.")
	var template globalstruct.Template
	template.IndexTemplate = index
	template.PageTemplate = page
	template.PostTemplate = post
	template.TemplateConfig = templateConfig
	return template
}

func LoadInclude() globalstruct.IncludePointStruct {
	var include globalstruct.IncludePointStruct
	include.Comment = string(utils.UnCheckOpenFile("./include/comment.html"))
	include.Header = string(utils.UnCheckOpenFile("./include/header.html"))
	include.Footer = string(utils.UnCheckOpenFile("./include/footer.html"))
	return include
}
