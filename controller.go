package main

import (
	"flag"
	"fmt"
	"os"
)

//var path = flag.String("path", os.Getenv("RECIPE"), "path to recipe <RECIPE> database")
//var n = flag.Int("n", 5, "number of results shown")
//var t = flag.Int("t", 0, "only results show with at least priority / start of t, default is 0")
//search

type options struct {
	path string
}

func main() {
	//TODO http://www.albertoleal.me/posts/golang-pipes.html
	mainSearch()
}

func mainSearch() {

	if err := os.Setenv("RECIPE", "./test/"); err != nil {
		handleError(err)
	}
	var path = flag.String("path", os.Getenv("RECIPE"), "path to recipe <RECIPE> database")
	var markdown = flag.Bool("md", false, "outputs recipe as MarkDown into stdout")
	var cli = flag.Bool("cli", false, "outputs recipe as string into stdout")

	flag.Parse()

	var db = XMLLazy{path: *path}
	//createData(&db)
	var se = InvertedIndexInMemory{}

	db.Init()
	se.Index(&db)

	for recipe := range se.Search([]string{"tomate"}) {
		if *markdown {
			fmt.Println(recipe.toMarkdown())
		} else if *cli {
			fmt.Println(recipe.toCLIString())
		}
	}
}

func createData(db DataBase) {
	db.Add(Recipe{Title: "Kartoffeln und Speck"})
	db.Add(Recipe{Title: "Nudeln mit Sosse und Ei"})
	db.Add(Recipe{Title: "Mehr mit Salz und Ei"})
	db.Add(Recipe{Title: "Nudeln und Tomaten und Ei"})
	db.Add(Recipe{Title: "Pfannenkuchen",
		Preparation: []string{"Eier und Milch zusammen verquirlen.",
			"Mehl in eine Schüssel geben, eier und Milch hinzugeben und verrühren.",
			"Danach geschmolzene Butter hinzugeben, umrühren bis ein weicher Teig entsteht.",
			"Teig eine Stunde ruhen lassen."},
		Ingredients: []Ingredient{{Name: "Eier", Amount: 6},
			{Name: "Mehl", Amount: 400, Unit: "g."},
			{Name: "Milch", Amount: 750, Unit: "ml."},
			{Name: "Butter", Amount: 1, Unit: "El."}}})
}
