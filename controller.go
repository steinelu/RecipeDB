package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

//var path = flag.String("path", os.Getenv("RECIPE"), "path to recipe <RECIPE> database")
//var n = flag.Int("n", 5, "number of results shown")
//var t = flag.Int("t", 0, "only results show with at least priority / start of t, default is 0")
//search

func main() {
	if err := os.Setenv("RECIPE", "./test/"); err != nil {
		handleError(err)
	}
	var path = flag.String("path", os.Getenv("RECIPE"), "path to recipe <RECIPE> database")
	var markdown = flag.Bool("md", false, "outputs recipe as MarkDown into stdout")
	var cli = flag.Bool("cli", false, "outputs recipe as string into stdout")
	var add = flag.Bool("add-xml", false, "adding an recipe via a pipe!")

	flag.Parse()

	var db = XMLLazy{path: *path}
	var se = InvertedIndexInMemory{}
	db.Init()
	se.Index(&db)
	tail := flag.Args()

	if len(tail) > 0 { // searching
		for recipe := range se.Search(tail) {
			if *markdown {
				fmt.Println(recipe.toMarkdown())
			} else if *cli {
				fmt.Println(recipe.toCLIString())
			}
		}
	}
	if *add {
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
