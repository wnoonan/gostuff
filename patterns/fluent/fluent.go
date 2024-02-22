package fluent

import "fmt"

type Doggo interface {
	Name(string) Breeder
	Speak() string
	Fetch() string
	Summary() string
}

type Colorer interface {
	Color(string) Doggo
}

type Breeder interface {
	Breed(string) Ager
}

type Ager interface {
	Age(int) Colorer
}

type Dog struct {
	name  string
	age   int
	breed string
	color string
}

func (d *Dog) Name(name string) Breeder {
	d.name = name

	return d
}

func (d *Dog) Breed(breed string) Ager {
	d.breed = breed

	return d
}

func (d *Dog) Age(age int) Colorer {
	d.age = age

	return d
}

func (d *Dog) Color(color string) Doggo {
	d.color = color

	return d
}

func (d *Dog) Speak() string {
	return "Woof!"
}

func (d *Dog) Fetch() string {
	return "I fetched it!"
}

func (d *Dog) Summary() string {
	return fmt.Sprintf("My name is %s, I am a %s, I am %d years old, and I am %s.", d.name, d.breed, d.age, d.color)
}
