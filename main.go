package main

import (
	"bytes"
	_ "embed"
	"fmt"
	"html/template"
	"log"
	"os"
	"strings"
	"time"

	wkhtml "github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type data struct {
	Name  string
	Month string
	Rows  []row
}

type row struct {
	Label string
	Index int
}

//go:embed tmpl/index.html
var tmpl string

func main() {
	t := template.Must(template.New("index.html").Parse(tmpl))

	htmlFiles := make([]string, 0)
	arr := buildCalendarData(2022)
	for _, data := range arr {
		buf := new(bytes.Buffer)
		err := t.Execute(buf, data)
		if err != nil {
			panic(err)
		}
		htmlFiles = append(htmlFiles, buf.String())
	}

	err := saveHtml("./out", htmlFiles)
	if err != nil {
		panic(err)
	}
}

func buildCalendarData(year int) []data {
	res := make([]data, 0)
	start := time.Date(year, 1, 1, 0, 0, 0, 0, time.UTC)
	for month := start; month.Year() == start.Year(); month = month.AddDate(0, 1, 0) {
		data := data{
			Name:  "Calendrier",
			Month: getMonthLabel(month.Month()),
			Rows:  make([]row, 0),
		}
		for d := month; d.Month() == month.Month(); d = d.AddDate(0, 0, 1) {
			data.Rows = append(data.Rows, row{
				Label: getDayLabel(d.Weekday()),
				Index: d.Day(),
			})
		}
		res = append(res, data)
	}
	return res
}

func saveHtml(folder string, htmlStrings []string) error {
	for i, str := range htmlStrings {
		if err := os.WriteFile(fmt.Sprintf("%s/%d.html", folder, i+1), []byte(str), 0644); err != nil {
			return err
		}
	}
	return nil
}

func toPDF(filepath string, htmlStrings []string) error {
	pdfg, err := wkhtml.NewPDFGenerator()
	pdfg.MarginTop.Set(0)
	pdfg.MarginBottom.Set(0)
	pdfg.MarginLeft.Set(0)
	pdfg.MarginRight.Set(0)
	if err != nil {
		return err
	}

	for _, htmlStr := range htmlStrings {
		pdfg.AddPage(wkhtml.NewPageReader(strings.NewReader(htmlStr)))
	}

	err = pdfg.Create()
	if err != nil {
		log.Fatal(err)
	}

	return pdfg.WriteFile(filepath)
}
