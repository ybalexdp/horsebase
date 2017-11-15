package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/cheggaaa/pb"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

// RaceData is race data object
type RaceData struct {
	RaceID     int           // Race ID
	Name       string        // Race name
	Course     int           // 0:Clockwise 1:Anti-Clockwise 2:Straight line
	Corner     int           // 0:Inner line 2:Outer line
	Distance   int           // Race distance
	Date       time.Time     // Day at the race
	Grade      int           // 1:Debut 2:Maidon 3:Under 500M 4:Under 1000M 5:Under 1600M 6:OPEN 7:G 8:G3 9:G2 10:G1
	Turf       string        // Turf
	RaceNumber int           // Race Number
	Day        int           // Passed days
	Surface    int           // 0:芝 1:ダート 2:障害
	Weather    int           // 0:晴 1:雨 2:雪 3:曇 4:小雨
	TrackCond  int           // 0:良 1:稍重 2:重 3:不良
	Horsenum   int           // 出走頭数
	AgeGr      int           // 0:2歳 1:3歳 2:3歳以上 3:4歳以上 4:その他
	SexGr      int           // 0:混合 1:牝馬限定
	Win        Win           // 単勝情報
	Place      Place         // 複勝情報
	Quinella   Quinella      // 馬連情報
	Exacta     Exacta        // 馬単情報
	QP         QuinellaPlace // ワイド情報
	Trio       Trio          // 3連複情報
	Trifecta   Trifecta      // 3連単情報
	Laps       []float64     // 1Fごとのタイム
}

// Horse is horse's
type Horse struct {
	HorseID    int    // netkeibaの馬ID
	Name       string //馬名
	Father     int    //父 StallionID
	FatherOfM  int    // 母父 StallionID
	FatherOfFM int    // 父母父 StallionID
	FatherOfMM int    // 母母父 StallionID
}

// RaceResultData :レース結果データ
type RaceResultData struct {
	RaceID       int       // netkeibaのレースID
	HorseID      int       // netkeibaの馬ID
	JockeyID     int       // 騎手ID
	Rank         int       // 着順
	Popularity   int       // 人気
	Odds         float64   // 単勝オッズ
	Age          int       // 年齢
	Weight       int       // 体重
	Bweight      float64   // 斤量
	Hnumber      int       // 馬番
	Wnumber      int       // 枠番
	LastThreeFur float64   // 上がり3ハロンのタイム
	Sex          int       // 性別: 0:牡 1:牝 2:騸 このレース当時の性別(牡→騸は変更があるため)
	Time         time.Time // 走破時計
	DifTime      float64   // 1着との着差
	POrder       [4]int    // 通過順(コーナー)
	Belonging    int       // 所属 0:関東 1:関西 2:地方 3:外国馬 このレース当時の所属(所属も変更あり)
}

// RaceCardData :馬柱用データ
type RaceCardData struct {
	RaceID     int    // netkeibaのレースID
	RaceNumber int    // Race Number
	Name       string // Race name
	Distance   int    // Race distance
	Grade      string // レースグレード
	Turf       string // Turf
	Day        int    // Passed days
	Surface    string // 芝/ダート/障害
	Weather    string // 晴/雨/雪/曇/小雨
	TrackCond  string // 良/稍重/重/不良
	Horsenum   int    // 出走頭数
	AgeGr      string // 2歳/3歳/3歳以上/4歳以上
	SexGr      string // 混合/牝馬限定
}

// RaceCardHorseData :出馬表内出走馬データ
type RaceCardHorseData struct {
	RaceID     int     // netkeibaのレースID
	HorseID    int     // netkeibaの馬ID
	JockeyID   int     // 騎手ID
	Popularity int     // 人気
	Odds       float64 // 単勝オッズ
	Age        int     // 年齢
	Weight     int     // 体重
	DiffWeight int     // 前走体重差
	Bweight    float64 // 斤量
	Hnumber    int     // 馬番
	Wnumber    int     // 枠番
	Sex        string  // 牡/牝/騸
	Belonging  int     // 所属 0:関東 1:関西 2:地方 3:外国馬 このレース当時の所属(所属も変更あり)
}

// Jockey :騎手データ
type Jockey struct {
	JockeyID int    // netkeibaの騎手ID
	Name     string // 騎手名
}

// Stallion :種牡馬データ
type Stallion struct {
	ID             int    // ID
	Name           string // 種牡馬名
	BloodTypeID    int    // 大系統
	SubBloodTypeID int    //小系統
}

// BloodType :血統データ
type BloodType struct {
	TypeName string // 血統名(例:ノーザンダンサー系)
}

// Win :単勝データ
type Win struct {
	Dividend   []int // 配当金
	Popularity []int // 人気
	HorseNum   []int // 馬番号
}

// Place :複勝データ
type Place struct {
	Dividend   []int // 配当金
	Popularity []int // 人気
	HorseNum   []int // 馬番号
}

// Quinella :馬連データ
type Quinella struct {
	Dividend   []int   // 配当金
	Popularity []int   // 人気
	HorseNum   [][]int // 馬番号
}

// Exacta :馬単データ
type Exacta struct {
	Dividend   []int   // 配当金
	Popularity []int   // 人気
	HorseNum   [][]int // 馬番号
}

// QuinellaPlace :ワイドデータ
type QuinellaPlace struct {
	Dividend   []int   // 配当金
	Popularity []int   // 人気
	HorseNum   [][]int // 馬番号
}

// Trio :三連複データ
type Trio struct {
	Dividend   []int   // 配当金
	Popularity []int   // 人気
	HorseNum   [][]int // 馬番号
}

// Trifecta :三連単データ
type Trifecta struct {
	Dividend   []int   // 配当金
	Popularity []int   // 人気
	HorseNum   [][]int // 馬番号
}

