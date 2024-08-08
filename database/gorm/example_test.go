package gorm_test

import (
	"fmt"
	"github.com/aiechoic/services/database/gorm"
	"github.com/aiechoic/services/ioc"
)

func ExampleGetGormDB() {
	c := ioc.NewContainer()
	defer c.Close()

	err := c.LoadConfig("../../configs", ioc.ConfigEnvTest)
	if err != nil {
		panic(err)
	}

	db := gorm.GetGormDB(c)

	fmt.Printf("%v\n", db == nil)

	// Output: false
}
