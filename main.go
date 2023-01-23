package main

import (
	"log"
	"os"
	"tangible-core/control"
	"tangible-core/public/builder"
	"tangible-core/public/common"
	"tangible-core/public/feed"
	PublicStruct "tangible-core/public/globalstruct"
	"tangible-core/server"
	"tangible-core/tui"
	"tangible-core/utils"

	"github.com/gin-gonic/gin"
	"github.com/urfave/cli/v2"
)

var commands []*cli.Command

func init() {
	if !utils.CheckFileExist("./config/system.json") {
		tui.SetupSystem()
	}
	gin.SetMode(gin.ReleaseMode)
	router := gin.Default()
	var systemConfig PublicStruct.SystemConfigStruct
	systemConfig = common.LoadSystemConfig()
	commands = []*cli.Command{
		{
			Name: "main",
			//Aliases: ,
			Usage:    "Run the frontend web page rendering service.",
			Category: "server",
			Flags:    []cli.Flag{},
			Action: func(c *cli.Context) error {
				router = server.StartWebService(systemConfig, router)
				if systemConfig.Service.Control {
					log.Println("Loading management API routes.")
					router = control.StartWebService(systemConfig, router)
				}
				if systemConfig.Service.Comment {

				}
				err := router.Run(systemConfig.ServerAddr)
				common.HandleError(err)
				return nil
			},
		},
		{
			Name: "control",
			//Aliases: ,
			Usage:    "Only run the management API.",
			Category: "server",
			Flags:    []cli.Flag{},
			Action: func(c *cli.Context) error {
				router = control.StartWebService(systemConfig, router)
				err := router.Run(systemConfig.ServerAddr)
				common.HandleError(err)
				return nil
			},
		},
		{
			Name: "gen-rss",
			//Aliases: ,
			Usage:    "Generate RSS file.",
			Category: "Generate",
			Flags:    []cli.Flag{},
			Action: func(c *cli.Context) error {
				feed.GenerateFeed(systemConfig)
				return nil
			},
		}, {
			Name: "build",
			//Aliases: ,
			Usage:    "Generate Page file.",
			Category: "Generate",
			Flags:    []cli.Flag{},
			Action: func(c *cli.Context) error {
				builder.BuildStaticPage(systemConfig)
				return nil
			},
		},
		{
			Name: "flush-cache",
			//Aliases: ,
			Usage:    "Clear the cache in redis.",
			Category: "Flush",
			Flags:    []cli.Flag{},
			Action: func(c *cli.Context) error {
				cacheHandle := server.LoadCache()
				server.Flush(cacheHandle)
				feed.GenerateFeed(systemConfig)
				log.Println("Start loading Rss...")
				server.Write(cacheHandle, "RSS", utils.OpenFile("./documents/rss.xml"))
				return nil
			},
		},
	}
}
func main() {
	app := cli.NewApp()
	app.Name = "tangibleblog"
	app.Usage = "A fast and lightweight blog framework based on Golang."
	app.HideVersion = false
	//app.Flags = flags
	app.Commands = commands
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
