package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"strings"
)

type ZditmResponse struct {
	Tresc     string
	Komunikat string
}

const homepage = `<!DOCTYPE html>
<head>
	<title>zditm</title>
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
	<style>
		body {
			font-family: monospace;
		}
	</style>
</head>
<body>
	<p>
		GET /json/:id
		<br />
		Example: <a href="http://{host}/json/30311">http://{host}/json/30311</a>
	</p>
	<p>
		GET /text/:id
		<br />
		Example: <a href="http://{host}/text/30311">http://{host}/text/30311</a>
		<br />
		Best used with a HTTP widget like <a href="https://apkpure.com/http-widget/com.axgs.httpwidget">this one</a>.
	</p>
	<p>
		ID can be found at every stop next to the name
		or online at <a href="https://www.zditm.szczecin.pl/pl/mapy/przystanki-i-pojazdy">https://www.zditm.szczecin.pl/pl/mapy/przystanki-i-pojazdy</a>.
	</p>
</body>
`

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		host := c.Request.Host
		c.Data(http.StatusOK, "text/html; charset=utf-8", []byte(
			strings.ReplaceAll(homepage, "{host}", host),
		))
	})

	r.GET("/json/:id", func(c *gin.Context) {
		id, _ := c.Params.Get("id")
		html, err := MakeRequest(id)
		if err != nil {
			c.AbortWithStatusJSON(500, err)
		}

		body, err := ParseHTML(html)
		if err != nil {
			c.AbortWithStatusJSON(500, err)
		}
		c.JSON(200, body)
		return
	})

	r.GET("/text/:id", func(c *gin.Context) {
		id, _ := c.Params.Get("id")
		html, err := MakeRequest(id)
		if err != nil {
			c.AbortWithStatusJSON(500, err)
		}

		body, err := ParseHTML(html)
		if err != nil {
			c.AbortWithStatusJSON(500, err)
		}
		var res string
		for _, row := range body {
			res += row.Line + ": " + row.Direction + " " + row.Time + "\n"
		}

		c.String(200, res)
		return
	})

	if err := r.Run(":38126"); err != nil {
		panic(err)
	}
}

func MakeRequest(id string) (string, error) {
	resp, err := http.Get(fmt.Sprintf("https://www.zditm.szczecin.pl/json/tablica.inc.php?lng=en&slupek=%s", id))
	if err != nil {
		return "", err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	res := ZditmResponse{}
	if err := json.Unmarshal(body, &res); err != nil {
		return "", err
	}

	return res.Tresc, nil
}
