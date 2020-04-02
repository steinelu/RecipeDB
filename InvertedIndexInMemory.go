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

func (self *InvertedIndexInMemory) Index(base DataBase) {
	self.db = base
	self.index = make(map[string][]int)
	self.idToTitle = make(map[int]string)

	for recipe := range self.db.Iterator() {
		rid := rand.Int()
		self.add(parseRecipe(recipe), rid)
		self.idToTitle[rid] = recipe.Filename()
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

func (self *InvertedIndexInMemory) add(terms []string, id int) {
	for _, term := range terms {
		self.index[term] = append(self.index[term], id)
	}
}

func (self *InvertedIndexInMemory) Search(terms []string) <-chan Recipe {
	// boolean retrival
	res := self.index[terms[0]]

	for _, term := range terms[0:] {
		res = intersect(res, self.index[term])
	}

	recipes := []string{}
	for _, id := range res {
		recipes = append(recipes, self.idToTitle[id])
	}
	//fmt.Println(recipes)
	return self.db.Get(recipes)
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

func max(x int, y int) int {
	if x > y {
		return x
	}
	return y
}
