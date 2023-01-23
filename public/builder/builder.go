package builder

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"tangible-core/public/common"
	"tangible-core/public/globalstruct"
	"tangible-core/public/template"
	"tangible-core/utils"
)

var (
	SystemConfig globalstruct.SystemConfigStruct
	PageList     []map[string]interface{}
	MenuList     []map[string]interface{}
	Template     globalstruct.Template
)

func BuildStaticPage(systemConfig globalstruct.SystemConfigStruct) {
	SystemConfig = systemConfig
	log.Println("Start loading documents...")
	PageList = processingMetadata(common.LoadDocuments())
	MenuList = common.LoadMenu()
	log.Println("Start loading templates...")

	Template = template.LoadTemplates(SystemConfig)
	_, _ = utils.CopyFile("./documents/rss.xml", "./dist/rss.xml")
	_ = os.MkdirAll("./dist/index", 0750)
	_ = os.MkdirAll("./dist/post", 0750)
	_ = os.MkdirAll("./dist/page", 0750)
	pageRow := template.GetPageRow(SystemConfig.Paging, len(PageList))
	i := 1
	for i <= pageRow {
		fileName := fmt.Sprintf("./dist/index/%d.html", i)
		err := os.WriteFile(fileName, buildIndex(i, pageRow), 0644)
		if err != nil {
			log.Panic(err)
			return
		}
		i++
	}
	_, _ = utils.CopyFile("./dist/index/1.html", "./dist/index.html")

	for _, item := range PageList {
		fileName := fmt.Sprintf("./dist/post/%s.html", item["Name"])
		err := os.WriteFile(fileName, buildPost(item), 0644)
		if err != nil {
			log.Panic(err)
		}
	}

	err := filepath.Walk("./documents/page/", func(filepath string, info os.FileInfo, err error) error {
		if strings.HasSuffix(filepath, ".json") {
			var item map[string]interface{}

			_ = json.Unmarshal(utils.OpenFile(filepath), &item)
			item["PageName"] = strings.Replace(path.Base(filepath), ".json", "", -1)
			fileName := fmt.Sprintf("./dist/page/%s.html", item["PageName"])
			err := os.WriteFile(fileName, buildPage(item), 0644)
			if err != nil {
				log.Panic(err)
			}

		}
		return nil
	})
	if err != nil {
		panic(err)
	}
}

func buildPage(metaData map[string]interface{}) []byte {
	var pageConfig globalstruct.PostStruct
	pageConfig.MenuList = MenuList
	fileName := fmt.Sprintf("./documents/page/%s", metaData["PageName"])
	if utils.CheckFileExist(fileName + ".html") {
		pageConfig.Content = string(utils.OpenFile(fileName + ".html"))
	} else if utils.CheckFileExist(fileName + ".md") {
		pageConfig.Content = template.BuildHTML(utils.OpenFile(fileName + ".md"))
	}
	return template.BuildPage(SystemConfig, pageConfig, Template)
}

func buildPost(metaData map[string]interface{}) []byte {
	var pageConfig globalstruct.PostStruct
	pageConfig.MenuList = MenuList
	pageConfig.MetaData = metaData
	fileName := fmt.Sprintf("./documents/post/%s.md", metaData["UUID"])
	pageConfig.Content = template.BuildHTML(utils.OpenFile(fileName))
	return template.BuildPost(SystemConfig, pageConfig, Template)
}

func buildIndex(indexPage int, pageRow int) []byte {
	var pageConfig globalstruct.IndexStruct
	var indexPageInfo globalstruct.IndexPageInfoStruct
	pageConfig.PageList = PageList
	pageConfig.MenuList = MenuList
	indexPageInfo.Title = "index"
	indexPageInfo.NowPage = indexPage
	indexPageInfo.PageRow = pageRow
	pageConfig.PageInfo = indexPageInfo
	var _, content = template.BuildIndex(SystemConfig, pageConfig, Template)
	return content
}

func processingMetadata(pageIndex []map[string]interface{}) []map[string]interface{} {
	formatPageIndex := pageIndex
	for key, value := range pageIndex {
		formatPageIndex[key]["FormatTime"] = utils.GetFormatTime(int64(value["Time"].(float64)), SystemConfig.TimeFormat)
	}
	log.Println("The article list was loaded successfully.")
	return formatPageIndex
}
