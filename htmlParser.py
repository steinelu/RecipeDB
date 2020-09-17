from bs4 import BeautifulSoup
from pprint import pprint
import sys
import subprocess

recipe = {}

def getHTMLFromWebsiteChromium(url):
    cmd = f"chromium --headless --disable-gpu {url}".split(" ")
    return subprocess.run(cmd)
    
def parseLocalHTML(html):
    with open(html) as html_:
        soup = BeautifulSoup(html_.read(), "html.parser")
        
        recipe["title"] = soup.h1.string.strip()
        parsePreparation(soup)
        parseIngredients(soup)

def parsePreparation(html):
    prep = html.find("ol", "recipe-preparation").find_all("p")
    recipe["preparation"] = [pr.string.strip() for pr in prep]

def parseIngredients(html):
    ingred = html.find_all("div", "recipe-ingredients")
    ingred = [ing.find_all("ul") for ing in ingred]
    ingred = [ing.find_all("span") for ing in ingred[0]]
    ingred = ingred[0]

    name = [ing.string for i, ing in enumerate(ingred) if i%2 == 1]
    amount = [ing.string for i, ing in enumerate(ingred) if i%2 == 0]
    recipe["ingredients"] = {n:a.split() if len(a.split()) > 1 else ["", ""] for n, a in zip(name, amount)}

def toXML(recipe):
    ing = "\n".join([f'\t\t<ingredient amount="{vals[0]}" unit="{vals[1]}">{name}</ingredient>' for name, vals in recipe["ingredients"].items()])
    pre = "\n".join([f'\t\t<step>{step}</step>' for step in recipe["preparation"]])
    xml = f"""
<recipe>
    <title>{recipe["title"]}</title>
    <ingredients>
{ing}
    </ingredients>
    <preparation>
{pre}
    </preparation>
</recipe>
""".strip()
    return xml



if __name__ == "__main__":

    if len(sys.argv) != 2:
        print("error")
        exit(-1)
    #print(getHTMLFromWebsiteChromium(sys.argv[1]))
    
    parseLocalHTML(sys.argv[1])

    print(toXML(recipe), encoding="utf-8")
