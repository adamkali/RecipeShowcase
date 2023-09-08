package main

import (
	"bytes"
	"net/http"

	"github.com/gin-gonic/gin"
)

//type RecipeController struct {
//	ID          uuid.UUID
//	Name        string
//	Ingredients []string
//}

func index(c *gin.Context) {
    var body bytes.Buffer

    

    c.Header("Content-Type", "text/html; charset=utf-8")
    c.Data(
        http.StatusOK,
    )
}

func main() {
	r := gin.Default()

    r.Static("./static")

    r.GET("/", index)
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
