package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
)

var markdown *bool
var cli *bool
var dbPath *string
var add *string
var bulkload *string

func getOptions() {
	dbPath = flag.String("path", os.Getenv("RECIPE"), "path to recipe <RECIPE> database")
	markdown = flag.Bool("md", false, "outputs recipe as MarkDown into stdout")
	cli = flag.Bool("cli", false, "outputs recipe as string into stdout")
	//adding = flag.Bool("add-xml", false, "adding an recipe via a pipe!")
	add = flag.String("add", "", "path to recipe file")
	bulkload = flag.String("load", "", "path to directory with xml/json files")
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

	var db = boltdb{path: *dbPath}
	var se = Index{}

	db.Init()
	se.Init(&db)

	if len(*bulkload) > 0 {
		loadBulkFiles(db, se, filepath.ToSlash(*bulkload))
	} else if len(*add) > 0 {
		recipe := loadSingleFile(*add)
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

func getFileContent(path string) []byte{
	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		handleError(err)
	}
	return content
}

func loadSingleFile(file string) Recipe{
	var recipe Recipe
	ext := path.Ext(file)
	content := getFileContent(file)
	if strings.ToLower(ext) == ".xml" {
		recipe = UnmarschalXMLRecipe(content)
	} else if strings.ToLower(ext) == ".json" {
		recipe = UnmarschalJSONRecipe(content)
	} else {
		fmt.Println("not a supported file format (extension) of ", file)
	}
	return recipe
}

func loadBulkFiles(db boltdb, se Index, dir string){
	entries, err := ioutil.ReadDir(dir)
	if err != nil{
		log.Fatal(err)
	}
	for _, file := range entries{
		recipe := loadSingleFile(dir + file.Name())
		db.Add(recipe)
		se.Index(recipe)
	}
}