package simple

import (
	"fmt"
	"github.com/aiechoic/services/ioc"
)

type Map struct{ Name string }

var mapProvider = ioc.NewProvider(func(c *ioc.Container) (*Map, error) {
	return &Map{}, nil
})

type Engine struct{ Name string }

var engineProvider = ioc.NewProvider(func(c *ioc.Container) (*Engine, error) {
	engine := &Engine{}
	return engine, nil
})

type Car struct {
	Engine *Engine
	Map    *Map
}

var carProvider = ioc.NewProvider(func(c *ioc.Container) (*Car, error) {
	singleMap := mapProvider.MustGet(c)       // get singleton Map
	newEngine := engineProvider.MustGetNew(c) // always create new Engine when Car is created
	return &Car{
		Engine: newEngine,
		Map:    singleMap,
	}, nil
})

func Example() {

	c := ioc.NewContainer()

	// create new Car, new Engine, new Map
	car := carProvider.MustGet(c)

	// create new Car, new Engine, use the same Map
	newCar := carProvider.MustGetNew(c)

	fmt.Printf("car == newCar: %v\n", car == newCar)
	fmt.Printf("car.Engine == newCar.Engine: %v\n", car.Engine == newCar.Engine)
	fmt.Printf("car.Map == newCar.Map: %v\n", car.Map == newCar.Map)

	// Output:
	// car == newCar: false
	// car.Engine == newCar.Engine: false
	// car.Map == newCar.Map: true
}
