package main

import (
	"database/sql"
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	layout  = "2006-01-02"
	tlayout = "04:05.0"
)

// Horsebase :horsebase.tomlの[horsebase]のパラメータ値
// を保持したオブジェクト
type Horsebase struct {
	Config Config `toml:"config"`
	DbInfo HBDB   `toml:"db"`
	Stdout io.Writer
	Stderr io.Writer
}

// Config :
// horsebase.tomlで定義した設定値
// を保持したオブジェクト
type Config struct {
	RaceHtmlPath  string `toml:"race_html_path"`
	HorseHtmlPath string `toml:"horse_html_path"`
	OldestDate    int    `toml:"oldest_date"`
}

type HBDB struct {
	DbUser string `toml:"dbuser"`
	DbPass string `toml:"dbpass"`
	db     *sql.DB
	tx     *sql.Tx
}

func (hb *Horsebase) New() *Horsebase {

	hb = &Horsebase{
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	_, err := toml.DecodeFile("./file/horsebase.toml", &hb)
	if err != nil {
		PrintError(hb.Stderr, "%s", err)
		os.Exit(1)
	}

	return hb
}

func (hb *Horsebase) Run(args []string) int {

	if len(os.Args) < 2 {
		PrintError(hb.Stderr, "Invalid Argument")
		os.Exit(1)
	}

	param := os.Args[1]

	switch param {
	// DB構築
	// 初回起動
	case "-init_db":
		hb.DbInfo = hb.DbInfo.New()
		defer hb.DbInfo.db.Close()

		if err := hb.DbInfo.InitDB(); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}

		// 血統データ登録
	case "-reg_bloodtype":
		hb.DbInfo = hb.DbInfo.New()
		defer hb.DbInfo.db.Close()

		var bt BloodTypeToml
		bt = bt.New()

		if err := bt.Btd.RegistBloodType(hb.DbInfo); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}

		// レースデータのURLを取得しracelist.txtに一覧化する
	case "-make_list":
		if err := hb.MakeRaceURLList(); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}

	case "-get_racehtml":
		if err := hb.GetRaceHTML(); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}

		// 5
	case "-reg_racedata":
		if err := hb.RegistRaceData(); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}

		//6
	case "-reg_horsedata":
		if err := hb.RegistHorseData(); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}

	case "-drop_db":
		hb.DbInfo = hb.DbInfo.New()
		defer hb.DbInfo.db.Close()

		if err := hb.DbInfo.DropDB(); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}
	case "-match_bloodtype":
		hb.DbInfo = hb.DbInfo.New()
		defer hb.DbInfo.db.Close()

		var bt BloodTypeToml
		bt = bt.New()

		if err := bt.MatchBloodType(hb.DbInfo); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}

	default:
		PrintError(hb.Stderr, "Invalid Argument")
		return 1
	}

	return 0
}
