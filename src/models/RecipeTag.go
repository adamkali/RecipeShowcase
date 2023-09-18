package models

import (
	"fmt"
	"net/http"

	"github.com/adamkali/RecipeShowcase/src"
	"github.com/gin-gonic/gin"
	"github.com/surrealdb/surrealdb.go"
)

type RecipeTagController struct {
	db *surrealdb.DB
}

func (controller *RecipeTagController) GetTags(ctx *gin.Context) {
	// get the query string for recipe_id
	recipe_id := ctx.Query("recipe_id")
	if recipe_id == "" {
		err := fmt.Errorf("No recipe_id provided")
		src.ErrWithCode(err, ctx, http.StatusBadRequest)
	}

	// get the recipe
	recipe_surreal, err := controller.db.Select(recipe_id)
	if err != nil {
		fmt.Println(err)
		src.Err(err, ctx)
	}

	var recipe Recipe
	err = surrealdb.Unmarshal(recipe_surreal, &recipe)
	if err != nil {
		fmt.Println(err)
		src.Err(err, ctx)
	}

	// get the tags
	tags_surreal, err := controller.db.Select("tags")
	if err != nil {
		fmt.Println(err)
		src.Err(err, ctx)
	}

	var tags []RecipeTag
	err = surrealdb.Unmarshal(tags_surreal, &tags)
	if err != nil {
		fmt.Println(err)
		src.Err(err, ctx)
	}

	src.OK("", ctx, gin.H{
		"ID":   recipe.ID,
		"Tags": recipe.Tags,
	})
}

func TagRouter(router *gin.RouterGroup, db *surrealdb.DB) {
	controller := RecipeTagController{db: db}
	recipes := router.Group("/tags")
	{
		recipes.GET("", controller.GetTags)
	}
}
