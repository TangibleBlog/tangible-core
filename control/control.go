package control

import (
	"encoding/json"
	"log"
	"net/http"
	"tangible-core/public/globalstruct"
	"tangible-core/utils"

	"github.com/gin-gonic/gin"
)

//goland:noinspection GoUnusedGlobalVariable
var (
	SystemConfig globalstruct.SystemConfigStruct
	Config       globalstruct.ControlConfigStruct
)

func StartWebService(config globalstruct.SystemConfigStruct, router *gin.Engine) *gin.Engine {

	SystemConfig = config
	err := json.Unmarshal(utils.OpenFile("./config/control.json"), &Config)
	if err != nil {
		log.Panic(err)
	}
	control := router.Group("/control")
	control.GET("/", func(context *gin.Context) {
		context.AbortWithStatus(http.StatusNoContent)
	})
	control.GET("/systemInfo", handleGetSystemInfo)
	control.GET("/list/:type", handleGetList)
	control.GET("/post/:uuid", handleGetPost)
	control.GET("/page/:name", handleGetPage)
	//The following routes need to use security verification
	control.Use(handlerSecurity())
	control.POST("/post/add", handleAddPost)
	control.POST("/post/edit/:uuid", handleEditPost)
	control.POST("/post/delete/:uuid", handleDelPost)
	control.POST("/page/edit/:name", handleEditPage)
	return router
}
