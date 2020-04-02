package main

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type XMLLazy struct {
	path  string
	files []os.FileInfo
}

func (xmlLazy *XMLLazy) Init() {
	var err error
	xmlLazy.files, err = ioutil.ReadDir(xmlLazy.path)
	if err != nil {
		handleError(err)
	}
}

func (xmlLazy XMLLazy) Add(recipe Recipe) {
	content, err := xml.MarshalIndent(recipe, "", "	")
	if err != nil {
		handleError(err)
	}
	err = ioutil.WriteFile(xmlLazy.path+recipe.Title+".xml", content, 0644)
	if err != nil {
		handleError(err)
	}
}

func (xmlLazy XMLLazy) Iterator() <-chan Recipe {
	var filenames []string
	for _, file := range xmlLazy.files {
		filenames = append(filenames, file.Name())
	}
	return xmlLazy.Get(filenames)
}

func (xmlLazy XMLLazy) Get(filenames []string) <-chan Recipe {
	ch := make(chan Recipe)
	go func() {
		defer close(ch)
		for _, fname := range filenames {
			content, err := ioutil.ReadFile(xmlLazy.path + fname)

			if err != nil {
				handleError(err)
			}

			recipe := xmlLazy.ParseXMLContent(content)
			recipe.SetFilename(fname)
			ch <- recipe
		}
	}()
	return ch
}

func (xmlLazy XMLLazy) ParseXMLContent(content []byte) Recipe {
	recipe := Recipe{}
	if err := xml.Unmarshal(content, &recipe); err != nil {
		handleError(err)
	}
	return recipe
}
