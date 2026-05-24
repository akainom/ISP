package main

import (
	"fmt"
	"lab10/graph"
	"lab10/graph/model"
	"log"
	"net/http"
	"os"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db *gorm.DB

func initDB() {
	var err error
	db, err = gorm.Open(sqlite.Open("./celebrities.db"), &gorm.Config{})
	if err != nil {
		log.Fatal("[FATAL] failed to connect database:", err)
	}

	db.AutoMigrate(&model.Celebrity{})
}

func main() {
	initDB()

	file, _ := os.OpenFile("10.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	log.SetOutput(file)

	srv := handler.NewDefaultServer(
		graph.NewExecutableSchema(
			graph.Config{
				Resolvers: &graph.Resolver{
					DB: db,
				},
			},
		),
	)

	http.Handle("/pg", playground.Handler("playground", "/gql"))
	http.Handle("/gql", srv)

	fmt.Println("server started on localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
