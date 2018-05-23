package main

import (
	"./router"
	"github.com/gin-gonic/gin"
	"runtime"
)

func main() {
    runtime.GOMAXPROCS(4)
    r:=gin.New();
    r.Use(gin.Logger());//, gin.Recovery()
	//r := gin.Default()
    r.GET("/", router.LevelDBCache)
    r.GET("/R",router.RedisCache)
	r.Run(":8081")
}