const (
	horseURL   = "http://db.netkeiba.com/horse/ped/"
	baseURL    = "http://db.netkeiba.com"
	racetopURL = "http://race.netkeiba.com"
	racetop    = "/?pid=race_top"
	racecard   = "/?pid=race_list"
)

const (
	AgeGrTwo = 0 + iota
	AgeGrThree
	AgeGrThreeOver
	AgeGrFourOver
)

const (
	BelongingEast = 0 + iota
	BelongingWest
	BelongingLocal
	BelongingForeign
)

const (
	ConditionGood = 0 + iota
	ConditionYielding
	ConditionSoft
	ConditionHeavy
)

const (
	InnerCourse = 0 + iota
	OuterCourse
)

const (
	RightTurn = 0 + iota
	LeftTurn
	Straight
)

const (
	GradeDebut = 0 + iota
	GradeNoWin
	Grade500
	Grade1000
	Grade1600
	GradeOP
	GradeG
	GradeG3
	GradeG2
	GradeG1
)

const (
	Male = 0 + iota
	Female
	Gelding
)

const (
	Turf = 0 + iota
	Dirt
	Hurdle
)

const (
	Sunny = 0 + iota
	Rainy
	Snowy
	Cloudy
	Drizzle
)

const (
	MixedRace = 0 + iota
	FemaleRace
)

// GetRaceHTML :
// file/racelist.txtに記載されているURLからHTML取得
func (hb *Horsebase) GetRaceHTML() error {
	var fpr *os.File

	fpr, err := os.Open(hb.dir + "/file/racelist.txt")
	if err != nil {
		return err
	}
	defer fpr.Close()

	// ./htmlフォルダ有無確認
	_, err = os.Stat(hb.dir + "/html")
	if err != nil {
		if err := os.Mkdir(hb.dir+"/html", 0777); err != nil {
			return err
		}
	}

	// ./html/raceフォルダ有無確認
	_, err = os.Stat(hb.Config.RaceHtmlPath)
	if err != nil {
		if err := os.Mkdir(hb.Config.RaceHtmlPath, 0777); err != nil {
			return err
		}
	}

	_, err = os.Stat(hb.Config.HorseHtmlPath)
	if err != nil {
		if err := os.Mkdir(hb.Config.HorseHtmlPath, 0777); err != nil {
			return err
		}
	}

	reader := bufio.NewReader(fpr)

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}

		if len(line) > 2 {
			raceID := getRaceIDfromHTML((string)(line))
			// 下記2つはnetkeibaのバグ
			if (raceID == "200808020398") || (raceID == "200808020399") {
				continue
			}

			//ファイル未取得の場合
			if !hb.raceExistenceCheck(raceID) {
				getHTML((string)(line), raceID, hb.Config.RaceHtmlPath)
			}
		}
	}
	return err
}

func (hb *Horsebase) GetRacecardHTML() error {
	var fpr *os.File

	fpr, err := os.Open(hb.dir + "/file/racecarddata.txt")
	if err != nil {
		return err
	}
	defer fpr.Close()

	// ./htmlフォルダ有無確認
	_, err = os.Stat(hb.dir + "/html")
	if err != nil {
		if err := os.Mkdir(hb.dir+"/html", 0777); err != nil {
			return err
		}
	}

	// ./html/raceフォルダ有無確認
	_, err = os.Stat(hb.Config.CardHtmlPath)
	if err != nil {
		if err := os.Mkdir(hb.Config.CardHtmlPath, 0777); err != nil {
			return err
		}
	}

	_, err = os.Stat(hb.Config.HorseHtmlPath)
	if err != nil {
		if err := os.Mkdir(hb.Config.HorseHtmlPath, 0777); err != nil {
			return err
		}
	}

	reader := bufio.NewReader(fpr)

	for {
		line, _, err := reader.ReadLine()
		if err == io.EOF {
			break
		} else if err != nil {
			return err
		}
		raceID := getRacecardIDfromHTML((string)(line))

		getHTML((string)(line), raceID, hb.Config.CardHtmlPath)
	}

	return err
}

func (hb *Horsebase) MakeRaceURLList() error {

	var racelist []string
	var raceURLlist []string

	file := hb.dir + "/file/racelist.txt"

	fp, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer fp.Close()

	writer := bufio.NewWriter(fp)

	// netkeiba.comからレースのURL一覧を取得
	racelist, err = hb.getRaceList()
	if err != nil {
		return err
	}

	for _, racelistURL := range racelist {
		raceURLlist, err = getRaceURL(racelistURL, raceURLlist)
		if err != nil {
			return err
		}
	}

	_, err = os.Stat(file)
	if err != nil {
		if err = os.Remove(file); err != nil {
			return err
		}
	}

	for _, raceURL := range raceURLlist {
		_, err = writer.WriteString(raceURL + "\n")
		if err != nil {
			return err
		}
	}
	writer.Flush()

	return err

}

// 出馬表作成用URLを取得しracecard.txtに保存する
func (hb *Horsebase) MakeRacecardURLList() error {
	var racelist []string

	file := hb.dir + "/file/racecard.txt"

	fp, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE, 0644)
	if err != nil {
		return err
	}

	defer fp.Close()

	writer := bufio.NewWriter(fp)

	// netkeiba.comからレースのURL一覧を取得
	racelist, err = hb.getRacecardList()
	if err != nil {
		return err
	}

	_, err = os.Stat(file)
	if err != nil {
		if err = os.Remove(file); err != nil {
			return err
		}
	}

	for _, raceURL := range racelist {
		_, err = writer.WriteString(raceURL + "\n")
		if err != nil {
			return err
		}
	}
	writer.Flush()

	return err
}

