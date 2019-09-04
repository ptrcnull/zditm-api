package main

import (
	"bytes"
	"errors"
	"golang.org/x/net/html"
	"io"
	"strings"
)

type Departure struct {
	Line string `json:"line"`
	Direction string `json:"direction"`
	Time string `json:"time"`
}

func iterate(el *html.Node, handler func(*html.Node)) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		handler(n)
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}
	f(el)
}

func getBody(doc *html.Node) (*html.Node, error) {
	var tbody *html.Node
	iterate(doc, func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "tbody" {
			tbody = node
		}
	})
	if tbody != nil {
		return tbody, nil
	}
	return nil, errors.New("missing <tbody> in the node tree")
}

func getRows(doc *html.Node) []*html.Node {
	var arr []*html.Node
	iterate(doc, func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "tr" {
			arr = append(arr, node)
		}
	})
	return arr
}

func getDeparture(row *html.Node) Departure {
	departure := Departure{}
	iterate(row, func(node *html.Node) {
		if node.Type == html.ElementNode && node.Data == "td" {
			for _, attrib := range node.Attr {
				if attrib.Key != "class" {
					continue
				}
				if attrib.Val == "gmvlinia" {
					departure.Line = getText(node)
				}
				if attrib.Val == "gmvkierunek" {
					departure.Direction = getText(node)
				}
				if attrib.Val == "gmvgodzina" {
					departure.Time = getText(node)
				}
			}
		}
	})
	return departure
}

func getText(node *html.Node) string {
	var text string
	iterate(node, func(node *html.Node) {
		if node.Type == html.TextNode {
			text = node.Data
		}
	})
	return text
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	_ = html.Render(w, n)
	return buf.String()
}


func ParseHTML(htm string) ([]Departure, error) {
	doc, err := html.Parse(strings.NewReader(htm))
	if err != nil {
		return nil, err
	}

	tbody, err := getBody(doc)
	if err != nil {
		return nil, err
	}

	rows := getRows(tbody)
	var departures []Departure
	for _, row := range rows {
		departures = append(departures, getDeparture(row))
	}
	//departures = append(departures, Departure{ Line: "", Direction: "", Time: renderNode(tbody) })
	return departures, nil
}
