package main

import (
	"net/http"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/handler"
)

var listingForGQL = graphql.NewObject(graphql.ObjectConfig{
	Name: "Listing",
	Fields: graphql.Fields{
		"Zip": &graphql.Field{
			Type: graphql.String,
		},
		"StreetName": &graphql.Field{
			Type: graphql.String,
		},
		"id": &graphql.Field{
			Type: graphql.Int,
		},
		"ListingPrice": &graphql.Field{
			Type: graphql.Int,
		},
		"squareFeet": &graphql.Field{
			Type: graphql.String,
		},
	},
})

var rootQuery = graphql.NewObject((graphql.ObjectConfig{
	Name: "RootQuery",
	Fields: graphql.Fields{
		"listing": &graphql.Field{
			Type:        listingForGQL,
			Description: "Get a Single listing by ID",
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.Int,
				},
			},
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return nil, nil
			},
		},
		"Prompt": &graphql.Field{
			Type:        graphql.String,
			Description: "Prompt Sent to the API. TODO: Return the list of Listings instead",
			Args: graphql.FieldConfigArgument{
				"prompt": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				prompt, isOK := p.Args["prompt"].(string)

				if isOK {
					returnString := "This what your prompt: " + prompt
					return returnString, nil
				}
				return "idk", nil
			},
		},
	},
}))

var Schema, _ = graphql.NewSchema(graphql.SchemaConfig{
	Query: rootQuery,
})

func main() {
	h := handler.New(&handler.Config{
		Schema:     &Schema,
		Pretty:     true,
		GraphiQL:   false,
		Playground: true,
	})

	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", nil)
}
