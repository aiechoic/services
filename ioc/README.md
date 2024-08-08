# IOC

The IOC library is a dependency injection and management tool implemented in Golang.

## Features:
* Dependency Injection: Simplifies object creation and management using the ioc library.
* Named Providers: Uses named providers to create and manage the same type of object.
* Health Checks: Integrates health check functionality to ensure the health status of system components.

## Provider
A Provider is a component responsible for creating instances of a specific type.  

```go
package main

import (
	"github.com/aiechoic/services/ioc"
)

type Engine struct{ Name string }

var engineProvider = ioc.NewProvider(func(c *ioc.Container) (*Engine, error) {
	engine := &Engine{}
	return engine, nil
})

```

## Container
A Container is a central registry that stores singleton objects.

```go
func main() {
    c := ioc.NewContainer()
    var engine *Engine = engineProvider.Get(c)
}
```


more examples can be found in the [examples](./examples) directory.


## License
This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
```