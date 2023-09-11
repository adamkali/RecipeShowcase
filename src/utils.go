package src

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/surrealdb/surrealdb.go"
)

type ErrorType struct { ErrorMessage string }

func Err(what error, c *gin.Context) {
    c.HTML(500, "utils/error_template.tmpl", ErrorType{
        ErrorMessage: what.Error(),
    })
}

func OK(templ string, c *gin.Context, what interface{}) {
    c.HTML(200, templ, what) 
}

func Migrate() (*surrealdb.DB, error) {
    var db *surrealdb.DB

    user := os.Getenv("DB_USER")
    pass := os.Getenv("DB_PASS")
    name := os.Getenv("DB_NAME")
    data := os.Getenv("DB_DATA")

    db, err := surrealdb.New("ws://localhost:8100/rpc")
    if err != nil {
        return db, err
    }

    println("Connecting to database...")
    fmt.Printf("User: %s\n", user)
    fmt.Printf("Pass: %s\n", pass)
    fmt.Printf("Namespace: %s\n", name)
    fmt.Printf("Database: %s\n", data)

    if _, err = db.Signin(map[string]interface{}{
        "user": user,
        "pass": pass,
    }); err != nil {
        return db, err
    }

    if _, err = db.Use(name, data); err != nil {
        return db, err
	}

    // load the migraten from the migrations folder
    var query string
    bs, err := os.ReadFile("recipe_db/migrations.surql")
    if err != nil {
        return db, err
    }

    query = string(bs)
    if _, err = db.Query(query, nil); err != nil {
        return db, err
    }

    return db, nil
}
