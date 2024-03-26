package main

import (
	"fmt"

	"github.com/wnoonan/gostuff/options/fluent"
)

func main() {
	fluentInterface()
}

func fluentInterface() {
	// dog := &fluent.Dog{}

	// dog.Name("Fido").Breed("Golden Retriever").Age(3).Color("Golden")

	// fmt.Println(dog.Summary())
	// fmt.Println(dog.Speak())
	// fmt.Println(dog.Fetch())

	// dog.Name("Rex").Breed("German Shepherd").Age(5).Color("Brown and Black")

	// fmt.Println(dog.Summary())

	pug := &fluent.Dog{}

	pug.Name("Pugsley").Breed("Pug").Age(2).Color("Grey")

	fmt.Println(pug.Summary())
}
