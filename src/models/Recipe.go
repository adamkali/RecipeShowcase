package models

import (
	"bytes"
	"fmt"
	"html/template"
	"net/http"

	"github.com/adamkali/RecipeShowcase/src"
	"github.com/gin-gonic/gin"
	"github.com/surrealdb/surrealdb.go"
)

type RecipeController struct {
	db *surrealdb.DB
}

// Ingredient represents the ingredients table.
type Ingredient struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Recipe represents the recipe table.
type Recipe struct {
	ID           string       `json:"id,omitempty"`
	Name         string       `json:"name"`
	Ingredients  []Ingredient `json:"ingredients"`
	Instructions []RecipeStep `json:"instructions"`
	PictureURL   string       `json:"picture_url"`
	Tags         []RecipeTag  `json:"tags"`
	Active       bool         `json:"active"`
}

type SurrealRecipe struct {
	ID           string       `json:"id,omitempty"`
	Name         string       `json:"name"`
	Ingredients  []string     `json:"ingredients"`
	Instructions []RecipeStep `json:"instructions"`
	PictureURL   string       `json:"picture_url"`
	Tags         []string     `json:"tags"`
	Active       bool         `json:"active"`
}

// RecipeStep represents the recipe_step table.
type RecipeStep struct {
	ID           string `json:"id"`
	Step         string `json:"step"`
	Type         string `json:"type"`
	ImageFileLoc string `json:"image_file_loc"`
}

// RecipeTag represents the recipe_tag table.
type RecipeTag struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Recipes struct{ RecipeList []Recipe }

// for SurrealRecipe get the Ingredients, Instructions, and Tags
// using the ids from surreal db and convert to the appropriate structs
// concurrently

func (sr SurrealRecipe) Convert(db *surrealdb.DB) (Recipe, error) {
	r := Recipe{
		ID:           sr.ID,
		Name:         sr.Name,
		PictureURL:   sr.PictureURL,
		Instructions: sr.Instructions,
		Ingredients:  []Ingredient{},
		Tags:         []RecipeTag{},
		Active:       sr.Active,
	}

	for _, id := range sr.Ingredients {
		var returnIngredient Ingredient
		println(id)
		ingredient, err := db.Select(id)
		if err != nil {
			return r, fmt.Errorf("Could not find: %s", id)
		}
		err = surrealdb.Unmarshal(ingredient, &returnIngredient)
		if err != nil {
			return r, err
		}
		r.Ingredients = append(r.Ingredients, returnIngredient)
	}

	for _, id := range sr.Tags {
		var returnRecipeTag RecipeTag
		tag, err := db.Select(id)
		if err != nil {
			return r, err
		}
		err = surrealdb.Unmarshal(tag, &returnRecipeTag)
		if err != nil {
			return r, err
		}
		r.Tags = append(r.Tags, returnRecipeTag)
	}

	return r, nil
}

// load the template from the file system
// and execute it with the recipe data
func (r Recipe) Render() (string, error) {
	out := ""

	buf, err := template.ParseFiles("../../templates/recipe/recipe.tmpl")
	if err != nil {
		return out, err
	}

	// execute the template
	buf.Execute(bytes.NewBufferString(out), r)

	return out, nil
}

// load the preview template from the file system
// and execute it with the recipe data
func (r Recipe) RenderPreview() (string, error) {
	out := ""

	buf, err := template.ParseFiles("../../templates/recipe/recipe_preview.tmpl")
	if err != nil {
		return out, err
	}

	// execute the template
	buf.Execute(bytes.NewBufferString(out), r)

	return out, nil
}

func (rc RecipeController) GetRecipies(c *gin.Context) {
	obj, err := rc.db.Select("recipe")
	if err != nil {
		fmt.Printf(err.Error())
		return
	}

	var objs []SurrealRecipe
	err = surrealdb.Unmarshal(obj, &objs)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	var results Recipes = Recipes{RecipeList: []Recipe{}}
	for _, i := range objs {
		r, err := i.Convert(rc.db)
		if err != nil {
			fmt.Printf("%s\n", err.Error())
			src.Err(err, c)
			return
		}
		if r.Active {
			results.RecipeList = append(results.RecipeList, r)
		}
	}

	src.OK("recipe/recipe-preview.tmpl", c, results)
}

func (rc RecipeController) GetRecipeForm(c *gin.Context) {
	surrealRecipe := SurrealRecipe{
		Name:         "New Recipe",
		Ingredients:  []string{},
		Instructions: []RecipeStep{},
		Tags:         []string{},
		Active:       false,
		PictureURL:   "",
	}

	object, err := rc.db.Create("recipe", surrealRecipe)
	if err != nil {
		fmt.Printf("[ERROR]\tGetRecipeForm::Creating =>\t%s\n", err.Error())
		src.Err(err, c)
		return
	}
    
    var empty []SurrealRecipe
    err = surrealdb.Unmarshal(object, &empty)
    if err != nil {
		fmt.Printf("[ERROR]\tGetRecipeForm::Unmarshal =>\t%s\n", err.Error())
		src.Err(err, c)
		return
    }

	recipe, err := empty[0].Convert(rc.db)
	if err != nil {
		fmt.Printf("[ERROR]\tGetRecipeForm::Converting =>\t%s\n", err.Error())
		src.Err(err, c)
		return
	}

	src.OK("recipe/recipe-form.tmpl", c, recipe)
}

