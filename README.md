nbclassifier-go
===============

A Naive Bayes classifier in Golang

Usage
-----

```go
package main

import (
	"fmt"

	"github.com/ajph/nbclassifier-go"
)

func main() {
	m := nbclassifier.New()

	m.NewClass("banana")
	m.Learn("banana", "sweet", "yellow", "long")

	m.NewClass("orange")
	m.Learn("orange", "sweet", "orange", "round")

	m.NewClass("apple")
	m.Learn("apple", "sweet", "green", "red", "round")

	//m.SaveToFile("./test.json") 

	/*
		// or load from existing file
		m, err := nbclassifier.LoadFromFile("./test.json")
		if err != nil {
			panic(err)
		} 
	*/

	w, unsure, _ := m.Classify("round", "sweet")

	fmt.Println(w.Id, unsure)
}
```

Todo
----
- Documentation
- Tests
