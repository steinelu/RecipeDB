package main

import (
	"container/heap"
	"fmt"
	bolt "go.etcd.io/bbolt"
	"log"
	"strings"
)

type Index struct {
	path string
	db *bolt.DB
	recipes DataBase
	bucketName string
}

func (index *Index) Init(recipeDB DataBase){
	index.recipes = recipeDB
	index.bucketName = "Index"

	db, err := bolt.Open("./tmp/test_index.db", 0600, nil)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Update(func(tx *bolt.Tx) error{
		_, err := tx.CreateBucketIfNotExists([]byte(index.bucketName))
		return err
	})
	index.db = db

}

func (index Index) Seperator() string{
	return ";"
}

func (index *Index) Index(recipe Recipe){
	terms := parseRecipe(recipe)
	for _, term := range terms {
		index.add(term, recipe.GetId())
	}
}

func (index *Index) add(term string, id string){
	key := []byte(term)

	err := index.db.Update(func(tx *bolt.Tx) error{
		indexBucket, err := tx.CreateBucketIfNotExists([]byte(index.bucketName))
		if err != nil{
			return err
		}

		val := string(indexBucket.Get(key))
		list := StringHeap(strings.Split(val, index.Seperator()))
		heap.Init(&list)
		heap.Push(&list, id)
		err = indexBucket.Put(key, []byte(strings.Join(list, index.Seperator())))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
}

func (index *Index)get(term string) []string {
	var val []byte
	err := index.db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte(index.bucketName))
		if bucket == nil {
			return fmt.Errorf("Bucket not found")
		}

		val = bucket.Get([]byte(term))
		return nil
	})

	if err != nil {
		log.Fatal(err)
	}
	// TODO fast enought?
	s := strings.Split(string(val), index.Seperator())
	// first item is always empty, because of empty get in get
	return s[1:]
}

func (index *Index) Search(terms []string) <- chan Recipe{
	// boolean retrival
	res := index.get(terms[0])
	for _, term := range terms[1:] {
		res = intersect(res, index.get(term))
	}


	//return (*index.recipes).Get(res)
	return index.recipes.Get(res)
}



func (index *Index) Close(){
	if index.db != nil{
		index.db.Close()
		index.db = nil
	}
}

type StringHeap []string

func (h StringHeap) Len() int           { return len(h) }
func (h StringHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h StringHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *StringHeap) Push(x interface{}) {
	*h = append(*h, x.(string))
}

func (h *StringHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func intersect(one []string, two []string) []string {
	if one == nil || two == nil || len(one) <= 0 || len(two) <= 0{
		return []string{}
	}
	intersection := StringHeap{}
	for i, j := 0, 0; i < len(one) && j < len(two); {
		if one[i] < two[j] {
			i++
		} else if one[i] > two[j] {
			j++
		} else {
			intersection = append(intersection, one[i])
			i++
			j++
		}
	}
	return intersection
}

