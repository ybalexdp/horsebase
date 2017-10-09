package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/BurntSushi/toml"
)

const (
	layout = "2006-01-02" // for Date Format
	//tlayout = "04:05.0"
)

type Horsebase struct {
	Config Config `toml:"config"`
	DbInfo HBDB   `toml:"db"`
	Stdout io.Writer
	Stderr io.Writer
}

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

// New generates a horsebase object
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

// Run executes process according to specified option
func (hb *Horsebase) Run(args []string) int {
	var err error

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
		update   bool
	)

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

	flag.BoolVar(&update, "update", false, "")
	flag.BoolVar(&update, "u", false, "")

	flag.Parse()

	// Create horsebase DB
	if initdb {
		hb.DbInfo, err = hb.DbInfo.New()
		if err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}
		defer hb.DbInfo.db.Close()

		if err := hb.DbInfo.InitDB(); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}

		// Store the bloodtype data defined in bloodtype.toml in horsebase DB
	} else if regblood {
		hb.DbInfo, err = hb.DbInfo.New()
		if err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}
		defer hb.DbInfo.db.Close()

		var btt BloodTypeToml
		btt = btt.New()

		if err := btt.Btd.RegistBloodType(hb.DbInfo); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}

		// Save the URL of the race data in racelist.txt
	} else if list {
		if err := hb.MakeRaceURLList(); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}

		// Get the HTML form the URL listed in racelist.txt
	} else if gethtml {

		if err := hb.GetRaceHTML(); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}

		// Scrape HTML and store race data in horsebase DB
	} else if regrace {
		if err := hb.RegistRaceData(); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}

		// Scrape HTML and store horse data in horsebase DB
	} else if reghorse {

		if err := hb.RegistHorseData(); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}

		// Delete horsebase DB
	} else if dropdb {
		if err := hb.destroy(); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}

		// Map bloodtype data and stallion data defined in bloodtype.toml
	} else if match {
		hb.DbInfo, err = hb.DbInfo.New()
		if err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}
		defer hb.DbInfo.db.Close()

		var btt BloodTypeToml
		btt = btt.New()

		if err := btt.MatchBloodType(hb.DbInfo); err != nil {
			PrintError(hb.Stderr, "%s", err)
			return 1
		}
	} else if update {

		if err := hb.update(); err != nil {
			PrintError(hb.Stderr, "%s", err)
		}
		return 1

		// Store all data
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

// build store all data
func (hb *Horsebase) build() error {
	var err error

	if err = hb.DbInfo.InitDB(); err != nil {
		return err
	}

	if err = hb.MakeRaceURLList(); err != nil {
		return err
	}

	if err = hb.GetRaceHTML(); err != nil {
		return err
	}

	if err = hb.RegistRaceData(); err != nil {
		PrintError(hb.Stderr, "%s", err)
		return fmt.Errorf("Please retry \n $ horsebase --reg_racedata")
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

func (hb *Horsebase) update() error {
	var err error
	hb.DbInfo, err = hb.DbInfo.New()
	if err != nil {
		return err
	}

	date, _ := hb.DbInfo.GetLatestDate()
	defer hb.DbInfo.db.Close()

	date = strings.Replace(date, "-", "", -1)
	oldestdate, _ := strconv.Atoi(date)
	hb.Config.OldestDate = oldestdate

	if err = hb.MakeRaceURLList(); err != nil {
		return err
	}

	if err = hb.GetRaceHTML(); err != nil {
		return err
	}

	if err = hb.RegistRaceData(); err != nil {
		PrintError(hb.Stderr, "%s", err)
		return fmt.Errorf("Please retry \n $ horsebase --reg_racedata")
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

// destroy delete horsebase DB
func (hb *Horsebase) destroy() error {
	var err error

	fmt.Print("All data will be deleted, is it OK?[y/n] ")

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()

		switch input {
		case "y", "Y":
			hb.DbInfo, err = hb.DbInfo.New()
			if err != nil {
				PrintError(hb.Stderr, "%s", err)
				return err
			}
			defer hb.DbInfo.db.Close()

			if err = hb.DbInfo.DropDB(); err != nil {
				return err
			}
			break

		case "n", "N":
			return err

		default:
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
