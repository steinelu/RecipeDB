package main

import (
	"strings"
)

type InvertedIndexInMemory struct {
	index map[string][]string
	//idToTitle map[int]string
	db DataBase
}

func (iiMem *InvertedIndexInMemory) Index(base DataBase) {
	iiMem.db = base
	iiMem.index = make(map[string][]string)

	for recipe := range iiMem.db.Iterator() {
		iiMem.add(parseRecipe(recipe), recipe.GetId())
	}

	//iiMem.index = iiMem.indexBuild.
}

func parseRecipe(recipe Recipe) []string {
	var tokens []string
	for _, tok := range strings.Fields(recipe.Title) {
		tokens = append(tokens, tokenize(tok))
	}
	return tokens
}

func tokenize(token string) string {
	return strings.ToLower(token)
}

func (iiMem *InvertedIndexInMemory) add(terms []string, recipe_hash string) {
	for _, term := range terms {
		//TODO insert into slice as heap, instead of append to list
		iiMem.index[term] = append(iiMem.index[term], recipe_hash)
	}
}

func (iiMem *InvertedIndexInMemory) Search(terms []string) <-chan Recipe {
	// boolean retrival
	//res := iiMem.index[terms[0]]

	//for _, term := range terms[0:] {
	//	res = intersect(res, iiMem.index[term])
	//}
	res := []string{}
	for _, term := range terms {
		res = append(res, iiMem.index[term]...)
	}
	var recipes []string

	for _, hash := range res {
		recipes = append(recipes, hash+".xml")
	}
	//fmt.Println(recipes)
	return iiMem.db.Get(recipes)
}

func intersect(one []int, two []int) []int {
	var intersection []int
	for i, j := 0, 0; i < len(one) && j < len(two); {
		if one[i] < two[j] {
			i++
		} else if one[i] > two[j] {
			j++
		} else {
			intersection = append(intersection, one[i])
			i++
			j++
		}
	}
	return intersection
}
