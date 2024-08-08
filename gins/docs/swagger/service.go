package swagger

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
    <meta charset="utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
    <meta name="description" content="SwaggerUI" />
    <title>SwaggerUI</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@latest/swagger-ui.css" />
  </head>
  <body>
  <div id="swagger-ui"></div>
  <script src="https://unpkg.com/swagger-ui-dist@latest/swagger-ui-bundle.js" crossorigin></script>
  <script src="https://unpkg.com/swagger-ui-dist@latest/swagger-ui-standalone-preset.js" crossorigin></script>
  <script>
    window.onload = () => {
	  window.ui = SwaggerUIBundle({
		url: "/openapi.json",
		dom_id: '#swagger-ui',
		deepLinking: true,
		presets: [
		  SwaggerUIBundle.presets.apis,
		  SwaggerUIStandalonePreset
		],
		plugins: [
		  SwaggerUIBundle.plugins.DownloadUrl
		],
		layout: "StandaloneLayout",
		persistAuthorization: true,
	  });
    };
  </script>
  </body>
</html>
`

func ServeAPI(s *gins.Server) {
	data, err := json.Marshal(s.API)
	if err != nil {
		panic(err)
	}
	s.Engine.GET("/docs", func(c *gin.Context) {
		_, _ = c.Writer.WriteString(html)
	})
	s.Engine.GET("/openapi.json", func(c *gin.Context) {
		c.JSON(200, json.RawMessage(data))
	})

	fmt.Printf("serve swagger at http://localhost:%d/docs\n", s.Port)
}
