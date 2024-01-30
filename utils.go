package main

import (
	"fmt"
	"io"
	"math"
	"net/http"

	"github.com/Jeffail/gabs/v2"
)

var svgStringFunc = func(text string, textColor string, backgroundColor string, strokeOffset string) string {
	return fmt.Sprintf(`<svg width="128" height="128" viewBox="-16 -16 160 160" version="1.1" xmlns="http://www.w3.org/2000/svg" style="transform:rotate(-90deg)">
    <circle r="54" cx="64" cy="64" fill="transparent" stroke="#e0e0e0" stroke-width="12" stroke-dasharray="339px" stroke-dashoffset="0"></circle>
    <circle r="54" cx="64" cy="64" stroke="%s" stroke-width="12" stroke-linecap="round" stroke-dashoffset="%s" fill="transparent" stroke-dasharray="339px"></circle>
    <text text-anchor="middle" transform="translate(52,64) rotate(90)" fill="%s" font-size="34px" font-weight="bold" style="font-family: Segoe UI;">%s</text>
</svg>
`, backgroundColor, strokeOffset, textColor, text)
}

const (
	circumfence float64 = 339
)

func GetSVG(url string) string {
	percentage, err := GetPerformancePercentage(url)
	if err != nil {
		return svgStringFunc("ERR", "#c00", "#ff4e43", "0px")
	}
	circlePercentage := circumfence
	if percentage > 0 {
		circlePercentage = math.Round(circumfence - (circumfence * percentage))
	}
	intPercentage := int(math.Round(percentage * 100))
	textColor := "#080" // Good
	percentageColor := "#0cce6a"
	if intPercentage <= 30 {
		textColor = "#c00" // Bad
		percentageColor = "#ff4e43"
	} else if intPercentage <= 60 {
		textColor = "#C33300" // Medium
		percentageColor = "#ffa400"
	}
	return svgStringFunc(fmt.Sprintf("%d", intPercentage), textColor, percentageColor, fmt.Sprintf("%fpx", circlePercentage))
}

func GetPerformancePercentage(url string) (float64, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("https://www.googleapis.com/pagespeedonline/v5/runPagespeed?url=%s", url), nil)
	if err != nil {
		fmt.Print("Req errored")
		fmt.Printf("Req errored: %s", err)
		return -1, err
	}

	req.Close = true
	req.Header.Add("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if resp != nil {
		defer resp.Body.Close()
	}
	if err != nil {
		fmt.Print("Do errored")
		fmt.Printf("Do errored: %s", err)
		return -1, err
	}
	byteValue, _ := io.ReadAll(resp.Body)
	jsonData, err := gabs.ParseJSON(byteValue)
	if err != nil {
		fmt.Print("Gabs errored")
		fmt.Printf("Gabs errored: %s", err)
		return -1, err
	}
	value, ok := jsonData.Search("lighthouseResult", "categories", "performance", "score").Data().(float64)
	if !ok {
		fmt.Print("Json search failed")
		return -1, fmt.Errorf("json search failed")
	}
	fmt.Printf("%s scored %v\n", url, value)
	return value, nil
}
