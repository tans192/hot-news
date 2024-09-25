package routes

import (
	"github.com/gin-gonic/gin"
	"hotNews/http/controllers"
	"hotNews/utils"
)

func Init() {
	gin.SetMode(utils.AppSetting.DebugMode)
	router := gin.Default()
	// v1 api
	v1 := router.Group("/v1")
	{
		v1.GET("/zhihu/top", controllers.QueryHtml)
		v1.GET("/zhihu/detail", controllers.Detail)
		v1.GET("/zhihu/wenku", controllers.Wenku)
	}
	port := utils.AppSetting.Port
	router.Run(":" + port)
}
