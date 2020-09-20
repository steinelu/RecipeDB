package main

import (
	"crypto/sha1"
	"fmt"
	"log"
	"math"
	"runtime/debug"
)

func handleError(err error) {
	log.Println(err)
	debug.PrintStack()
	log.Fatal()
}

type SearchEngine interface {
	Init(*DataBase)
	Search([]string) <-chan Recipe
	Index(Recipe)
	Close()
}

type DataBase interface {
	Init()
	Add(Recipe)
	//Iterator(func(Recipe, []byte))
	Get([]string) <-chan Recipe
	Close()
}

type Recipe struct {
	//filename    string
	Source 		string		 `xml:"href,attr" json:"href"`
	Title       string       `xml:"title" json:"title"`
	Ingredients []Ingredient `xml:"ingredients>ingredient" json:ingredients`
	Preparation []string     `xml:"preparation>step" json:"preparation"` // TODO saving order of steps
}

type Ingredient struct {
	Name   string  `xml:",chardata" json:"name"`
	Amount float64 `xml:"amount,attr" json:"amount"`
	Unit   string  `xml:"unit,attr" json:"unit"`
}

func (ingredient Ingredient) getAmountUnit() string {
	if ingredient.Amount == 0 {
		return ""
	}
	if ingredient.Amount-math.Floor(ingredient.Amount) == 0 {
		return fmt.Sprintf("%.0f %s ", ingredient.Amount, ingredient.Unit)
	}
	return fmt.Sprintf("%.3f %s ", ingredient.Amount, ingredient.Unit)
}

func (self Recipe) toMarkdown() string {
	prep := ""
	for _, step := range self.Preparation {
		prep = prep + fmt.Sprintf("1. %s\n", step)
	}
	ingred := ""
	for _, ingredient := range self.Ingredients {
		ingred = ingred + ingredient.getAmountUnit() + ingredient.Name + "\n"
	}
	return fmt.Sprintf("## %s\n### Preparation:\n%s\n\n### Ingredients:\n%s\n", self.Title, prep, ingred)
}

func (self Recipe) toCLIString() string {
	prep := ""
	for i, step := range self.Preparation {
		prep = prep + fmt.Sprintf(" [%d.] %s\n", i+1, step)
	}

	ingred := ""
	for _, ingredient := range self.Ingredients {
		ingred = ingred + fmt.Sprintf(" - %f %s %s\n", ingredient.Amount, ingredient.Unit, ingredient.Name)
	}

	return fmt.Sprintf("%s\nPreparation:\n%s\nIngredients:\n%s\n", self.Title, prep, ingred)
}

func (self *Recipe) GetId() string {
	h := sha1.New()
	hash := fmt.Sprintf("%v", self)
	h.Write([]byte(hash))
	return fmt.Sprintf("%x", h.Sum(nil))
}
