# RecipeDB
small recipe database mostly in golang

How to use:
- Add a recipe (in a XML format) `go run . -add "/path/to/recipe.xml"`
- Add a recipe (in a json format)`go run . -add "/path/to/recipe.json"`
- Add a bunch of recipes located in one folder `go run . -add "/path/to/recipes/"` nutil now the last slash is needed
- Search for a recipe `go run . term1 term2` -> show the title of the recipe containing these words
- Search for a recipe and show it in a specific format `go run . -md term1 term2 term3 ...` 
