package main

import (
	"github.com/bilibili-danmaku/handler"
	"github.com/bilibili-danmaku/search"
	"github.com/gin-gonic/gin"
)

const (
	initDataSet = "2018-09-29-titles-sort.csv"
)

func init() {
	// search.SearchHelper = search.InitWithFiles(initDataSet)
	search.SearchEngine = search.InitEngineWithFiles(initDataSet)
}

func main() {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/danmaku", handler.Danmaku)
	r.Run(":9898") // listen and serve on 0.0.0.0:9898
}
