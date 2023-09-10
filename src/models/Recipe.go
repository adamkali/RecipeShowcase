package models

import (
	"bytes"
	"fmt"
	"html/template"
	"sync"

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
}

type SurrealRecipe struct {
    ID           string         `json:"id"`
    Name         string         `json:"name"`
    Ingredients  []string       `json:"ingredients"`
    Instructions []RecipeStep   `json:"instructions"`
    PictureURL   string         `json:"picture_url"`
    Tags         []string       `json:"tags"`
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
    }

    fmt.Printf("%d\n", len(sr.Ingredients))
    fmt.Printf("%d\n", len(sr.Tags))
    ingredients := make(chan []Ingredient, len(sr.Ingredients))
    tags := make(chan []RecipeTag, len(sr.Tags))
    convert_error := make(chan error, 2)

    var wg sync.WaitGroup
    var errorWg sync.WaitGroup

    wg.Add(2)
    errorWg.Add(2)

    go func() {
        println("starting ingredients goroutine")
        defer wg.Done()
        var ingredientList []Ingredient
        for _, id := range sr.Ingredients {
            var ingredient Ingredient
            println(id)
            i, err := db.Select(id)
            if err != nil {
                println(err.Error())
                convert_error <- err
                return // Exit the goroutine if an error occurs.
            }
            err = surrealdb.Unmarshal(i, &ingredient)
            if err != nil {
                println(err.Error())
                convert_error <- err
                return // Exit the goroutine if an error occurs.
            }
            ingredientList = append(ingredientList, ingredient)
        }
        fmt.Printf("%v\n", ingredientList)
        ingredients <- ingredientList
        errorWg.Done()
        println("finished ingredients goroutine")
    }()
    
    go func() {
        defer wg.Done()
        var tagList []RecipeTag
        for _, id := range sr.Tags {
            var tag RecipeTag
            i, err := db.Select(id)
            if err != nil {
                convert_error <- err
                return // Exit the goroutine if an error occurs.
            }
            err = surrealdb.Unmarshal(i, &tag)
            if err != nil {
                convert_error <- err
                return // Exit the goroutine if an error occurs.
            }
            tagList = append(tagList, tag)
        }
        tags <- tagList
        errorWg.Done()
        println("finished tags goroutine")
    }()

    // Wait for both goroutines to complete.
    wg.Wait()

    // Close channels after all goroutines have finished.
    close(ingredients)
    close(tags)

    // Wait for error handling goroutines to complete.
    errorWg.Wait()

    // Check if there were any errors.
    select {
    case err := <-convert_error:
        if err != nil {
            return r, err
        }
    default:
    }

    r.Ingredients = <-ingredients
    r.Tags = <-tags

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
    // get the recipes from the database
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


    var results Recipes
    resultChan := make(chan Recipe, len(objs)) // Create a channel for results

    var wg sync.WaitGroup

    for _, i := range objs {
        wg.Add(1) // Increment the WaitGroup counter for each goroutine.

        go func(sr SurrealRecipe) {
            defer wg.Done() // Decrement the WaitGroup counter when the goroutine completes
            r, err := sr.Convert(rc.db)
            if err != nil {
                fmt.Printf(err.Error())
                src.Err(err, c)
                return
            }
            // Send the result to the channel.
            resultChan <- r
        }(i)
    }

    // Create a goroutine to close the resultChan when all goroutines are done.
    go func() {
        wg.Wait()
        close(resultChan)
    }()

    // Collect results from the channel.
    for r := range resultChan {
        results.RecipeList = append(results.RecipeList, r)
    }

    if err != nil {
        fmt.Printf(err.Error())
        src.Err(err, c)
        return
    }

    src.OK("recipe/recipe-preview.tmpl", c, results)
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
