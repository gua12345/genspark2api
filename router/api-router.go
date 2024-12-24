package router

import (
    "genspark2api/controller"
    "genspark2api/middleware"
    "github.com/gin-gonic/gin"
)

func SetApiRouter(router *gin.Engine) {
    router.Use(middleware.CORS())
    //router.Use(gzip.Gzip(gzip.DefaultCompression))
    router.Use(middleware.RequestRateLimit())

    // OpenAI API 路由组
    v1Router := router.Group("/v1")
    v1Router.Use(middleware.OpenAIAuth())
    v1Router.POST("/chat/completions", controller.ChatForOpenAI)
    v1Router.POST("/images/generations", controller.ImagesForOpenAI)
    v1Router.GET("/models", controller.OpenaiModels)

    // token 相关路由
    tokenController := &controller.TokenController{}
    
    // 网页路由
    router.GET("/:password", tokenController.TokenPage)
    
    // token管理API路由
    router.GET("/:password/token/list", tokenController.GetTokens)     // 查看所有 token
    router.POST("/:password/token/append", tokenController.AppendToken) // 追加 token
    router.POST("/:password/token/clear", tokenController.ClearTokens)  // 清空 token
}
