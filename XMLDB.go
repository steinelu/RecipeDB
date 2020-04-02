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

func (self *XMLLazy) Init() {
	var err error
	self.files, err = ioutil.ReadDir(self.path)
	if err != nil {
		handleError(err)
	}
}

func (self XMLLazy) Add(recipe Recipe) {
	content, err := xml.MarshalIndent(recipe, "", "	")
	if err != nil {
		handleError(err)
	}
	err = ioutil.WriteFile(self.path+recipe.Title+".xml", content, 0644)
	if err != nil {
		handleError(err)
	}
}

func (self XMLLazy) Iterator() <-chan Recipe {
	filenames := []string{}
	for _, file := range self.files {
		filenames = append(filenames, file.Name())
	}
	return self.Get(filenames)
}

func (self XMLLazy) Get(filenames []string) <-chan Recipe {
	ch := make(chan Recipe)
	go func() {
		defer close(ch)
		for _, fname := range filenames {
			recipe := Recipe{}
			content, err := ioutil.ReadFile(self.path + fname)

			if err != nil {
				handleError(err)
			}

			if err := xml.Unmarshal(content, &recipe); err != nil {
				handleError(err)
			}
			recipe.SetFilename(fname)
			ch <- recipe
		}
	}()
	return ch
}
