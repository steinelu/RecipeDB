package main

import (
	"fmt"
	"log"
	"runtime/debug"
)

func handleError(err error) {
	log.Println(err)
	debug.PrintStack()
	log.Fatal()
}

type SearchEngine interface {
	Index(DataBase)
	Search([]string) <-chan Recipe
}

type DataBase interface {
	Init()
	Add(Recipe)
	Iterator() <-chan Recipe
	Get([]string) <-chan Recipe
}

type Recipe struct {
	filename    string
	id          int
	Title       string       `xml:"title"`
	Ingredients []Ingredient `xml:"ingredients>ingredient"`
	Preparation []string     `xml:"preparation>step"` // TODO saving order of steps
}

type Ingredient struct {
	Name   string  `xml:",chardata"`
	Amount float32 `xml:"amount,attr"`
	Unit   string  `xml:"unit,attr"`
}

func (self Recipe) toMarkdown() string {
	prep := ""
	for _, step := range self.Preparation {
		prep = prep + fmt.Sprintf("1. %s\n", step)
	}

	ingred := ""
	for _, ingredient := range self.Ingredients {
		ingred = ingred + fmt.Sprintf("* %d %s %s\n", ingredient.Amount, ingredient.Unit, ingredient.Name)
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
		ingred = ingred + fmt.Sprintf(" - %d %s %s\n", ingredient.Amount, ingredient.Unit, ingredient.Name)
	}

	return fmt.Sprintf("%s\nPreparation:\n%s\nIngredients:\n%s\n", self.Title, prep, ingred)
}

func (self Recipe) Filename() string {
	return self.filename
}

func (self *Recipe) SetFilename(name string) {
	self.filename = name
}
