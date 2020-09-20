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

	var db = boltdb{path: *path}
	var se = Index{}

	db.Init()
	se.Init(&db)

	if len(*adding) > 0 {
		recipe := getFileContent(*adding)
		db.Add(recipe)
		se.Index(recipe)
	} else if len(flag.Args()) > 0 { // Searching
		res := se.Search(flag.Args())
		handleSearchResults(res)
	} else {
		fmt.Println("Dry run ...")
	}

	db.Close()
	se.Close()
}

func handleSearchResults(recipes <-chan Recipe) {
	for recipe := range recipes {
		if *markdown {
			fmt.Println(recipe.toMarkdown())
		} else if *cli {
			fmt.Println(recipe.toCLIString())
		} else {
			fmt.Println(recipe.GetId(), ": ", recipe.Title)
		}
	}
}

func getFileContent(path string) Recipe{
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		handleError(err)
	}
	return UnmarschalXMLRecipe(content)
}