// 競走馬データを登録する
func (hb *Horsebase) RegistHorseData() error {

	var horse Horse
	var stallion [4]Stallion

	// 取得済みのHTML読み込み
	files, err := ioutil.ReadDir(hb.Config.HorseHtmlPath)
	if err != nil {
		return err
	}

	count := len(files)
	bar := pb.StartNew(count)

	hb.DbInfo, err = hb.DbInfo.New()
	if err != nil {
		PrintError(hb.Stderr, "%s", err)
		return err
	}
	defer hb.DbInfo.db.Close()

	for _, file := range files {
		bar.Increment()

		fp, err := os.Open(hb.Config.HorseHtmlPath + file.Name())
		if err != nil {
			return err
		}
		defer fp.Close()

		doc, err := goquery.NewDocumentFromReader(fp)
		if err != nil {
			return err
		}

		horseID := strings.Split(file.Name(), ".")[0]

		horse.HorseID, _ = strconv.Atoi(horseID)

		horseCheck, err := hb.DbInfo.HorseExistenceCheck(horseID)
		if err != nil {
			return err
		}

		if horseCheck {
			fp.Close()
			continue
		}

		s := doc.Find("h1").First()
		horse.Name = strings.TrimSpace(s.Text())

		mFlag := false
		fmFlag := false

		doc.Find("tr > td.b_ml").Each(func(_ int, s *goquery.Selection) {
			rowspan, _ := s.Attr("rowspan")
			_, width := s.Attr("width")

			s = s.Children().First()

			if rowspan == "16" {
				stallion[0].Name = strings.Split(strings.TrimSpace(s.Text()), "\n")[0]
				hb.DbInfo.InsertStallion(stallion[0].Name)
				horse.Father, _ = hb.DbInfo.GetId("stallion", stallion[0].Name)
			}

			if rowspan == "8" {
				if mFlag {
					stallion[1].Name = strings.Split(strings.TrimSpace(s.Text()), "\n")[0]
					hb.DbInfo.InsertStallion(stallion[1].Name)
					horse.FatherOfM, _ = hb.DbInfo.GetId("stallion", stallion[1].Name)
				} else {
					mFlag = true
				}
			}

			if !width && rowspan == "4" {
				if fmFlag {
					stallion[3].Name = strings.Split(strings.TrimSpace(s.Text()), "\n")[0]
					hb.DbInfo.InsertStallion(stallion[3].Name)
					horse.FatherOfMM, _ = hb.DbInfo.GetId("stallion", stallion[3].Name)
				} else {
					stallion[2].Name = strings.Split(strings.TrimSpace(s.Text()), "\n")[0]
					err = hb.DbInfo.InsertStallion(stallion[2].Name)
					horse.FatherOfFM, _ = hb.DbInfo.GetId("stallion", stallion[2].Name)
					fmFlag = true
				}
			}

		})
		if err != nil {
			return err
		}
		err = hb.DbInfo.UpdateHorse(horse)
		if err != nil {
			return err
		}
		fp.Close()

	}

	bar.FinishPrint("Registed Horse Data")

	return err
}

