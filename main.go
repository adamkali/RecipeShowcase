package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func index(c *gin.Context) {
    var html, err = ioutil.ReadFile("./pages/index.html")
    if err != nil { log.Fatal(err) }

    c.Header("Content-Type", "text/html; charset=utf-8")
    c.Data( http.StatusOK, "text/html; charset=utf-8", html,)
}


func main() {
	r := gin.Default()

    r.LoadHTMLGlob("templates/**/*")
    r.Static("/static", "./static")  

    r.GET("/", index)
	r.Run(":8000")
}
