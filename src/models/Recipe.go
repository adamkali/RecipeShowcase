package models

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/surrealdb/surrealdb.go"
	"github.com/adamkali/RecipeShowcase/src"
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

type Recipes struct { Recipes []Recipe }

// for SurrealRecipe get the Ingredients, Instructions, and Tags 
// using the ids from surreal db and convert to the appropriate structs
// concurrently
func (sr SurrealRecipe) Convert(db *surrealdb.DB) (Recipe, error) {
    // TODO

    println("Starting Recipe Convert")
    r := Recipe{
        ID: sr.ID,
        Name: sr.Name,
        PictureURL: sr.PictureURL,
    }

    ingredients := make(chan []Ingredient, len(sr.Ingredients))
    tags := make(chan []RecipeTag, len(sr.Tags))
    convert_error := make(chan error, 1)

    var wg sync.WaitGroup

    go func() {
        defer wg.Done()
        for _, id := range sr.Ingredients {
            var ingredient Ingredient
            println(id)
            i, err := db.Select(id)
            if err != nil {
                convert_error <- err
            }
            err = surrealdb.Unmarshal(i,&ingredient)
            fmt.Printf("Ingredient: %v\n", i)
            if err != nil {

                convert_error <- err
            }
            ingredients <- []Ingredient{ingredient}
        }
        close(ingredients)
        convert_error <- nil
    }()
    
    go func() {
        defer wg.Done()
        for _, id := range sr.Tags {
            var tag RecipeTag
            println(id)
            i, err := db.Select(id)
            if err != nil {
                convert_error <- err
            }
            fmt.Printf("RecipeTag: %v\n", i)
            err = surrealdb.Unmarshal(i,&tag)
            if err != nil {
                convert_error <- err 
            }
            tags <- []RecipeTag{tag}
        }
        close(tags)
        convert_error <- nil
    }()

    wg.Add(2)

    if err := <-convert_error; err != nil {
        return r, err
    }

    return r, nil
}

// load the template from the file system
// and execute it with the recipe data
func (r Recipe) Render() (string, error) {
    out := ""
    
    buf, err := template.ParseFiles("templates/recipe/recipe.tmpl")
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
    
    buf, err := template.ParseFiles("templates/recipe/recipe_preview.tmpl")
    if err != nil {
        return out, err
    }

    // execute the template 
    buf.Execute(bytes.NewBufferString(out), r) 

    return out, nil
}

func (rs Recipes) Render() (string, error) {
    var out strings.Builder

    // Start the outer <div> element with Tailwind CSS classes for column layout
    out.WriteString("<div class=\"grid grid-cols-1 sm:grid-cols-2 md:grid-cols-3 lg:grid-cols-4 gap-4\">\n\t")

    for _, recipe := range rs.Recipes {
        str, err := recipe.RenderPreview()
        if err != nil {
            return "", err
        }
        out.WriteString(str)
        out.WriteString("\n")
    }

    // Close the outer <div> element
    out.WriteString("</div>")

    return out.String(), nil
}

func (rc RecipeController) GetRecipies(c *gin.Context) {
    // get the recipes from the database
    obj, err := rc.db.Select("recipe")
    if err != nil {
        // TODO: change this to use an error page
        c.Data(500,"text/html; charset=utf-8", []byte(src.RenderError(err.Error())))
        return
    }


    var objs []SurrealRecipe
    err = surrealdb.Unmarshal(obj, &objs)
    if err != nil {
        fmt.Printf("Could Not Unmarshal Data: %s\n", err.Error())
        c.Data(500,"text/html; charset=utf-8", []byte(src.RenderError(err.Error())))
        return
    }

    fmt.Printf("Converted Recipe: %v\n", objs) // Add this line for debugging
    var results Recipes
    for _, i := range objs {
        fmt.Printf("Recipe: %v\n", i)
        r, err := i.Convert(rc.db)
        if err != nil {
            fmt.Printf("Could Not Convert Data: %s\n", err.Error())
            c.Data(500,"text/html; charset=utf-8", []byte(src.RenderError(err.Error())))
            return
        }
        fmt.Printf("Converted Recipe: %v\n", r) // Add this line for debugging
        results.Recipes = append(results.Recipes, r)
    }

    html_body, err := results.Render()
    if err != nil {
        c.Data(500,"text/html; charset=utf-8", []byte(src.RenderError(err.Error())))
        return
    }

    // return html_body as the html_body
    c.Data(200, "text/html; charset=utf-8", []byte(html_body))
}

func RecipeRouter(router *gin.RouterGroup, db *surrealdb.DB) {
    rc := RecipeController {
        db: db,
    }
    recipe := router.Group("/recipes")
    {
        recipe.GET("/", rc.GetRecipies)
    }
}
