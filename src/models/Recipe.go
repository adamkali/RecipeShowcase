package models

import (
	"bytes"
	"fmt"
	"html/template"

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
    ID           string        `json:"id"`
    Name         string        `json:"name"`
    Ingredients  []Ingredient  `json:"ingredients"`
    Instructions []RecipeStep  `json:"instructions"`
    PictureURL   string        `json:"picture_url"`
    Tags         []RecipeTag   `json:"tags"`
    Active       bool          `json:"active"`
}

type SurrealRecipe struct {
    ID           string         `json:"id"`
    Name         string         `json:"name"`
    Ingredients  []string       `json:"ingredients"`
    Instructions []RecipeStep   `json:"instructions"`
    PictureURL   string         `json:"picture_url"`
    Tags         []string       `json:"tags"`
    Active       bool           `json:"active"`
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

type Recipes struct { RecipeList []Recipe }

// for SurrealRecipe get the Ingredients, Instructions, and Tags 
// using the ids from surreal db and convert to the appropriate structs
// concurrently

func (sr SurrealRecipe) Convert(db *surrealdb.DB) (Recipe, error) {
    r := Recipe{
        ID:         sr.ID,
        Name:       sr.Name,
        PictureURL: sr.PictureURL,
        Instructions: sr.Instructions,
        Ingredients: []Ingredient{},
        Tags: []RecipeTag{},
        Active: sr.Active,
    }

    for _, id := range(sr.Ingredients) {
        var returnIngredient Ingredient
        println(id)
        ingredient, err := db.Select(id)
        if err  != nil {
            return r, fmt.Errorf("Could not find: %s", id)
        }
        err = surrealdb.Unmarshal(ingredient, &returnIngredient)
        if err  != nil {
            return r, err
        }
        r.Ingredients = append(r.Ingredients, returnIngredient)
    }
    
    for _, id := range(sr.Tags) {
        var returnRecipeTag RecipeTag
        tag, err := db.Select(id)
        if err  != nil {
            return r, err
        }
        err = surrealdb.Unmarshal(tag, &returnRecipeTag)
        if err  != nil {
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


    var results Recipes = Recipes{ RecipeList: []Recipe{}, }

    for _, i := range objs {
        r, err := i.Convert(rc.db)
        if err != nil {
            fmt.Printf("%s\n",err.Error())
            src.Err(err, c)
            return
        }
        if r.Active {
            results.RecipeList = append(results.RecipeList, r)
        }
    }

    src.OK("recipe/recipe-preview.tmpl", c, results)
}

func (rc RecipeController) GetRecipe(c *gin.Context) {
    recipe := Recipe{
        Name: "New Recipe",
        Ingredients: []Ingredient{},
        Instructions: []RecipeStep{},
        Tags: []RecipeTag{},
        Active: false,
        PictureURL: "",
        
    }
    object, err := rc.db.Create("recipe", recipe)
    if err != nil {
        fmt.Printf(err.Error())
        src.Err(err, c)
        return
    }
    var recipeSurreal SurrealRecipe
    err = surrealdb.Unmarshal(object, &recipeSurreal)
    if err != nil {
        fmt.Printf(err.Error())
        src.Err(err, c)
        return
    }

    recipe, err = recipeSurreal.Convert(rc.db)
    if err != nil {
        fmt.Printf(err.Error())
        src.Err(err, c)
        return
    }

    src.OK("recipe/recipe-form.tmpl", c, recipe)
}

func RecipeRouter(router *gin.RouterGroup, db *surrealdb.DB) {
    rc := RecipeController {
        db: db,
    }
    recipe := router.Group("/recipes")
    {
        recipe.GET("", rc.GetRecipies)
    }
}
