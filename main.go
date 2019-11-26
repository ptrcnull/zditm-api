package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
)

type ZditmResponse struct {
	Tresc string
	Komunikat string
}

func main() {
	r := gin.Default()

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

