package dingtalk_test

import (
	"github.com/aiechoic/services/ioc"
	"github.com/aiechoic/services/message/notifier/dingtalk"
)

var markdown = `# Test
- test1
- test2

[Google](https://www.google.com/)
`

func ExampleClient_Notify() {
	c := ioc.NewContainer()
	err := c.LoadConfig("../../../configs", ioc.ConfigEnvTest)
	if err != nil {
		panic(err)
	}
	client := dingtalk.GetClient(c)

	err = client.Notify(markdown)
	if err != nil {
		panic(err)
	}

	// Output:

}
