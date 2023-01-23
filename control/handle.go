package control

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"tangible-core/control/db"
	_struct "tangible-core/control/struct"
	"tangible-core/public/common"
	"tangible-core/public/template"
	"tangible-core/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func handleGetSystemInfo(c *gin.Context) {
	systemConfig := common.LoadSystemConfig()
	result := make(map[string]interface{})
	result["Author"] = systemConfig.Author
	result["System"] = systemConfig.Project
	common.DataResult(c, http.StatusOK, result)
}

func handleGetList(c *gin.Context) {
	requestType := c.Param("type")
	var result []map[string]interface{}
	if requestType == "post" {
		result = common.LoadDocuments()
	} else {
		result = common.LoadMenu()
	}
	common.DataResult(c, http.StatusOK, result)
}

func handleGetPost(c *gin.Context) {
	uuidStr := c.Param("uuid")
	fileName := fmt.Sprintf("./documents/post/%s.md", uuidStr)
	if !utils.CheckFileExist(fileName) {
		common.CodeResult(c, http.StatusNotFound)
	}
	var pageConfig _struct.PostStruct
	pageIndex := common.LoadDocuments()
	for _, value := range pageIndex {
		if value["UUID"].(string) == uuidStr {
			pageConfig.MetaData = value
			break
		}
	}
	pageConfig.Content = string(utils.OpenFile(fileName))
	common.DataResult(c, http.StatusOK, pageConfig)
}

func handleGetPage(c *gin.Context) {
	postName := c.Param("name")
	var pageConfig _struct.PostStruct
	fileName := fmt.Sprintf("./documents/page/%s", postName)
	if !utils.CheckFileExist(fileName + ".json") {
		common.CodeResult(c, http.StatusNotFound)
	}
	err := json.Unmarshal(utils.OpenFile(fileName+".json"), &pageConfig.MetaData)
	if err != nil {
		log.Println(err)
		common.CodeResult(c, http.StatusNotFound)
	}
	if utils.CheckFileExist(fileName + ".html") {
		pageConfig.Content = string(utils.OpenFile(fileName + ".html"))
	} else if utils.CheckFileExist(fileName + ".md") {
		pageConfig.Content = string(utils.OpenFile(fileName + ".md"))
	} else {
		common.CodeResult(c, http.StatusNotFound)
	}
	common.DataResult(c, http.StatusOK, pageConfig)
}

func handleAddPost(c *gin.Context) {
	uuidObj := uuid.New()
	fileName := fmt.Sprintf("./documents/post/%s.md", uuidObj.String())
	var body _struct.PostBody
	err := c.BindJSON(&body)
	if err != nil {
		common.CodeResult(c, http.StatusInternalServerError)
	}
	body.MetaData["UUID"] = uuidObj.String()
	if _, ok := body.MetaData["Name"]; !ok {
		body.MetaData["Name"] = body.MetaData["Title"]
	}
	excerpt := ""
	if value, ok := body.MetaData["Excerpt"]; ok {
		excerpt = value.(string)
	} else {
		text := template.BuildPlainText([]byte(body.Content))
		if len(text) >= 140 {
			expertRune := []rune(text)
			excerpt = string(expertRune[:140])
		} else {
			excerpt = text
		}
	}
	body.MetaData["Excerpt"] = excerpt
	body.MetaData["Time"] = time.Now().Unix()
	metaFileName := "./documents/index.json"
	pageIndex := common.LoadDocuments()

	pageIndex = append([]map[string]interface{}{body.MetaData}, pageIndex...)
	marshal, _ := json.Marshal(pageIndex)
	err = os.WriteFile(metaFileName, marshal, 0600)
	if err != nil {
		common.CodeResult(c, http.StatusInternalServerError)
	}
	err = os.WriteFile(fileName, []byte(body.Content), 0600)
	if err != nil {
		common.CodeResult(c, http.StatusInternalServerError)
	}
	metadata, _ := json.Marshal(body.MetaData)
	history := db.History{
		UUID:       uuidObj.String(),
		MetaData:   string(metadata),
		Content:    body.Content,
		CreateTime: time.Now().Unix(),
	}
	db.Dbc.Create(&history)
	common.CodeResult(c, http.StatusOK)
}

