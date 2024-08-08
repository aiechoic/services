package redoc

import (
	"encoding/json"
	"fmt"
	"github.com/aiechoic/services/gins"
	"github.com/gin-gonic/gin"
)

var html = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>API Documentation</title>
	<redoc spec-url="/redoc-api.json"></redoc>
    <script src="https://cdn.redoc.ly/redoc/latest/bundles/redoc.standalone.js"> </script>
</head>
<body>
</body>
</html>
`

func ServeAPI(s *gins.Server) {
	data, err := json.Marshal(s.API)
	if err != nil {
		panic(err)
	}
	s.Engine.GET("/redoc", func(c *gin.Context) {
		_, _ = c.Writer.WriteString(html)
	})

	s.Engine.GET("/redoc-api.json", func(c *gin.Context) {
		c.JSON(200, json.RawMessage(data))
	})

	fmt.Printf("serve redoc at http://localhost:%d/redoc\n", s.Port)
}
