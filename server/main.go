package server

import (
	"fmt"
	"log"
	"tangible-core/public/common"
	"tangible-core/public/feed"
	"tangible-core/public/globalstruct"
	"tangible-core/public/template"
	"tangible-core/utils"
	"time"

	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
)

//goland:noinspection GoSnakeCaseUsage
var (
	SystemConfig globalstruct.SystemConfigStruct
	PageList     []map[string]interface{}
	PageMetaVO   map[string]map[string]interface{}
	UUIDList     map[string]string
	MenuList     []map[string]interface{}
	Cache        Struct
	Template     globalstruct.Template
)

func StartWebService(systemConfig globalstruct.SystemConfigStruct, router *gin.Engine) *gin.Engine {
	err := pongo2.RegisterFilter("json_script", common.ToJSONFilter)
	if err != nil {
		log.Fatalln(err)
	}

	SystemConfig = systemConfig
	Cache = LoadCache()
	go loadTemplates()
	go loadDocument()
	go loadRss()
	log.Println("Listening and serving HTTP on " + SystemConfig.ServerAddr)
	router.GET("/", indexHandle)
	router.GET("/index/:page", indexHandle)
	router.GET("/post/:post", postHandle)
	router.GET("/page/:page", pageHandle)
	router.GET("/rss/", rssHandle)
	router.GET("/feed/", rssHandle)
	router.StaticFS("/static/", gin.Dir("./templates/static", false))
	router.StaticFS("/public/", gin.Dir(fmt.Sprintf("./templates/%s/public", SystemConfig.Theme), false))
	router.DELETE("/cache/", flushHandle)
	router.DELETE("/cache/:flag", flushHandle)
	return router
}
func flushHandle(c *gin.Context) {
	if c.Param("flag") == "rss" {
		FlushCache(true)
		return
	}
	FlushCache(false)
}

func FlushCache(genRss bool) {
	if genRss {
		feed.GenerateFeed(SystemConfig)
	}
	Flush(Cache)
	loadTemplates()
	loadDocument()
	loadRss()
}
func loadRss() {
	log.Println("Start loading Rss...")
	if !utils.CheckFileExist("./documents/rss.xml") {
		feed.GenerateFeed(SystemConfig)
	}
	Write(Cache, "RSS", utils.OpenFile("./documents/rss.xml"))
}
func loadTemplates() {
	log.Println("Start loading templates...")
	Template = template.LoadTemplates(SystemConfig)
}

func loadDocument() {
	log.Println("Start loading documents...")
	PageList, PageMetaVO, UUIDList = processingMetadata(common.LoadDocuments())
	MenuList = common.LoadMenu()
}

func processingMetadata(pageIndex []map[string]interface{}) ([]map[string]interface{}, map[string]map[string]interface{}, map[string]string) {
	formatPageIndex := pageIndex
	uuidList := make(map[string]string)
	pageMetaVO := make(map[string]map[string]interface{})

	for key, value := range pageIndex {
		uuidList[value["Name"].(string)] = value["UUID"].(string)
		formatPageIndex[key]["FormatTime"] = utils.GetFormatTime(int64(value["Time"].(float64)), SystemConfig.TimeFormat)
		formatPageIndex[key]["ISOTime"] = time.Unix(int64(value["Time"].(float64)), 0).Format(time.RFC3339)
		pageMetaVO[value["UUID"].(string)] = value
	}
	log.Println("The article list was loaded successfully.")
	return formatPageIndex, pageMetaVO, uuidList
}
