package main

import (
	"log"
	"mircool"
	"net/http"
)

type User struct {
	UserName *string
	PassWord *string
}

func main() {
	server := mircool.NewServer()
	server.GET("/index", func(ctx *mircool.Context) {
		ctx.String(http.StatusOK, "Hello World!")
	})
	server.POST("/login", func(c *mircool.Context) {
		var user User
		if err := c.BindJson(c.Req.Body, &user); err != nil {
			log.Println("参数绑定失败")
		}
		c.JSON(http.StatusOK, mircool.M{
			"userName": user.UserName,
			"passWord": user.PassWord,
		})
	})
	server.Run(":30000")
}