package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/BurntSushi/toml"
)

const (
	layout  = "2006-01-02"
	tlayout = "04:05.0"
)

// Horsebase :
// horsebase.tomlの[horsebase]のパラメータ値を保持した構造体
type Horsebase struct {
	Config Config `toml:"config"`
	DbInfo HBDB   `toml:"db"`
	Stdout io.Writer
	Stderr io.Writer
}

// Config :
// horsebase.tomlで定義した設定値を保持した構造体
type Config struct {
	RaceHtmlPath  string `toml:"race_html_path"`
	HorseHtmlPath string `toml:"horse_html_path"`
	OldestDate    int    `toml:"oldest_date"`
}

// HBDB :
// DBアクセス用構造体
type HBDB struct {
	DbUser string `toml:"dbuser"`
	DbPass string `toml:"dbpass"`
	db     *sql.DB
	tx     *sql.Tx
}

// New :
// Horsebaseオブジェクト生成
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

	if len(args) != 2 {
		PrintError(hb.Stderr, "Invalid Argument")
		os.Exit(1)
	}

	var (
		initdb   bool
		regblood bool
		list     bool
		gethtml  bool
		regrace  bool
		reghorse bool
		dropdb   bool
		match    bool
		build    bool
	)

	//f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)

	flag.Usage = func() {
		fmt.Println(help)
	}
	flag.BoolVar(&initdb, "init_db", false, "")
	flag.BoolVar(&initdb, "i", false, "")

	flag.BoolVar(&regblood, "reg_bloodtype", false, "")

	flag.BoolVar(&list, "list", false, "")
	flag.BoolVar(&list, "l", false, "")

	flag.BoolVar(&gethtml, "get_racehtml", false, "")

	flag.BoolVar(&regrace, "reg_racedata", false, "")

	flag.BoolVar(&reghorse, "reg_horsedata", false, "")

	flag.BoolVar(&dropdb, "drop_db", false, "")
	flag.BoolVar(&dropdb, "d", false, "")

	flag.BoolVar(&match, "match_bloodtype", false, "")
	flag.BoolVar(&match, "m", false, "")

	flag.BoolVar(&build, "build", false, "")
	flag.BoolVar(&build, "b", false, "")

	flag.Parse()

	if initdb {
		// DB構築
		// 初回起動
		hb.DbInfo = hb.DbInfo.New()
		defer hb.DbInfo.db.Close()

		if err := hb.DbInfo.InitDB(); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}
	} else if regblood {
		// 血統データ登録
		hb.DbInfo = hb.DbInfo.New()
		defer hb.DbInfo.db.Close()

		var btt BloodTypeToml
		btt = btt.New()

		if err := btt.Btd.RegistBloodType(hb.DbInfo); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}
	} else if list {
		// レースデータのURLを取得しracelist.txtに保存する
		if err := hb.MakeRaceURLList(); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}
	} else if gethtml {

		if err := hb.GetRaceHTML(); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}
	} else if regrace {
		if err := hb.RegistRaceData(); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}
	} else if reghorse {

		if err := hb.RegistHorseData(); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}
	} else if dropdb {
		if err := hb.destroy(); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}
	} else if match {
		hb.DbInfo = hb.DbInfo.New()
		defer hb.DbInfo.db.Close()

		var btt BloodTypeToml
		btt = btt.New()

		if err := btt.MatchBloodType(hb.DbInfo); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}

	} else if build {
		if err := hb.build(); err != nil {
			PrintError(hb.Stderr, "%s", err)
		}
		return 1
	} else {
		PrintError(hb.Stderr, "Invalid Argument")
		return 1
	}

	return 0
}

func (hb *Horsebase) build() error {
	var err error
	hb.DbInfo = hb.DbInfo.New()
	defer hb.DbInfo.db.Close()

	if err = hb.DbInfo.InitDB(); err != nil {
		return err
	}

	if err = hb.MakeRaceURLList(); err != nil {
		return err
	}

	if err = hb.GetRaceHTML(); err != nil {
		return err
	}

	if err = hb.RegistHorseData(); err != nil {
		return err
	}

	if err = hb.RegistHorseData(); err != nil {
		return err
	}

	var btt BloodTypeToml
	btt = btt.New()

	if err = btt.MatchBloodType(hb.DbInfo); err != nil {
		return err
	}

	return err
}

func (hb *Horsebase) destroy() error {
	var err error

	fmt.Println("All data will be deleted, is it OK?[y/n] ")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		switch input {
		case "y", "Y":
			hb.DbInfo = hb.DbInfo.New()
			defer hb.DbInfo.db.Close()

			if err = hb.DbInfo.DropDB(); err != nil {
				return err
			}

		case "n", "N":
			return err
		}
	}
	return err
}

var help = `usage: horsebase [options]

Options:
  --build,-b            Stores all data

  --init_db,-i          Create horsebase DB
  --reg_bloodtype       Store the bloodtype data defined in bloodtype.toml in horsebase DB
  --make_list,-l        Save the URL of the race data in racelist.txt
  --get_racehtml        Gets the HTML form the URL listed in racelist.txt
  --reg_racedata        Scrape HTML and store race data in horsebase DB
  --reg_horsedata       Scrape HTML and store horse data in horsebase DB
  --drop_db,-d          Delete horsebase DB
  --match_bloodtype.-m  Map bloodtype data and stallion data defined in bloodtype.toml
`
