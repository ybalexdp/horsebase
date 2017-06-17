package main

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func Test_GetRaceHTML(t *testing.T) {
	hb := &Horsebase{}
	hb = hb.New()

	err := hb.GetRaceHTML()
	if err != nil {
		t.Fatalf("GetRaceHTML error:%s", err)
	}
}

func Test_MakeRaceURLList(t *testing.T) {
	hb := &Horsebase{}
	hb = hb.New()

	err := hb.MakeRaceURLList()
	if err != nil {
		t.Fatalf("MakeRaceURLList error:%s", err)
	}
}

func Test_RegistHorseData(t *testing.T) {
	hb := &Horsebase{}
	hb = hb.New()

	err := hb.RegistHorseData()
	if err != nil {
		t.Fatalf("RegistHorseData error:%s", err)
	}
}

/*
func Test_RegistRaceData(t *testing.T) {
	hb := &Horsebase{}
	hb = hb.New()

	err := hb.RegistRaceData()
	if err != nil {
		t.Fatalf("RegistRaceData error:%s", err)
	}

}
*/

func Test_getRaceIDfromHTML(t *testing.T) {

	hb := &Horsebase{}
	hb = hb.New()

	url := getRaceIDfromHTML("http://db.netkeiba.com/race/201703010103")

	expected := "201703010103"

	if url != expected {
		t.Fatalf("RaceID is not expected:%s", url)
	}
}

func Test_calcDifTime(t *testing.T) {
	var result RaceResultData
	var ftime time.Time
	var text string
	text = "0:20.8"
	var racedata RaceData

	fp, err := os.Open("./test/201709030109.html")
	if err != nil {
		t.Fatalf("%s", err)
	}
	defer fp.Close()

	doc, err := goquery.NewDocumentFromReader(fp)
	if err != nil {
		t.Fatalf("%s", err)
	}

	s := doc.Find("title").First()
	title := strings.Split(s.Text(), "｜")

	racedata.Date = getRaceDate(title[1])

	s = doc.Find("td").First().Next().Next().Next().Next().Next().Next().Next()

	ftime = racedata.getRaceTime(text)

	result.Time = racedata.getRaceTime(s.Text())

	result.DifTime = calcDifTime(result, ftime)

	if 0 > result.DifTime {
		t.Fatalf("calcDifTime error:%s", err)
		t.Fatalf("Diftime:%f", result.DifTime)
	}
}
