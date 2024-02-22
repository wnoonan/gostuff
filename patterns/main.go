package main

import (
	"fmt"

	"github.com/wnoonan/gostuff/options/fluent"
	"github.com/wnoonan/gostuff/options/options"
)

func main() {
	functionalOptions()
	fluentInterface()
}

func functionalOptions() {
	// conventional oven
	pepperoniPizza := options.NewPizza(options.WithCheese(), options.WithPepperoni(), options.WithThinCrust())

	pepperoniPizza.Cook()
	fmt.Println(pepperoniPizza.Deliver())

	// substitute another oven
	cheesePizza := options.NewPizza(options.WithCheese(), options.WithStuffedCrust(), options.WithOven(&options.FireBrickOven{Temperature: 900}))

	cheesePizza.Cook()
	fmt.Println(cheesePizza.Deliver())
}

func fluentInterface() {
	dog := &fluent.Dog{}

	dog.Name("Fido").Breed("Golden Retriever").Age(3).Color("Golden")

	fmt.Println(dog.Summary())
	fmt.Println(dog.Speak())
	fmt.Println(dog.Fetch())

	dog.Name("Rex").Breed("German Shepherd").Age(5).Color("Brown and Black")

	fmt.Println(dog.Summary())
}
