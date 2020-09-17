package main

import (
	"encoding/xml"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"log"
)


func UnmarschalXMLRecipe(content []byte) Recipe {
	recipe := Recipe{}
	if err := xml.Unmarshal(content, &recipe); err != nil {
		handleError(err)
	}
	return recipe
}

func MarschalXMLRecipe(recipe Recipe) []byte{
	content, err := xml.MarshalIndent(recipe, "", "	")
	if err != nil {
		handleError(err)
	}
	return content
}

type boltdb struct{
	path string
	db_ptr *bolt.DB
}

func (db *boltdb)Init(){
	db_, err := bolt.Open("./tmp/test.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db_.Update(func(tx *bolt.Tx) error{
		_, err := tx.CreateBucketIfNotExists([]byte("Recipes"))
		return err
	})

	db.db_ptr = db_
}

func (db *boltdb) Close(){
	if db.db_ptr != nil{
		db.db_ptr.Close()
		db.db_ptr = nil
	}
}

func (db *boltdb) Add(recipe Recipe){
	key := []byte(recipe.GetId())
	val := MarschalXMLRecipe(recipe)

	err := db.db_ptr.Update(func(tx *bolt.Tx) error{
		bucket, err := tx.CreateBucketIfNotExists([]byte("Recipes"))
		if err != nil{
			return err
		}

		err = bucket.Put(key, val)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

}

func (db *boltdb) Get(keys []string) <- chan Recipe{
	recipes := make(chan Recipe)
	go func(){
		for _, key := range keys{
			fmt.Println(key)
			recipes <- db.get([]byte(key))
		}
		defer close(recipes)
	}()
	return recipes
}

func (db *boltdb) get(key []byte) Recipe{
	var val []byte
	err := db.db_ptr.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Recipes"))
		if bucket == nil {
			return fmt.Errorf("Bucket B1 not found")
		}

		val = bucket.Get(key)
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}

	return UnmarschalXMLRecipe(val)
}

func (db *boltdb) Iterator(apply func(Recipe, []byte)){
	err := db.db_ptr.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("Recipes"))
		if bucket == nil {
			return fmt.Errorf("Bucket 'Recipes' not found")
		}

		err := bucket.ForEach(func(k, v []byte) error {
			apply(UnmarschalXMLRecipe(v), k)
			return nil
		})

		if err != nil {
			return fmt.Errorf("Error while reading all recipes!: ")
		}
		return nil
	})

	if err != nil{
		log.Fatal(err)
	}
}