func (hb *Horsebase) RegistRaceData() error {

	var data []string
	var horse Horse
	var result RaceResultData
	var ftime time.Time
	var jockey Jockey

	files, err := ioutil.ReadDir(hb.Config.RaceHtmlPath)
	if err != nil {
		return err
	}
	count := len(files)
	bar := pb.StartNew(count)

	for _, file := range files {

		var racedata RaceData
		bar.Increment()

		hb.DbInfo, err = hb.DbInfo.New()
		if err != nil {
			PrintError(hb.Stderr, "%s", err)
			return err
		}
		defer hb.DbInfo.db.Close()

		//レースID取得
		raceID := strings.Split(file.Name(), ".")[0]
		if len(raceID) != 12 {
			hb.DbInfo.db.Close()
			continue
		}
		racedata.RaceID, _ = strconv.Atoi(raceID)

		// 既に登録済みのレースデータであれば解析しない
		raceCheck, err := hb.DbInfo.RaceExistenceCheck(raceID)
		if err != nil {
			return err
		}

		if raceCheck {
			hb.DbInfo.db.Close()
			continue
		}

		fp, err := os.Open(hb.Config.RaceHtmlPath + file.Name())
		if err != nil {
			return err
		}
		defer fp.Close()

		doc, err := goquery.NewDocumentFromReader(fp)
		if err != nil {
			return err
		}

		//fmt.Println("Start:", racedata.RaceID)

		//レース番号
		s := doc.Find("dt").First()
		racedata.RaceNumber, _ = strconv.Atoi(s.Text()[1:3])

		//レース名
		s = doc.Find("title").First()
		title := strings.Split(s.Text(), "｜")
		racedata.Name = title[0]

		//レース開催日
		racedata.Date = getRaceDate(title[1])

		//レースクラス(重賞の場合)
		s = doc.Find("h1").First()
		var gradeFlag bool
		gradeFlag, racedata.Grade = getRaceGrade(s.Text())

		//芝 or ダート or 障害
		s = doc.Find("diary_snap_cut > span").First()
		data = strings.Split(s.Text(), "/")
		racedata.Surface = convSurface(data[0][:3])

		//距離,コーナー,コース
		index := strings.Index(data[0], "m")
		racedata.Distance, _ = strconv.Atoi(data[0][index-4 : index])

		// 障害でなければ
		if racedata.Surface != 2 {
			racedata.getCourseInfo(data)
		}

		// 天気
		weather := strings.Split(data[1], ":")[1][1:]
		racedata.Weather = convWeather(strings.TrimSpace(weather))

		// 馬場状態
		cond := strings.Split(data[2], ":")[1][1:]
		racedata.TrackCond = convCond(strings.TrimSpace(cond))

		// 競馬場
		s = doc.Find("p.smalltxt").First()
		data = strings.Split(s.Text(), "回")
		racedata.Turf = data[1][:6]

		// 第何日の開催か
		racedata.Day, _ = strconv.Atoi(data[1][6:7])

		// レース年齢
		grade := getGradeStr(data)
		racedata.AgeGr, index = getAgeGr(grade)

		// レースクラス(非重賞)
		if !gradeFlag {
			racedata.Grade = convGrade(strings.Split(grade[index:], "  ")[0])
		}

		racedata.SexGr = getSexGr(grade)

		s = doc.Find("td").First()
		racedata.Horsenum = 0

		// レースデータは途中でエラーになると次の実行で不整合となるため
		// 失敗したらロールバックできるようにする
		hb.DbInfo.tx, err = hb.DbInfo.db.Begin()
		if err != nil {
			return err
		}

		defer func() {
			if err := recover(); err != nil {
				hb.DbInfo.tx.Rollback()
			}
		}()

		err = hb.DbInfo.InsertRaceData(racedata)
		if err != nil {
			return err
		}

		result.RaceID = racedata.RaceID

		for {
			// 着順
			result.Rank, err = strconv.Atoi(s.Text())
			if err != nil {

				// 降着の場合は着順を抜き取って続ける
				if strings.Contains(s.Text(), "降") {
					result.Rank, _ = strconv.Atoi(strings.Split(s.Text(), "(")[0])
				} else {
					// 除外や中止の場合ループを抜ける
					err = hb.DbInfo.UpdateHorseNum(racedata)
					if err != nil {
						return err
					}
					break
				}
			}
			racedata.Horsenum++

			// 枠番
			s = s.Next()
			result.Wnumber, _ = strconv.Atoi(s.Text())

			// 馬番
			s = s.Next()
			result.Hnumber, _ = strconv.Atoi(s.Text())

			// horse id(netkeibaの)
			s = s.Next()
			attr, _ := s.Children().Attr("href")
			horseID := strings.Split(attr, "/")[2]
			horse.HorseID, _ = strconv.Atoi(horseID)

			// 馬名
			horse.Name = strings.TrimRight(s.Text()[1:], "\n")
			result.HorseID = horse.HorseID

			// データ未登録の馬は馬データの登録
			horseCheck, err := hb.DbInfo.HorseExistenceCheck(horseID)
			if err != nil {
				return err
			}

			if !horseCheck {
				err = hb.DbInfo.InsertHorse(horse)
				if err != nil {
					return err
				}

				// getHorseData内でHTTP GETするため
				// インターバルをおく
				time.Sleep(3000 * time.Millisecond)
				err = hb.getHorseData(horse.HorseID)
				if err != nil {
					return err
				}
			}

			// 性別
			s = s.Next()
			result.Sex = convSex(s.Text()[:3])
			// 年齢
			result.Age, _ = strconv.Atoi(s.Text()[3:])

			// 斤量
			s = s.Next()
			result.Bweight, _ = strconv.ParseFloat(s.Text(), 64)

			// 騎手名
			s = s.Next()
			jockey.Name, _ = s.Children().Attr("title")
			hb.DbInfo.InsertJockey(jockey)

			result.JockeyID, err = hb.DbInfo.GetId("jockey", jockey.Name)
			if err != nil {
				return err
			}

			// 走破タイム
			s = s.Next()
			result.Time = racedata.getRaceTime(s.Text())

			// 1着のタイムを保持して1着とのタイム差の計算に使用
			if result.Rank == 1 {
				ftime = result.Time
			}

			// 着差
			result.DifTime = calcDifTime(result, ftime)

			// 通過順位
			s = s.Next().Next().Next() // タイム指数不要
			result.POrder = getPassOrder(s.Text())

			// ラスト3F
			s = s.Next()
			result.LastThreeFur, _ = strconv.ParseFloat(s.Text(), 64)

			// 単勝オッズ
			s = s.Next()
			result.Odds, _ = strconv.ParseFloat(s.Text(), 64)

			// 人気
			s = s.Next()
			result.Popularity, _ = strconv.Atoi(s.Text())

			// 馬体重
			s = s.Next()
			result.Weight, _ = strconv.Atoi(strings.Split(s.Text(), "(")[0])

			// 所属
			s = s.Next().Next().Next().Next()
			result.Belonging = convBelonging(s.Text()[2:5])

			err = hb.DbInfo.InsertRaceresult(result)
			if err != nil {
				return err
			}

			s = s.Parent().Next().Children().First()

			attr, _ = s.Attr("class")
			if attr != "txt_r" {
				err = hb.DbInfo.UpdateHorseNum(racedata)
				if err != nil {
					return err
				}
				break
			}
		}

		// 単勝馬番
		s = doc.Find("th.tan").First().Next()
		racedata.Win.HorseNum = getDividendInfo(s.Text())

		// 単勝配当金
		s = s.Next()
		racedata.Win.Dividend = getDividendInfo(s.Text())

		// 単勝人気
		s = s.Next()
		racedata.Win.Popularity = getDividendInfo(s.Text())

		// 複勝馬番
		s = doc.Find("th.fuku").First().Next()
		racedata.Place.HorseNum = getDividendInfo(s.Text())

		// 複勝配当金
		s = s.Next()
		racedata.Place.Dividend = getDividendInfo(s.Text())

		// 複勝人気
		s = s.Next()
		racedata.Place.Popularity = getDividendInfo(s.Text())

		// 馬連馬番
		s = doc.Find("th.uren").First().Next()
		racedata.Quinella.HorseNum = getHorseNum(s.Text(), "-")
		// 馬連配当金
		s = s.Next()
		racedata.Quinella.Dividend = getDividendInfo(s.Text())

		// 馬連人気
		s = s.Next()
		racedata.Quinella.Popularity = getDividendInfo(s.Text())

		// ワイド馬番
		s = doc.Find("th.wide").First().Next()
		if strings.Contains(s.Text(), "-") {
			racedata.QP.HorseNum = getHorseNum(s.Text(), "-")

			// ワイド配当金
			s = s.Next()
			racedata.QP.Dividend = getDividendInfo(s.Text())

			// ワイド人気
			s = s.Next()
			racedata.QP.Popularity = getDividendInfo(s.Text())
		}

		// 馬単馬番
		s = doc.Find("th.utan").First().Next()
		if strings.Contains(s.Text(), "→") {
			racedata.Exacta.HorseNum = getHorseNum(s.Text(), "→")

			// 馬単配当金
			s = s.Next()
			racedata.Exacta.Dividend = getDividendInfo(s.Text())

			// 馬単人気
			s = s.Next()
			racedata.Exacta.Popularity = getDividendInfo(s.Text())
		}

		// 三連複馬番
		s = doc.Find("th.sanfuku").First().Next()
		if strings.Contains(s.Text(), "-") {
			racedata.Trio.HorseNum = getHorseNum(s.Text(), "-")

			// 三連複配当金
			s = s.Next()
			racedata.Trio.Dividend = getDividendInfo(s.Text())

			// 三連複人気
			s = s.Next()
			racedata.Trio.Popularity = getDividendInfo(s.Text())
		}

		// 三連単馬番
		s = doc.Find("th.santan").First().Next()

		if strings.Contains(s.Text(), "→") {
			racedata.Trifecta.HorseNum = getHorseNum(s.Text(), "→")
			// 三連単配当金
			s = s.Next()
			racedata.Trifecta.Dividend = getDividendInfo(s.Text())

			// 三連単人気
			s = s.Next()
			racedata.Trifecta.Popularity = getDividendInfo(s.Text())
		}

		s = doc.Find("td.race_lap_cell").First()
		racedata.Laps = getLaps(s.Text())

		hb.registDividendInfo(racedata)
		err = hb.DbInfo.tx.Commit()
		if err != nil {
			return err
		}
		hb.DbInfo.db.Close()
		fp.Close()

	}
	bar.FinishPrint("Registed Race Data")

	return err
}