func handleEditPost(c *gin.Context) {
	var body _struct.PostBody
	err := c.BindJSON(&body)
	if err != nil {
		common.CodeResult(c, http.StatusInternalServerError)
	}
	uuidStr := c.Param("uuid")
	var pageConfig _struct.PostStruct
	var indexKey int
	pageIndex := common.LoadDocuments()
	for key, value := range pageIndex {
		if value["UUID"].(string) == uuidStr {
			pageConfig.MetaData = value
			indexKey = key
			break
		}
	}
	metaFileName := "./documents/index.json"
	for k, v := range body.MetaData {
		pageConfig.MetaData[k] = v
	}
	excerpt := ""
	if value, ok := body.MetaData["Excerpt"]; ok {
		excerpt = value.(string)
	} else {
		text := template.BuildPlainText([]byte(body.Content))
		if len(text) >= 140 {
			expertRune := []rune(text)
			excerpt = string(expertRune[:140])
		} else {
			excerpt = text
		}
	}
	pageConfig.MetaData["Excerpt"] = excerpt
	pageIndex[indexKey] = pageConfig.MetaData
	marshal, _ := json.Marshal(pageIndex)
	err = os.WriteFile(metaFileName, marshal, 0600)
	if err != nil {
		common.CodeResult(c, http.StatusInternalServerError)
	}
	fileName := fmt.Sprintf("./documents/post/%s.md", uuidStr)
	err = os.WriteFile(fileName, []byte(body.Content), 0644)
	if err != nil {
		common.CodeResult(c, http.StatusInternalServerError)
		return
	}
	metadata, _ := json.Marshal(pageConfig.MetaData)
	history := db.History{
		UUID:       uuidStr,
		MetaData:   string(metadata),
		Content:    body.Content,
		CreateTime: time.Now().Unix(),
	}
	db.Dbc.Create(&history)
	common.CodeResult(c, http.StatusOK)

}

func handleEditPage(c *gin.Context) {
	var body _struct.PostBody
	err := c.BindJSON(&body)
	if err != nil {
		common.CodeResult(c, http.StatusInternalServerError)
	}
	pageName := c.Param("name")
	var pageConfig _struct.PostStruct
	fileName := fmt.Sprintf("./documents/page/%s", pageName)
	err = json.Unmarshal(utils.OpenFile(fileName+".json"), &pageConfig.MetaData)
	if err != nil {
		log.Println(err)
		common.CodeResult(c, http.StatusNotFound)
	}
	for k, v := range body.MetaData {
		pageConfig.MetaData[k] = v
	}
	excerpt := ""
	if value, ok := body.MetaData["Excerpt"]; ok {
		excerpt = value.(string)
	} else {
		text := template.BuildPlainText([]byte(body.Content))
		if len(text) >= 140 {
			expertRune := []rune(text)
			excerpt = string(expertRune[:140])
		} else {
			excerpt = text
		}
	}
	pageConfig.MetaData["Excerpt"] = excerpt
	bytes, err := json.Marshal(pageConfig.MetaData)
	if err != nil {
		return
	}
	_ = os.WriteFile(fileName+".json", bytes, 0644)

	if utils.CheckFileExist(fileName + ".html") {
		err = os.WriteFile(fileName+".html", []byte(body.Content), 0644)
	} else if utils.CheckFileExist(fileName + ".md") {
		err = os.WriteFile(fileName+".md", []byte(body.Content), 0644)
	} else {
		common.CodeResult(c, http.StatusNotFound)
		return
	}
	metadata, _ := json.Marshal(pageConfig.MetaData)
	history := db.History{
		UUID:       pageName,
		MetaData:   string(metadata),
		Content:    body.Content,
		CreateTime: time.Now().Unix(),
	}
	db.Dbc.Updates(&history)
	common.CodeResult(c, http.StatusOK)

}

func handleDelPost(c *gin.Context) {
	uuidStr := c.Param("uuid")
	metaFileName := "./documents/index.json"
	var indexKey int
	pageIndex := common.LoadDocuments()
	for key, value := range pageIndex {
		if value["UUID"].(string) == uuidStr {
			indexKey = key
			break
		}
	}
	pageIndex = append(pageIndex[:indexKey], pageIndex[indexKey+1:]...)
	marshal, _ := json.Marshal(pageIndex)
	err := os.WriteFile(metaFileName, marshal, 0600)
	if err != nil {
		common.CodeResult(c, http.StatusInternalServerError)
		c.Abort()
	}
	fileName := fmt.Sprintf("./documents/post/%s.md", uuidStr)
	err = os.Remove(fileName)
	if err != nil {
		common.CodeResult(c, http.StatusInternalServerError)
		c.Abort()
	}
	common.CodeResult(c, http.StatusOK)
}
