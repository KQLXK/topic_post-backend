package main

import (
	"github.com/gin-gonic/gin"
	"youthcamp/lesson02/project/controller"
	"youthcamp/lesson02/project/repository"
)

func main() {

	if err := Init(); err != nil {
		panic(err)
	}

	r := gin.Default()

	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"msg": "pong",
		})
	})

	r.GET("/community/page/get/:id", func(c *gin.Context) {
		topicId := c.Param("id")
		data := controller.QueryPageInfo(topicId)
		c.JSON(200, data)
	})

	r.POST("community/post/do", func(c *gin.Context) {
		uid, _ := c.GetPostForm("uid")
		content, _ := c.GetPostForm("content")
		topicId, _ := c.GetPostForm("topicId")
		data := controller.PublishPost(topicId, uid, content)
		c.JSON(200, data)
	})

	r.Run(":8080")

}

func Init() error {
	if err := repository.InitDB(); err != nil {
		return err
	}
	return nil
}