func (rc RecipeController) PutRecipeName(c *gin.Context) {
	recipeID := c.Param("id")

	recipe, err := rc.db.Select(recipeID)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	var recipeSurreal SurrealRecipe
	err = surrealdb.Unmarshal(recipe, &recipeSurreal)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	recipeSurreal.Name = c.PostForm("name")
	recipe, err = rc.db.Update(recipeID, recipeSurreal)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	var recipeStruct Recipe
	err = surrealdb.Unmarshal(recipe, &recipeStruct)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	src.OK("recipe/recipe-form-name.tmpl", c, recipeStruct)
}

func (rc RecipeController) PutRecipeIngredient(c *gin.Context) {

	recipeID := c.Param("id")
	recipe, err := rc.db.Select(recipeID)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	var recipeSurreal SurrealRecipe
	err = surrealdb.Unmarshal(recipe, &recipeSurreal)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	ingredient := Ingredient{
		Name:        c.PostForm("name"),
		Description: c.PostForm("description"),
	}
	ingredientID, err := rc.db.Create("ingredient", ingredient)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	str := fmt.Sprintf("%s", ingredientID)
	recipeSurreal.Ingredients = append(recipeSurreal.Ingredients, str)
	recipe, err = rc.db.Update(recipeID, recipeSurreal)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	var recipeStruct Recipe
	err = surrealdb.Unmarshal(recipe, &recipeStruct)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	src.OK("recipe/recipe-form-ingredient.tmpl", c, recipeStruct)
}

func (rc RecipeController) PutRecipeInstruction(c *gin.Context) {

	recipeID := c.Param("id")
	recipe, err := rc.db.Select(recipeID)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	var recipeSurreal SurrealRecipe
	err = surrealdb.Unmarshal(recipe, &recipeSurreal)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	instruction := RecipeStep{
		Step:         c.PostForm("step"),
		Type:         c.PostForm("type"),
		ImageFileLoc: c.PostForm(""),
	}
	recipeSurreal.Instructions = append(recipeSurreal.Instructions, instruction)
	recipe, err = rc.db.Update(recipeID, recipeSurreal)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	var recipeStruct Recipe
	err = surrealdb.Unmarshal(recipe, &recipeStruct)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	src.OK("recipe/recipe-form-instruction.tmpl", c, recipeStruct)
}

func (rc RecipeController) PutRecipeTag(c *gin.Context) {

	recipeID := c.Param("id")
	recipe, err := rc.db.Select(recipeID)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	var recipeSurreal SurrealRecipe
	err = surrealdb.Unmarshal(recipe, &recipeSurreal)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	tag := RecipeTag{
		Name: c.PostForm("name"),
	}
	tagID, err := rc.db.Create("tag", tag)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	str := fmt.Sprintf("%s", tagID)
	recipeSurreal.Tags = append(recipeSurreal.Tags, str)
	recipe, err = rc.db.Update(recipeID, recipeSurreal)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	var recipeStruct Recipe
	err = surrealdb.Unmarshal(recipe, &recipeStruct)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	src.OK("recipe/recipe-form-tag.tmpl", c, recipeStruct)
}

func (rc RecipeController) SaveRecipe(c *gin.Context) {
	// switch the recipe to Active
	recipeID := c.Param("id")
	recipe, err := rc.db.Select(recipeID)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	var recipeSurreal SurrealRecipe
	err = surrealdb.Unmarshal(recipe, &recipeSurreal)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	recipeSurreal.Active = true
	recipe, err = rc.db.Update(recipeID, recipeSurreal)
	if err != nil {
		fmt.Printf(err.Error())
		src.Err(err, c)
		return
	}

	// redirect to the to index
	c.Redirect(http.StatusMovedPermanently, "/")
}

func RecipeRouter(router *gin.RouterGroup, db *surrealdb.DB) {
	rc := RecipeController{
		db: db,
	}
	recipes := router.Group("/recipes")
	{
		recipes.GET("", rc.GetRecipies)
		recipes.GET("/new", rc.GetRecipeForm)
	}
	recipe := router.Group("/recipe")
	{
		recipe.PUT("/:id/name", rc.PutRecipeName)
		recipe.PUT("/:id/ingredient", rc.PutRecipeIngredient)
		recipe.PUT("/:id/instruction", rc.PutRecipeInstruction)
		recipe.PUT("/:id/tag", rc.PutRecipeTag)
		recipe.POST("/:id", rc.SaveRecipe)
	}
}
