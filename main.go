package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/gocolly/colly"
)

type EtfInfo struct {
	Title              string
	Replication        string
	Earnings           string
	TotalExpenceRatio  string
	TrackingDifference string
	FundSize           string
}

func main() {

	etfInfo := EtfInfo{}

	scrapeUrl := "https://www.trackingdifferences.com/ETF/ISIN/IE00B1XNHC34"

	c := colly.NewCollector(colly.AllowedDomains("www.trackingdifferences.com"))

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("Accept-Language", "en-US;q=0.9")
		fmt.Printf("Visiting %s\n", r.URL)
	})

	c.OnError(func(r *colly.Response, e error) {
		fmt.Printf("Error while scrapping : %s\n", e.Error())
	})

	c.OnHTML("h1.page-title", func(h *colly.HTMLElement) {
		etfInfo.Title = h.Text
	})

	c.OnHTML("div[descfloat] p[desc]", func(h *colly.HTMLElement) {
		selection := h.DOM

		childNodes := selection.Children().Nodes
		if len(childNodes) == 3 {
			description := cleanDesc(selection.Find("span.desctitle").Text())

			value := selection.FindNodes(childNodes[2]).Text()

			switch description {

			case "Replication":
				etfInfo.Replication = value
				break
			case "TER":
				etfInfo.TotalExpenceRatio = value
				break
			case "TD":
				etfInfo.TrackingDifference = value
				break
			case "Earnings":
				etfInfo.FundSize = value
				break
			}
		}
	})

	c.OnScraped(func(r *colly.Response) {
		enc := json.NewEncoder(os.Stdout)

		enc.SetIndent("", " ")

		enc.Encode(etfInfo)
	})

	c.Visit(scrapeUrl)
}

func cleanDesc(s string) string {
	return strings.TrimSpace(s)
}
