package main

import (
	"math/rand"
	"strings"
)

type InvertedIndexInMemory struct {
	index     map[string][]int
	idToTitle map[int]string
	db        DataBase
}

func (iiMem *InvertedIndexInMemory) Index(base DataBase) {
	iiMem.db = base
	iiMem.index = make(map[string][]int)
	iiMem.idToTitle = make(map[int]string)

	for recipe := range iiMem.db.Iterator() {
		rid := rand.Int()
		iiMem.add(parseRecipe(recipe), rid)
		iiMem.idToTitle[rid] = recipe.Filename()
	}
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

func (iiMem *InvertedIndexInMemory) add(terms []string, id int) {
	for _, term := range terms {
		iiMem.index[term] = append(iiMem.index[term], id)
	}
}

func (iiMem *InvertedIndexInMemory) Search(terms []string) <-chan Recipe {
	// boolean retrival
	res := iiMem.index[terms[0]]

	for _, term := range terms[0:] {
		res = intersect(res, iiMem.index[term])
	}

	var recipes []string

	for _, id := range res {
		recipes = append(recipes, iiMem.idToTitle[id])
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
