package main

import (
	"log"
	"net/http"
	"os"

	"github.com/adamkali/RecipeShowcase/src"
	"github.com/adamkali/RecipeShowcase/src/models"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func index(c *gin.Context) {
	var html, err = os.ReadFile("./pages/index.html")
	if err != nil {
		log.Fatal(err)
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	c.Data(http.StatusOK, "text/html; charset=utf-8", html)
}

func main() {
	r := gin.Default()

	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := src.Migrate()
	if err != nil {
		log.Fatal(err)
	}

	models.RecipeRouter(&r.RouterGroup, db)

	r.LoadHTMLGlob("templates/**/*")
	r.Static("/static", "./static")

	r.GET("/", index)
	r.Run(":8000")
}
