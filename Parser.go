package main

import (
	"encoding/json"
	"encoding/xml"
	"strings"
)

func tokenize(token string) string {
	return strings.ToLower(token)
}

func parseRecipe(recipe Recipe) []string {
	var tokens []string
	for _, tok := range strings.Fields(recipe.Title) {
		tokens = append(tokens, tokenize(tok))
	}

	for _, ingredient := range recipe.Ingredients {
		for _, tok := range strings.Fields(ingredient.Name) {
			tokens = append(tokens, tokenize(tok))
		}
	}
	return tokens
}

func UnmarschalXMLRecipe(content []byte) Recipe {
	recipe := Recipe{}
	if err := xml.Unmarshal(content, &recipe); err != nil {
		handleError(err)
	}
	return recipe
}

func MarschalXMLRecipe(recipe Recipe) []byte{
	content, err := xml.MarshalIndent(recipe, "", "	")
	if err != nil {
		handleError(err)
	}
	return content
}

func UnmarschalJSONRecipe(content []byte) Recipe {
	recipe := Recipe{}
	if err := json.Unmarshal(content, &recipe); err != nil {
		handleError(err)
	}
	return recipe
}

func MarschalJSONRecipe(recipe Recipe) []byte{
	content, err := json.MarshalIndent(recipe, "", "	")
	if err != nil {
		handleError(err)
	}
	return content
}