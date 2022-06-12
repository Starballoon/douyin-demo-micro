package main

import (
	"douyin-demo-micro/cmd/api/handlers"
	"github.com/gin-gonic/gin"
)

func initRouter(r *gin.Engine) {
	auth, identify, err := handlers.AuthWithHelper()
	if err != nil {
		panic(err)
	}

	apiRouter := r.Group("/douyin")

	//basic apis
	apiRouter.GET("/feed/", identify, handlers.Feed)
	apiRouter.POST("/user/login/", auth.LoginHandler)
	// 这里把注册和登录串起来了，注册把合法的用户信息落库，Login再从库里查是否合法，响应由LoginHandler返回
	apiRouter.POST("/user/register/", handlers.Register, auth.LoginHandler)
	//
	// Authentication only needed next
	apiRouter.Use(auth.MiddlewareFunc())
	apiRouter.GET("/user/", handlers.UserInfo)
	apiRouter.POST("/publish/action/", handlers.Publish)
	apiRouter.GET("/publish/list/", handlers.PublishList)
	//
	//// extra apis - I
	apiRouter.POST("/favorite/action/", handlers.FavoriteAction)
	apiRouter.GET("/favorite/list/", handlers.FavoriteList)
	apiRouter.POST("/comment/action/", handlers.CommentAction)
	apiRouter.GET("/comment/list/", handlers.CommentList)
	//
	//// extra apis - II
	apiRouter.POST("/relation/action/", handlers.RelationAction)
	apiRouter.GET("/relation/follow/list/", handlers.FollowList)
	apiRouter.GET("/relation/follower/list/", handlers.FollowerList)
}
