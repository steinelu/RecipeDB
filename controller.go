package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

type Options struct {
	path     string
	markdown bool
	cli      bool
	add      bool
	tail     []string
}

var arguments Options

func getOptions() {
	arguments.path = *flag.String("path", os.Getenv("RECIPE"), "path to recipe <RECIPE> database")
	arguments.markdown = *flag.Bool("md", false, "outputs recipe as MarkDown into stdout")
	arguments.cli = *flag.Bool("cli", false, "outputs recipe as string into stdout")
	arguments.add = *flag.Bool("add-xml", false, "adding an recipe via a pipe!")
	arguments.tail = flag.Args()
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

	fmt.Println(arguments.path, len(arguments.path))
	var db = XMLLazy{path: arguments.path}
	var se = InvertedIndexInMemory{}
	db.Init()
	se.Index(&db)

	if len(arguments.tail) > 0 { // searching
		for recipe := range se.Search(arguments.tail) {
			if arguments.markdown {
				fmt.Println(recipe.toMarkdown())
			} else if arguments.cli {
				fmt.Println(recipe.toCLIString())
			} else {
				fmt.Println(recipe.Title)
			}
		}
	}

	if arguments.add {
		x := shellPipeInput()
		recipe := db.ParseXMLContent([]byte(x))
		fmt.Println(recipe)
		db.Add(recipe)
	}
}

func shellPipeInput() string {
	info, err := os.Stdin.Stat()
	if err != nil {
		log.Fatal(err)
	}
	if info.Mode()&os.ModeCharDevice != 0 || info.Size() <= 0 {
		fmt.Println("Pipes!")
		os.Exit(1)
	}
	reader := bufio.NewReader(os.Stdin)
	var output []rune
	for {
		input, _, err := reader.ReadRune()
		if err != nil && err == io.EOF {
			break
		}
		output = append(output, input)
	}
	return string(output)
}

//CreateData
/*func CreateData(db DataBase) {
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
*/
