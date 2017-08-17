package main

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
)

/*
func Test_GetRaceHTML(t *testing.T) {
	hb := &Horsebase{}
	hb = hb.New()

	err := hb.GetRaceHTML()
	if err != nil {
		t.Fatalf("GetRaceHTML error:%s", err)
	}
}
*/

/*
func Test_MakeRaceURLList(t *testing.T) {
	hb := &Horsebase{}
	hb = hb.New()

	err := hb.MakeRaceURLList()
	if err != nil {
		t.Fatalf("MakeRaceURLList error:%s", err)
	}
}
*/

/*
func Test_RegistHorseData(t *testing.T) {
	hb := &Horsebase{}
	hb = hb.New()

	err := hb.RegistHorseData()
	if err != nil {
		t.Fatalf("RegistHorseData error:%s", err)
	}
}
*/

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

func Test_checkOldestData(t *testing.T) {

	hb := &Horsebase{}
	hb = hb.New()

	if !hb.checkOldestData("/?pid=race_top&date=20170506") {
		t.Fatalf("checkOldestData error")
	}
}

func Test_convAgeGr(t *testing.T) {

	agegr := convAgeGr("2歳")
	if agegr != AgeGrTwo {
		t.Fatalf("convAgeGr error:%d", agegr)
	}

	agegr = convAgeGr("3歳")
	if agegr != AgeGrThree {
		t.Fatalf("convAgeGr error:%d", agegr)
	}

	agegr = convAgeGr("3歳以上")
	if agegr != AgeGrThreeOver {
		t.Fatalf("convAgeGr error:%d", agegr)
	}

	agegr = convAgeGr("4歳以上")
	if agegr != AgeGrFourOver {
		t.Fatalf("convAgeGr error:%d", agegr)
	}

	agegr = convAgeGr("5歳以上")
	if agegr != -1 {
		t.Fatalf("convAgeGr error:%d", agegr)
	}

}

func Test_convBelonging(t *testing.T) {

	belonging := convBelonging("東")
	if belonging != BelongingEast {
		t.Fatalf("convBelonging error:%d", belonging)
	}

	belonging = convBelonging("西")
	if belonging != BelongingWest {
		t.Fatalf("convBelonging error:%d", belonging)
	}

	belonging = convBelonging("地")
	if belonging != BelongingLocal {
		t.Fatalf("convBelonging error:%d", belonging)
	}

	belonging = convBelonging("外")
	if belonging != BelongingForeign {
		t.Fatalf("convBelonging error:%d", belonging)
	}

	belonging = convBelonging("他")
	if belonging != -1 {
		t.Fatalf("convBelonging error:%d", belonging)
	}

}

func Test_convCond(t *testing.T) {

	cond := convCond("良")
	if cond != ConditionGood {
		t.Fatalf("convCond error:%d", cond)
	}

	cond = convCond("稍重")
	if cond != ConditionYielding {
		t.Fatalf("convCond error:%d", cond)
	}

	cond = convCond("重")
	if cond != ConditionSoft {
		t.Fatalf("convCond error:%d", cond)
	}

	cond = convCond("不良")
	if cond != ConditionHeavy {
		t.Fatalf("convCond error:%d", cond)
	}

	cond = convCond("他")
	if cond != -1 {
		t.Fatalf("convCond error:%d", cond)
	}

}

func Test_convCorner(t *testing.T) {

	corner := convCorner("内")
	if corner != InnerCourse {
		t.Fatalf("convCorner error:%d", corner)
	}

	corner = convCorner("外")
	if corner != OuterCourse {
		t.Fatalf("convCorner error:%d", corner)
	}

	corner = convCorner("他")
	if corner != -1 {
		t.Fatalf("convCorner error:%d", corner)
	}

}

func Test_convCourse(t *testing.T) {

	course := convCourse("右")
	if course != RightTurn {
		t.Fatalf("convCourse error:%d", course)
	}

	course = convCourse("左")
	if course != LeftTurn {
		t.Fatalf("convCourse error:%d", course)
	}

	course = convCourse("直線")
	if course != Straight {
		t.Fatalf("convCourse error:%d", course)
	}

	course = convCourse("他")
	if course != -1 {
		t.Fatalf("convCourse error:%d", course)
	}

}

func Test_convGrade(t *testing.T) {

	grade := convGrade("新馬")
	if grade != GradeDebut {
		t.Fatalf("convGrade error:%d", grade)
	}

	grade = convGrade("未勝利")
	if grade != GradeNoWin {
		t.Fatalf("convGrade error:%d", grade)
	}

	grade = convGrade("500万下")
	if grade != Grade500 {
		t.Fatalf("convGrade error:%d", grade)
	}

	grade = convGrade("1000万下")
	if grade != Grade1000 {
		t.Fatalf("convGrade error:%d", grade)
	}

	grade = convGrade("1600万下")
	if grade != Grade1600 {
		t.Fatalf("convGrade error:%d", grade)
	}

	grade = convGrade("オープン")
	if grade != GradeOP {
		t.Fatalf("convGrade error:%d", grade)
	}

	grade = convGrade("G")
	if grade != GradeG {
		t.Fatalf("convGrade error:%d", grade)
	}

	grade = convGrade("G3")
	if grade != GradeG3 {
		t.Fatalf("convGrade error:%d", grade)
	}

	grade = convGrade("G2")
	if grade != GradeG2 {
		t.Fatalf("convGrade error:%d", grade)
	}

	grade = convGrade("G1")
	if grade != GradeG1 {
		t.Fatalf("convGrade error:%d", grade)
	}

	grade = convGrade("他")
	if grade != -1 {
		t.Fatalf("convGrade error:%d", grade)
	}

}

func Test_convSex(t *testing.T) {

	sex := convSex("牡")
	if sex != Male {
		t.Fatalf("convSex error:%d", sex)
	}

	sex = convSex("牝")
	if sex != Female {
		t.Fatalf("convSex error:%d", sex)
	}

	sex = convSex("セ")
	if sex != Gelding {
		t.Fatalf("convSex error:%d", sex)
	}

	sex = convSex("他")
	if sex != -1 {
		t.Fatalf("convSex error:%d", sex)
	}

}

func Test_convSurface(t *testing.T) {

	surface := convSurface("芝")
	if surface != Turf {
		t.Fatalf("convSurface error:%d", surface)
	}

	surface = convSurface("ダ")
	if surface != Dirt {
		t.Fatalf("convSurface error:%d", surface)
	}

	surface = convSurface("障")
	if surface != Hurdle {
		t.Fatalf("convSurface error:%d", surface)
	}

	surface = convSurface("他")
	if surface != -1 {
		t.Fatalf("convSurface error:%d", surface)
	}

}