func (hb *Horsebase) RegistRacecardData() error {

	var err error
	var rchd RaceCardHorseData

	files, err := ioutil.ReadDir(hb.Config.CardHtmlPath)
	if err != nil {
		return err
	}

	for _, file := range files {

		var racecarddata RaceCardData

		hb.DbInfo, err = hb.DbInfo.New()
		if err != nil {
			PrintError(hb.Stderr, "%s", err)
			return err
		}
		defer hb.DbInfo.db.Close()

		//レースID取得
		raceID := strings.Split(file.Name(), ".")[0]
		if len(raceID) != 12 {
			hb.DbInfo.db.Close()
			continue
		}
		racecarddata.RaceID, _ = strconv.Atoi(raceID)

		fmt.Println("RaceID:", racecarddata.RaceID)

		fp, err := os.Open(hb.Config.CardHtmlPath + file.Name())
		if err != nil {
			return err
		}
		defer fp.Close()

		doc, err := goquery.NewDocumentFromReader(fp)
		if err != nil {
			return err
		}

		//fmt.Println("Start:", racedata.RaceID)

		//レース番号
		s := doc.Find("dl.racedata > dt").First()
		//fmt.Println("tttt:" + strings.Split(s.Text(), "R")[0][1:])
		racecarddata.RaceNumber, _ = strconv.Atoi(strings.Split(s.Text(), "R")[0][1:])

		fmt.Println("RaceNumber:", racecarddata.RaceNumber)

		s = doc.Find("dd > h1").First()
		racecarddata.Name = s.Text()

		fmt.Println("RaceName:" + racecarddata.Name)

		s = doc.Find("p > span").First()

		data := strings.Split(s.Text(), "/")
		racecarddata.Surface = data[0][:3]
		fmt.Println("Surface:" + racecarddata.Surface)

		//距離,コーナー,コース
		racecarddata.Distance, _ = strconv.Atoi(data[0][3:7])
		fmt.Println("Distance:", racecarddata.Distance)

		// 天気
		weather := strings.Split(data[1], "：")
		if len(weather) > 1 {
			racecarddata.Weather = strings.TrimSpace(weather[1])
		} else {
			racecarddata.Weather = "不明"
		}
		fmt.Println("Weather:" + racecarddata.Weather)

		// 馬場状態
		cond := strings.Split(data[2], "：")
		if len(cond) > 1 {
			racecarddata.TrackCond = strings.TrimSpace(cond[1])
		} else {
			racecarddata.TrackCond = "不明"
		}
		fmt.Println("TrackCond:" + racecarddata.TrackCond)

		// 競馬場
		s = doc.Find("div.race_otherdata").Children().First()
		other_data := strings.Replace(s.Text(), " ", ",", -1)
		data = strings.Split(other_data, "回")
		racecarddata.Turf = data[1][0:6]
		fmt.Println("Turf:" + racecarddata.Turf)

		racecarddata.Day, _ = strconv.Atoi(data[1][6:7])
		fmt.Println("Day:", racecarddata.Day)

		// レース年齢
		data = strings.Split(other_data, ",")
		fmt.Println("data:" + data[0])
		if strings.Contains(data[1], "障害") {
			racecarddata.AgeGr = strings.Split(data[1], "障害")[1]
		} else {
			racecarddata.AgeGr = data[1]
		}
		fmt.Println("AgeGr:" + racecarddata.AgeGr)

		// レースクラス
		//gradecheck := strings.Split(data[1], "歳")
		fmt.Println("other_data:" + other_data)
		if len(data[2]) > 2 {
			attr, _ := doc.Find("dd > h1 > img").First().Attr("src")
			if len(attr) < 1 {
				if strings.Contains(data[2], "上") {
					racecarddata.Grade = strings.Split(data[2], "上")[1]
				} else {
					racecarddata.Grade = strings.Split(data[2], "歳")[1]
				}
			} else {
				racecarddata.Grade = strings.Split(attr, "_")[2]
			}
		} else {
			if strings.Contains(racecarddata.Name, "障害") {
				racecarddata.Grade = strings.Split(racecarddata.Name, "障害")[1]
			} else if strings.Contains(racecarddata.Name, "上") {
				racecarddata.Grade = strings.Split(racecarddata.Name, "上")[1]
			} else if strings.Contains(racecarddata.Name, "歳") {
				racecarddata.Grade = strings.Split(racecarddata.Name, "歳")[1]
			} else {
				return fmt.Errorf("Race Grade is unknown")
			}
		}
		fmt.Println("Grade:" + racecarddata.Grade)

		if strings.Contains(s.Next().Text(), "牝") {
			racecarddata.SexGr = "牝"
		} else {
			racecarddata.SexGr = "混"
		}
		fmt.Println("SexGr:" + racecarddata.SexGr)

		rchd.Hnumber = 0
		var wakunum []int
		i := 0
		doc.Find("tr.bml1").Each(func(_ int, s *goquery.Selection) {
			s = s.First().Children().First()
			waku, _ := strconv.Atoi(s.Text())
			wakunum = append(wakunum, waku)
		})

		doc.Find("div > a").Each(func(_ int, s *goquery.Selection) {
			horseURL, _ := s.First().Attr("href")
			if !strings.Contains(horseURL, "javascript") {
				rchd.HorseID, _ = strconv.Atoi(strings.Split(horseURL, "/")[4])
				fmt.Println("HorseID:", rchd.HorseID)
				racecarddata.Horsenum++

				// 性別
				s = s.Parent().Parent().Next()
				info := s.Text()
				rchd.Sex = info[0:3]
				fmt.Println("Sex:" + rchd.Sex)

				// 年齢
				rchd.Age, _ = strconv.Atoi(info[3:])
				fmt.Println("Age:", rchd.Age)

				// 斤量
				s = s.Next()
				rchd.Bweight, _ = strconv.ParseFloat(s.Text(), 64)
				fmt.Println("Bweight:", rchd.Bweight)

				// 騎手
				s = s.Next().Children().First()
				attr, _ := s.Attr("href")
				jockeyID := strings.Split(attr, "/")[4]
				rchd.JockeyID, _ = strconv.Atoi(jockeyID)
				fmt.Println("JockeyID:", rchd.JockeyID)

				// 体重
				s = s.Parent().Next().Next()
				weight := strings.Split(s.Text(), "(")
				rchd.Weight, err = strconv.Atoi(weight[0])
				if err == nil {
					fmt.Println("Weight:", rchd.Weight)

					// 前走体重差
					diff_weight := strings.Split(weight[1], ")")[0]
					rchd.DiffWeight, _ = strconv.Atoi(diff_weight)
					fmt.Println("DiffWeight:", rchd.DiffWeight)

				}

				// オッズ
				s = s.Next()
				rchd.Odds, _ = strconv.ParseFloat(s.Text(), 64)
				fmt.Println("Odds:", rchd.Odds)

				// 人気
				s = s.Next()
				rchd.Popularity, _ = strconv.Atoi(s.Text())
				fmt.Println("Popularty:", rchd.Popularity)

				rchd.Wnumber = wakunum[i]
				// 枠番
				fmt.Println("WakuNumber:", rchd.Wnumber)
				i++
				// 馬番

				fmt.Println("HorseNumber:", rchd.Hnumber)

			}

		})

	}

	return err
}

