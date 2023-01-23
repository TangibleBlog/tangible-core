package template

import (
	"log"
	"net/http"
	"strings"
	"tangible-core/public/globalstruct"
	"tangible-core/utils"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/flosch/pongo2/v6"
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/html"
	"github.com/gomarkdown/markdown/parser"
)

func BuildHTML(markdownFileByte []byte) string {
	extensions :=
		parser.CommonExtensions | parser.Tables | parser.FencedCode | parser.Autolink | parser.Strikethrough | parser.MathJax
	withExtensions := parser.NewWithExtensions(extensions)

	return string(markdown.ToHTML(markdownFileByte, withExtensions, nil))

}

func BuildIndex(systemConfig globalstruct.SystemConfigStruct, page globalstruct.IndexStruct, template globalstruct.Template) (int, []byte) {
	indexList := page.PageList
	if page.PageInfo.NowPage == 0 || page.PageInfo.NowPage > page.PageInfo.PageRow {
		return http.StatusNotFound, nil
	}
	if len(page.PageList) > systemConfig.Paging {
		startNum := -systemConfig.Paging + (page.PageInfo.NowPage * systemConfig.Paging)
		endNum := startNum + systemConfig.Paging
		if endNum >= len(page.PageList) {
			endNum = len(page.PageList)
		}
		indexList = page.PageList[startNum:endNum]
	}
	page.PageList = indexList
	context := pongo2.Context{"System": systemConfig, "NowTime": utils.GetTimeObj(time.Now()), "Page": page, "Template": template.TemplateConfig}
	out, err := template.IndexTemplate.ExecuteBytes(context)
	if err != nil {
		log.Panic(err)
	}
	return http.StatusOK, out
}

func BuildPost(systemConfig globalstruct.SystemConfigStruct, page globalstruct.PostStruct, template globalstruct.Template) []byte {
	context := pongo2.Context{"System": systemConfig, "NowTime": utils.GetTimeObj(time.Now()), "Page": page, "Template": template.TemplateConfig}
	out, err := template.PostTemplate.ExecuteBytes(context)
	if err != nil {
		log.Panic(err)
	}
	return out
}
func BuildPage(systemConfig globalstruct.SystemConfigStruct, page globalstruct.PostStruct, template globalstruct.Template) []byte {
	context := pongo2.Context{"System": systemConfig, "NowTime": utils.GetTimeObj(time.Now()), "Page": page, "Template": template.TemplateConfig}
	out, err := template.PageTemplate.ExecuteBytes(context)
	if err != nil {
		log.Panic(err)
	}
	return out
}

func BuildRestfulIndex(systemConfig globalstruct.SystemConfigStruct, page globalstruct.IndexStruct, template globalstruct.Template) (int, map[string]any) {
	indexList := page.PageList
	if page.PageInfo.NowPage == 0 || page.PageInfo.NowPage > page.PageInfo.PageRow {
		return http.StatusNotFound, map[string]any{}
	}
	if len(page.PageList) > systemConfig.Paging {
		startNum := -systemConfig.Paging + (page.PageInfo.NowPage * systemConfig.Paging)
		endNum := startNum + systemConfig.Paging
		if endNum >= len(page.PageList) {
			endNum = len(page.PageList)
		}
		indexList = page.PageList[startNum:endNum]
	}
	page.PageList = indexList
	context := map[string]any{"System": systemConfig, "NowTime": utils.GetTimeObj(time.Now()), "Page": page, "Template": template.TemplateConfig}
	return http.StatusOK, context
}

func GetPageRow(paging int, pageListSize int) int {
	quotient, denominator := divMod(pageListSize, paging)
	pageRow := quotient
	if denominator != 0 && pageListSize > paging {
		pageRow = quotient + 1
	}
	if pageRow == 0 {
		pageRow = 1
	}
	return pageRow
}

func divMod(numerator, denominator int) (quotient, remainder int) {
	quotient = numerator / denominator // integer division, decimals are truncated
	remainder = numerator % denominator
	return
}

func BuildPlainText(markdownFileByte []byte) string {
	htmlFlags := html.CommonFlags | html.SkipHTML | html.SkipImages
	opts := html.RendererOptions{Flags: htmlFlags}
	renderer := html.NewRenderer(opts)
	htmlRaw := string(markdown.ToHTML(markdownFileByte, nil, renderer))
	dom, err := goquery.NewDocumentFromReader(strings.NewReader("<div id='content'>" + htmlRaw + "</div>"))
	if err != nil {
		log.Println(err)
		return ""
	}
	return dom.Find("#content").Text()
}
