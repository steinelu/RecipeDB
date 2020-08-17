package main

import (
	"container/heap"
	"strings"
)

type InvertedIndexInMemory struct {
	index_ map[string]*StringHeap
	db     DataBase
}

func (iiMem *InvertedIndexInMemory) Index(base DataBase) {
	iiMem.db = base
	iiMem.index_ = make(map[string]*StringHeap)

	for recipe := range iiMem.db.Iterator() {
		iiMem.add(parseRecipe(recipe), recipe.GetId())
	}
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

func tokenize(token string) string {
	return strings.ToLower(token)
}

func (iiMem *InvertedIndexInMemory) add(terms []string, recipe_hash string) {
	for _, term := range terms {
		if iiMem.index_[term] == nil {
			iiMem.index_[term] = &StringHeap{}
			heap.Init(iiMem.index_[term])
		}

		heap.Push(iiMem.index_[term], recipe_hash)
	}
}

func unique(list []string) []string {
	uniq := []string{list[0]}
	for _, elem := range list {
		if elem != uniq[len(uniq)-1] {
			uniq = append(uniq, elem)
		}
	}
	return uniq
}

func (iiMem *InvertedIndexInMemory) Search(terms []string) <-chan Recipe {
	// boolean retrival
	res := *iiMem.index_[terms[0]]

	for _, term := range terms[0:] {
		res = *intersect(res, *iiMem.index_[term])
	}
	//fmt.Println(res)
	var recipes []string
	for _, hash := range unique(res) {
		recipes = append(recipes, hash+".xml")
	}

	return iiMem.db.Get(recipes)
}

func intersect(one []string, two []string) *StringHeap {
	intersection := StringHeap{}
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
	return &intersection
}

type StringHeap []string

func (h StringHeap) Len() int           { return len(h) }
func (h StringHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h StringHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *StringHeap) Push(x interface{}) {
	*h = append(*h, x.(string))
}

func (h *StringHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}