// 1着馬との着差を算出する
func calcDifTime(result RaceResultData, ftime time.Time) float64 {
	var diftime float64

	diftimeStr := result.Time.Sub(ftime).String()

	if strings.Contains(diftimeStr, "ms") {
		i, _ := strconv.Atoi(strings.Split(diftimeStr, "ms")[0][:1])
		diftime = (float64)(i) / 10
	} else {
		diftime, _ = strconv.ParseFloat(strings.Split(diftimeStr, "s")[0], 64)
	}

	return diftime

}

func (hb *Horsebase) checkOldestData(date string) bool {
	target, _ := strconv.Atoi(strings.Split(date, "=")[2])
	result := false

	if target > hb.Config.OldestDate {
		result = true
	}

	return result
}

func convAgeGr(agegr string) int {
	switch agegr {
	case "2歳":
		return AgeGrTwo
	case "3歳":
		return AgeGrThree
	case "3歳以上":
		return AgeGrThreeOver
	case "4歳以上":
		return AgeGrFourOver
	default:
		return -1
	}
}

func convBelonging(bel string) int {
	switch bel {
	case "東":
		return BelongingEast
	case "西":
		return BelongingWest
	case "地":
		return BelongingLocal
	case "外":
		return BelongingForeign
	default:
		return -1
	}
}

func convCond(cond string) int {
	switch cond {
	case "良":
		return ConditionGood
	case "稍重":
		return ConditionYielding
	case "重":
		return ConditionSoft
	case "不良":
		return ConditionHeavy
	default:
		return -1
	}
}

