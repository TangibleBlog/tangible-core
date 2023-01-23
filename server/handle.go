package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"tangible-core/public/globalstruct"
	"tangible-core/public/net"
	"tangible-core/public/template"
	"tangible-core/utils"
	"time"

	"github.com/gin-gonic/gin"
)

func indexHandle(c *gin.Context) {
	index := c.Params.ByName("page")
	var indexPage int
	if index == "" {
		indexPage = 1
	} else {
		var err error
		indexPage, err = strconv.Atoi(index)
		if err != nil {
			log.Println(err)
		}
	}
	cacheKey := "index:" + index
	found, result := Read(Cache, cacheKey)
	if found && !strings.Contains(c.ContentType(), "application/json") {
		c.Data(http.StatusOK, net.CONTENT_TYPE_TEXT, result)
		return
	}
	var pageConfig globalstruct.IndexStruct
	var indexPageInfo globalstruct.IndexPageInfoStruct
	pageConfig.PageList = PageList
	pageConfig.MenuList = MenuList
	indexPageInfo.Title = "index"
	indexPageInfo.NowPage = indexPage
	indexPageInfo.PageRow = template.GetPageRow(SystemConfig.Paging, len(PageList))
	pageConfig.PageInfo = indexPageInfo
	if strings.Contains(c.ContentType(), "application/json") {
		System := SystemConfig
		System.Renderer = nil
		System.ServerAddr = ""
		var status, content = template.BuildRestfulIndex(SystemConfig, pageConfig, Template)
		c.JSON(status, content)
	} else {
		var status, content = template.BuildIndex(SystemConfig, pageConfig, Template)
		if status == http.StatusNotFound {
			c.AbortWithStatus(http.StatusNotFound)
			return
		}
		Write(Cache, cacheKey, content)
		c.Data(status, net.CONTENT_TYPE_TEXT, content)
	}
}

func rssHandle(c *gin.Context) {
	found, content := Read(Cache, "RSS")
	if found {
		c.Data(http.StatusOK, net.CONTENT_TYPE_RSS, content)
		return
	}
	if utils.CheckFileExist("./documents/rss.xml") {
		c.Data(http.StatusOK, net.CONTENT_TYPE_RSS, utils.OpenFile("./documents/rss.xml"))
	}
	c.AbortWithStatus(http.StatusNotFound)
}

func postHandle(c *gin.Context) {
	postName := c.Param("post")
	uuid := UUIDList[postName]
	cacheKey := "post:" + uuid
	found, result := Read(Cache, cacheKey)
	if found && !strings.Contains(c.ContentType(), "application/json") {
		log.Print("Use cached data.")
		c.Data(http.StatusOK, net.CONTENT_TYPE_TEXT, result)
		return
	}

	var pageConfig globalstruct.PostStruct
	pageConfig.MenuList = MenuList
	if _, ok := UUIDList[postName]; !ok {
		fileName := fmt.Sprintf("./documents/page/%s.json", postName)
		if utils.CheckFileExist(fileName) {
			c.Redirect(301, "/page/"+postName)
			return
		}
		c.AbortWithStatus(404)
		return
	}
	pageConfig.MetaData = PageMetaVO[uuid]
	fileName := fmt.Sprintf("./documents/post/%s.md", uuid)
	if !utils.CheckFileExist(fileName) {
		c.AbortWithStatus(404)
		return
	}

	pageConfig.MetaData, pageConfig.Content = getPageContent(pageConfig, fileName)
	if strings.Contains(c.ContentType(), "application/json") {
		System := SystemConfig
		System.Renderer = nil
		System.ServerAddr = ""
		content := map[string]any{"System": System, "NowTime": utils.GetTimeObj(time.Now()), "Page": pageConfig, "Template": Template.TemplateConfig}
		c.JSON(http.StatusOK, content)
	} else {
		content := template.BuildPost(SystemConfig, pageConfig, Template)
		Write(Cache, cacheKey, content)
		c.Data(http.StatusOK, net.CONTENT_TYPE_TEXT, content)
	}
}

func pageHandle(c *gin.Context) {
	postName := c.Param("page")
	cacheKey := "page:" + postName
	found, result := Read(Cache, cacheKey)
	if found && !strings.Contains(c.ContentType(), "application/json") {
		log.Print("Use cached data.")
		c.Data(http.StatusOK, net.CONTENT_TYPE_TEXT, result)
		return
	}
	var pageConfig globalstruct.PostStruct
	pageConfig.MenuList = MenuList
	fileName := fmt.Sprintf("./documents/page/%s", postName)
	if !utils.CheckFileExist(fileName + ".json") {
		c.AbortWithStatus(404)
		return
	}
	err := json.Unmarshal(utils.OpenFile(fileName+".json"), &pageConfig.MetaData)
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(404)
		return
	}
	if utils.CheckFileExist(fileName + ".html") {
		pageConfig.Content = string(utils.OpenFile(fileName + ".html"))
	} else if utils.CheckFileExist(fileName + ".md") {
		pageConfig.MetaData, pageConfig.Content = getPageContent(pageConfig, fileName+".md")
	} else {
		c.AbortWithStatus(404)
		return
	}
	if strings.Contains(c.ContentType(), "application/json") {
		System := SystemConfig
		System.Renderer = nil
		System.ServerAddr = ""
		content := map[string]any{"System": System, "NowTime": utils.GetTimeObj(time.Now()), "Page": pageConfig, "Template": Template.TemplateConfig}
		c.JSON(http.StatusOK, content)
	} else {
		content := template.BuildPage(SystemConfig, pageConfig, Template)
		Write(Cache, cacheKey, content)
		c.Data(http.StatusOK, net.CONTENT_TYPE_TEXT, content)
	}
}

func getPageContent(pageConfig globalstruct.PostStruct, fileName string) (map[string]interface{}, string) {
	var result string
	meta := pageConfig.MetaData
	if _, ok := pageConfig.MetaData["Renderer"]; ok {
		if value, ok := SystemConfig.Renderer[pageConfig.MetaData["Renderer"].(string)]; ok {
			meta, result = template.ExtensionRenderer(value, pageConfig.MetaData, utils.OpenFile(fileName))
		}
	} else {
		result = template.BuildHTML(utils.OpenFile(fileName))
	}
	return meta, result
}
