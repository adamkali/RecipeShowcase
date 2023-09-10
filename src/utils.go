package src

import (
	"fmt"
	"html/template"
	"os"
	"strings"

	"github.com/surrealdb/surrealdb.go"
)

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

func RenderError(errorMessage string) string {
    buf, err := template.New("error").Parse("templates/error_template.tmpl")
    if err != nil {
        return "Oops!"
    }

    var output strings.Builder
    err = buf.Execute(&output, struct{ ErrorMessage string }{errorMessage})
    if err != nil {
        return "Oops"
    }

    return output.String()
}