func convCorner(corner string) int {
	switch corner {
	case "内":
		return InnerCourse
	case "外":
		return OuterCourse
	default:
		return -1
	}
}

func convCourse(course string) int {
	switch course {
	case "右":
		return RightTurn
	case "左":
		return LeftTurn
	case "直線":
		return Straight
	default:
		return -1
	}
}

func convGrade(grade string) int {
	if strings.Contains(grade, "J.") {
		grade = strings.Split(grade, ".")[1]
	}
	switch grade {
	case "新馬":
		return GradeDebut
	case "未勝利":
		return GradeNoWin
	case "500万下":
		return Grade500
	case "1000万下":
		return Grade1000
	case "1600万下":
		return Grade1600
	case "オープン":
		return GradeOP
	case "G":
		return GradeG
	case "G3":
		return GradeG3
	case "G2":
		return GradeG2
	case "G1":
		return GradeG1
	default:
		return -1
	}
}

func convSex(sex string) int {
	switch sex {
	case "牡":
		return Male
	case "牝":
		return Female
	case "セ":
		return Gelding
	default:
		return -1
	}
}

func convSurface(surface string) int {
	switch surface {
	case "芝":
		return Turf
	case "ダ":
		return Dirt
	case "障":
		return Hurdle
	default:
		return -1
	}
}

func convWeather(weather string) int {
	switch weather {
	case "晴":
		return Sunny
	case "雨":
		return Rainy
	case "雪":
		return Snowy
	case "曇":
		return Cloudy
	case "小雨":
		return Drizzle
	default:
		return -1
	}
}

func getGradeStr(data []string) string {
	var grade string

	if strings.Contains(data[1][7:], "障害") {
		grade = strings.Split(data[1][7:], "害")[1]
	} else {
		grade = strings.Split(data[1][7:], "系")[1]
	}

	return grade
}

func getAgeGr(grade string) (int, int) {

	if strings.Contains(grade, "上") {
		ageGr := convAgeGr(grade[:10])
		return ageGr, 10
	}
	ageGr := convAgeGr(grade[:4])
	return ageGr, 4

}

func (racedata *RaceData) getCourseInfo(data []string) {

	if strings.Contains(data[0], "直線") {
		racedata.Course = convCourse(data[0][3:9])
	} else if strings.Contains(data[0], "外") || strings.Contains(data[0], "内") {
		racedata.Corner = convCorner(data[0][7:10])
		racedata.Course = convCourse(data[0][3:6])
	}

}

func (hb *Horsebase) getHorseData(id int) error {
	horseID := strconv.Itoa(id)
	url := horseURL + horseID
	err := getHTML(url, horseID, hb.Config.HorseHtmlPath)
	return err
}

/*
  HTMLファイル取得
*/
func getHTML(url string, id string, htmlPath string) error {

	file := htmlPath + id + ".html"

	fp, err := os.OpenFile(file, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer fp.Close()

	writer := bufio.NewWriter(fp)

	html, err := getResBody(url)
	if err != nil {
		return err
	}

	html = strings.Replace(html, "<br />", "/", -1)
	_, err = writer.WriteString(html)
	if err != nil {
		return err
	}
	writer.Flush()

	return err
}

/*
  馬券情報入手
*/
func getDividendInfo(text string) []int {
	var divs []int
	div := strings.Split(text, "/")
	for _, v := range div {
		v = strings.Replace(v, ",", "", -1)
		i, _ := strconv.Atoi(v)
		divs = append(divs, i)
	}
	return divs
}

func getDividendStr(text string) []string {
	var divs []string
	div := strings.Split(text, "/")
	for _, v := range div {
		divs = append(divs, v)
	}
	return divs
}

func getHorseNum(text string, symbol string) [][]int {
	var horseNums [][]int
	ticket := getDividendStr(text)
	for _, n := range ticket {
		horseNums = append(horseNums, splitHorseNum(n, symbol))
	}

	return horseNums
}

func getLaps(text string) []float64 {
	var laps []float64
	lap := strings.Split(text, "-")
	for _, v := range lap {
		f, _ := strconv.ParseFloat(strings.TrimSpace(v), 64)
		laps = append(laps, f)
	}
	return laps
}

func getPassOrder(text string) [4]int {
	var passOrder [4]int
	porderStr := strings.Split(text, "-")

	for i, o := range porderStr {
		passOrder[i], _ = strconv.Atoi(o)
	}
	return passOrder
}

func getRaceDate(text string) time.Time {
	date := strings.Split(text, "|")[0]
	year, _ := strconv.Atoi(date[:4])
	month, _ := strconv.Atoi(date[7:9])
	day, _ := strconv.Atoi(date[12:14])

	return time.Date(year, (time.Month)(month), day, 0, 0, 0, 0, time.Local)
}

func getRaceIDfromHTML(url string) string {
	return strings.Split(url, "/")[4]
}

func getRacecardIDfromHTML(url string) string {
	return strings.Split(url, "=")[2][1:13]
}

func (hb *Horsebase) getRaceList() ([]string, error) {

	var racelist []string
	doc, err := goquery.NewDocument(baseURL + racetop)
	if err != nil {
		return nil, err
	}

	s := doc.Find("li.rev").First().Children().Next()
	attr, _ := s.Attr("href")

	racelist, err = hb.addRaceList(baseURL+attr, racelist)
	if err != nil {
		return nil, err
	}

	return racelist, nil
}

func (hb *Horsebase) getRacecardList() ([]string, error) {

	var racelist []string
	doc, err := goquery.NewDocument(racetopURL + racecard)
	if err != nil {
		return nil, err
	}

	doc.Find("dt > a").Each(func(_ int, s *goquery.Selection) {
		attr, _ := s.Attr("href")

		if !isOtherURI(attr) {

			if isRaceOldURI(attr) {
				attr = convRacecardURI(attr)
			}

			racelist = append(racelist, racetopURL+attr)

		}

	})

	return racelist, nil
}

func isRaceOldURI(uri string) bool {
	urlsplit := strings.Split(uri, "&")
	if len(urlsplit) == 3 {
		return true
	}
	return false
}

func isOtherURI(uri string) bool {
	if strings.Contains(uri, "http://") {
		return true
	}
	return false
}

func convRacecardURI(uri string) string {
	urisplit := strings.Split(uri, "&")
	uri = urisplit[0] + "_old&" + urisplit[1]
	return uri
}

func (hb *Horsebase) addRaceList(url string, racelist []string) ([]string, error) {

	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	s := doc.Find("li.rev").First().Children().Next()
	attr, _ := s.Attr("href")

	doc.Find("a").Each(func(_ int, s *goquery.Selection) {
		listURL, _ := s.Attr("href")
		if strings.Contains(listURL, "/race/list/") {
			racelist = append(racelist, baseURL+listURL)
		}
		racelist = removeDuplicate(racelist)
	})

	if hb.checkOldestData(attr) {
		racelist, err = hb.addRaceList(baseURL+attr, racelist)
		if err != nil {
			return nil, err
		}
	}

	return racelist, nil
}

func getRaceGrade(text string) (bool, int) {
	gradecheck := strings.Split(text, "(")
	if len(gradecheck) > 1 {
		// 天皇賞(春) or (秋)の場合
		if len(gradecheck) == 3 {
			grade := convGrade(strings.Split(gradecheck[2], ")")[0])
			return true, grade
		}
		grade := convGrade(strings.Split(gradecheck[1], ")")[0])
		return true, grade
	}
	return false, -1
}

func (racedata RaceData) getRaceTime(text string) time.Time {
	min, _ := strconv.Atoi(strings.Split(text, ":")[0])
	sec, _ := strconv.Atoi(strings.Split(text, ":")[1][:2])
	ns, _ := strconv.Atoi(strings.Split(text, ".")[1])
	ns = ns * 100000000
	return time.Date(racedata.Date.Year(), racedata.Date.Month(), racedata.Date.Day(), 0, min, sec, ns, time.Local)
}

func getRaceURL(url string, raceURLlist []string) ([]string, error) {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		return nil, err
	}

	doc.Find("dd > a").Each(func(_ int, s *goquery.Selection) {
		raceURL, _ := s.Attr("href")
		if strings.Contains(raceURL, "/race/") && !strings.Contains(raceURL, "/pay/") && !strings.Contains(raceURL, "/sum/") {
			raceURLlist = append(raceURLlist, baseURL+raceURL)
		}
		raceURLlist = removeDuplicate(raceURLlist)
	})

	return raceURLlist, nil
}

