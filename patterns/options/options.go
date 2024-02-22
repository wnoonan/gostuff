package options

import "fmt"

type Crust int
type Topping int

const (
	Thin Crust = iota
	Thick
	Stuffed
)

const (
	Pepperoni Topping = iota
	Sausage
	Mushroom
	GreenOnion
	Onion
	Cheese
	Sauce
)

type Pie interface {
	Cook()
	Deliver() string
}

type Pizza struct {
	Size     int
	Crust    Crust
	Toppings []Topping
	Count    int
	Oven     Oven
}

type PizzaOption func(*Pizza)

type Oven interface {
	Done() bool
	Heat(amount int)
}

type FireBrickOven struct {
	Temperature int
}

type ConventionalOven struct {
	Temperature int
}

func NewPizza(options ...PizzaOption) *Pizza {
	// We set some defaults, you'll get a 12 inch thing crust empty pie by default
	// and it will be cooked in a conventional oven at 400 degrees.
	p := &Pizza{
		Size:     12,
		Crust:    Thin,
		Toppings: []Topping{},
		Count:    1,
		Oven:     &ConventionalOven{Temperature: 400},
	}

	// We then loop over the options and apply them to the pizza, this is where end user customizations come in.
	for _, option := range options {
		option(p)
	}

	return p
}

func WithSize(s int) PizzaOption {
	return func(p *Pizza) {
		p.Size = s
	}
}

func WithThinCrust() PizzaOption {
	return func(p *Pizza) {
		p.Crust = Thin
	}
}

func WithThickCrust() PizzaOption {
	return func(p *Pizza) {
		p.Crust = Thick
	}
}

func WithStuffedCrust() PizzaOption {
	return func(p *Pizza) {
		p.Crust = Stuffed
	}
}

func WithPepperoni() PizzaOption {
	return func(p *Pizza) {
		p.Toppings = append(p.Toppings, Pepperoni)
	}
}

func WithSausage() PizzaOption {
	return func(p *Pizza) {
		p.Toppings = append(p.Toppings, Sausage)
	}
}

func WithMushroom() PizzaOption {
	return func(p *Pizza) {
		p.Toppings = append(p.Toppings, Mushroom)
	}
}

func WithGreenOnion() PizzaOption {
	return func(p *Pizza) {
		p.Toppings = append(p.Toppings, GreenOnion)
	}
}

func WithOnion() PizzaOption {
	return func(p *Pizza) {
		p.Toppings = append(p.Toppings, Onion)
	}
}

func WithCheese() PizzaOption {
	return func(p *Pizza) {
		p.Toppings = append(p.Toppings, Cheese)
	}
}

func WithSauce() PizzaOption {
	return func(p *Pizza) {
		p.Toppings = append(p.Toppings, Sauce)
	}
}

func WithOven(o Oven) PizzaOption {
	return func(p *Pizza) {
		p.Oven = o
	}
}

func (p *Pizza) Cook() {
	for !p.Oven.Done() {
		p.Oven.Heat(10)
	}

}

func (p *Pizza) Deliver() string {
	return fmt.Sprintf("Your %d inch pizza with a %s crust and %v toppings is ready!", p.Size, p.Crust.String(), p.Toppings)
}

func (fbo *FireBrickOven) Done() bool {
	return fbo.Temperature > 500
}

func (co *ConventionalOven) Done() bool {
	return co.Temperature > 500
}

func (fbo *FireBrickOven) Heat(amount int) {
	fbo.Temperature += amount
}

func (co *ConventionalOven) Heat(amount int) {
	co.Temperature += amount
}

func (t Topping) String() string {
	switch t {
	case Pepperoni:
		return "Pepperoni"
	case Sausage:
		return "Sausage"
	case Mushroom:
		return "Mushroom"
	case GreenOnion:
		return "Green Onion"
	case Onion:
		return "Onion"
	case Cheese:
		return "Cheese"
	case Sauce:
		return "Sauce"
	default:
		return "Unknown"
	}
}

func (c Crust) String() string {
	switch c {
	case Thin:
		return "Thin"
	case Thick:
		return "Thick"
	case Stuffed:
		return "Stuffed"
	default:
		return "Unknown"
	}
}
