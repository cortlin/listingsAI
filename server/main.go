package main

import (
	"fmt"

	"github.com/cortlin/mls-ai/db"
	"github.com/graphql-go/graphql"
)

type stuff = db.Listing

var rootQuery = graphql.NewObject((graphql.ObjectConfig{
	Name:   "RootQuery",
	Fields: graphql.Fields{},
}))

func main() {
	fmt.Println("Hello World")
}
