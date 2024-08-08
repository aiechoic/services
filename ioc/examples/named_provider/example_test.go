package main

import (
	"fmt"
	"github.com/aiechoic/services/ioc"
)

type Wheel struct {
	Name string
}

var wheelProviders = ioc.NewProviders[*Wheel]()

func GetWheelProvider(provider string) *ioc.Provider[*Wheel] {
	return wheelProviders.GetProvider(provider, func(c *ioc.Container) (*Wheel, error) {
		return &Wheel{
			Name: "Wheel " + provider,
		}, nil
	})
}

func GetWheel(c *ioc.Container, name string) *Wheel {
	return GetWheelProvider(name).MustGet(c)
}

func Example() {

	c := ioc.NewContainer()
	defer c.Close()
	providerNames := []string{
		"front-left",
		"front-right",
		"rear-left",
		"rear-right",
	}
	for _, providerName := range providerNames {
		wheel := GetWheel(c, providerName)
		fmt.Println(wheel.Name)
	}

	// Output:
	// Wheel front-left
	// Wheel front-right
	// Wheel rear-left
	// Wheel rear-right
}
