package main

import (
	"errors"
	"fmt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"io/ioutil"
	"net/http"
	"strings"
)

type responseBody struct {
	To        string `json:"To"`
	From      string `json:"From"`
	Date      string `json:"Date"`
	Subject   string `json:"subject"`
	MessageID string `json:"Message-Id"`
}

const (
	mainBodyStr  = "Content-Type: multipart/mixed;"
	msgBorderStr = "msg_border"
	errBodyParse = "body doesn't contain required data"
)

func main() {
	// Echo instance
	e := echo.New()

	// Middleware
	//e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.POST("/mail", hello)

	// Start server
	e.Logger.Fatal(e.Start(":4343"))
}

// Handler
func hello(c echo.Context) error {
	x := c.Request()
	body := x.Body
	bodyBytes, err := ioutil.ReadAll(body)
	if err != nil {
		fmt.Printf("error %+v", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	response, err := parseBody(string(bodyBytes))

	if err != nil {
		fmt.Printf("error %+v", err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, response)
}

func parseBody(body string) (*responseBody, error) {
	response := &responseBody{}
	bodySplit := strings.Split(body, mainBodyStr)
	if len(bodySplit) < 2 {
		return nil, errors.New(errBodyParse)
	}
	mainData := strings.Split(bodySplit[1], msgBorderStr)
	if len(mainData) < 2 {
		return nil, errors.New(errBodyParse)
	}
	dataSplit := strings.Split(mainData[1], "\n")
	for _, k := range dataSplit {
		s := strings.Split(k, ":")
		if len(s) < 2 {
			continue
		}
		switch s[0] {
		case "To":
			response.To = strings.TrimSpace(s[1])
		case "From":
			response.From = strings.TrimSpace(s[1])
		case "Date":
			response.Date = strings.TrimSpace(s[1])
		case "Subject":
			response.Subject = strings.TrimSpace(s[1])
		case "Message-ID":
			response.MessageID = strings.TrimSpace(s[1])
		default:

		}
	}
	return response, nil
}
