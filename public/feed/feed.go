package feed

import (
	"fmt"
	"log"
	"os"
	"tangible-core/public/common"
	"tangible-core/public/globalstruct"
	"tangible-core/public/template"
	"tangible-core/utils"
	"time"

	"github.com/gorilla/feeds"
)

func GenerateFeed(config globalstruct.SystemConfigStruct) {

	feed := &feeds.Feed{
		Title:       config.Project["Name"].(string),
		Link:        &feeds.Link{Href: config.Project["URL"].(string)},
		Description: config.Project["Description"].(string),
		Author:      &feeds.Author{Name: config.Author["Name"].(string)},
		Created:     time.Now(),
	}

	feed.Items = []*feeds.Item{}
	documents := common.LoadDocuments()
	for _, item := range documents {
		feed.Items = append(feed.Items, &feeds.Item{
			Id:          item["UUID"].(string),
			Title:       item["Title"].(string),
			Link:        &feeds.Link{Href: fmt.Sprintf("%s/post/%s", config.Project["URL"], item["Name"])},
			Description: item["Excerpt"].(string),
			Created:     time.Unix(int64(item["Time"].(float64)), 0),
			Content:     template.BuildHTML(utils.OpenFile(fmt.Sprintf("./documents/post/%s.md", item["UUID"]))),
		})

	}
	atom, err := feed.ToAtom()
	if err != nil {
		log.Fatal(err)
	}
	err = os.WriteFile("./documents/rss.xml", []byte(atom), 0600)
	if err != nil {
		return
	}
	log.Println("Generate feed is complete.")
}
