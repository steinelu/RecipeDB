package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

var markdown *bool
var cli *bool
var path *string
//var adding *bool
var adding *string

func getOptions() {
	path = flag.String("path", os.Getenv("RECIPE"), "path to recipe <RECIPE> database")
	markdown = flag.Bool("md", false, "outputs recipe as MarkDown into stdout")
	cli = flag.Bool("cli", false, "outputs recipe as string into stdout")
	//adding = flag.Bool("add-xml", false, "adding an recipe via a pipe!")
	adding = flag.String("add-recipe-xml", "", "path to xml file of the recipe")
	flag.Parse()
}

func setDefaultPath() {
	if len(os.Getenv("RECIPE")) <= 0 {
		if err := os.Setenv("RECIPE", "./test/"); err != nil {
			handleError(err)
		}
	}
}

func main() {
	setDefaultPath()
	getOptions()

	//fmt.Println(arguments.path, len(arguments.path))
	var db = boltdb{path: *path}
	var se = InvertedIndexInMemory{}

	db.Init()


	if len(flag.Args()) > 0 { // searching
		se.Index(&db)
		res := se.Search(flag.Args())
		handleSearchResults(res)
	} else if len(*adding) > 0 {
		file, err := os.Open(*adding)
		if err != nil {
			log.Fatal(err)
		}
		content, err := ioutil.ReadAll(file)

		if err != nil {
			handleError(err)
		}
		db.Add(UnmarschalXMLRecipe(content))
	}
}

func handleSearchResults(recipes <-chan Recipe) {
	for recipe := range recipes {
		if *markdown {
			fmt.Println(recipe.toMarkdown())
		} else if *cli {
			fmt.Println(recipe.toCLIString())
		} else {
			fmt.Println(recipe.GetId(), ":", recipe.Title)
		}
	}
}