func getResBody(url string) (string, error) {
	var html string
	res, err := http.Get(url)
	if err != nil {
		return html, err
	}
	defer res.Body.Close()

	utfbody := transform.NewReader(bufio.NewReader(res.Body), japanese.EUCJP.NewDecoder())
	body, err := ioutil.ReadAll(utfbody)
	if err != nil {
		return html, err
	}
	buf := bytes.NewBuffer(body)
	html = buf.String()
	if err != nil {
		return html, err
	}
	return html, nil
}

func getSexGr(data string) int {
	if strings.Contains(data, "牝") {
		return FemaleRace
	} else {
		return MixedRace
	}
}

/*
func (hb *Horsebase) getStallionInfo(text string, stallion *Stallion) error {
	var err error
	stallion.Name = strings.Split(strings.TrimSpace(text), "\n")[0]
	stallion.Id, err = hb.DbInfo.GetStallionId(stallion.Name)
	return err
}
*/

/*
  すでにHTMLが取得済みかどうかを確認する
  TODO 将来的にはDBから確認
*/
func (hb *Horsebase) raceExistenceCheck(raceID string) bool {
	file := hb.Config.RaceHtmlPath + raceID + ".html"

	_, err := os.Stat(file)
	return err == nil
}

func removeDuplicate(args []string) []string {
	results := make([]string, 0, len(args))
	encountered := map[string]bool{}
	for i := range args {
		if !encountered[args[i]] {
			encountered[args[i]] = true
			results = append(results, args[i])
		}
	}
	return results
}

func splitHorseNum(ticket string, symbol string) []int {
	var horseNum []int
	horseNumStr := strings.Split(ticket, symbol)
	for _, numStr := range horseNumStr {
		num, _ := strconv.Atoi(strings.TrimSpace(numStr))
		horseNum = append(horseNum, num)
	}
	return horseNum
}

/*
func (hb *Horsebase) regStallionInfo(name string) error {
	err := hb.DbInfo.InsertStallion(name)
	return err
}
*/

func (hb *Horsebase) registDividendInfo(rd RaceData) {
	for i := range rd.Win.HorseNum {
		hb.DbInfo.InsertWinData(rd, i)
	}

	for i := range rd.Place.HorseNum {
		hb.DbInfo.InsertPlaceData(rd, i)
	}

	for i := range rd.Quinella.HorseNum {
		hb.DbInfo.InsertQuinellaData(rd, i)
	}

	for i := range rd.Exacta.HorseNum {
		hb.DbInfo.InsertExactaData(rd, i)
	}

	for i := range rd.QP.HorseNum {
		hb.DbInfo.InsertQPData(rd, i)
	}

	for i := range rd.Trio.HorseNum {
		hb.DbInfo.InsertTrioData(rd, i)
	}

	for i := range rd.Trifecta.HorseNum {
		hb.DbInfo.InsertTrifectaData(rd, i)
	}

}
