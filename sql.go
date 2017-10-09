package main

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/BurntSushi/toml"
	_ "github.com/go-sql-driver/mysql"
	//"fmt"
)

func (hbdb HBDB) New() (HBDB, error) {
	_, err := toml.DecodeFile(path.Dir(os.Args[0])+"/file/horsebase.toml", &hbdb)
	if err != nil {
		return hbdb, err
	}
	hbdb.db, err = sql.Open("mysql", hbdb.DbUser+":"+hbdb.DbPass+"@/horsebase")
	if err != nil {
		return hbdb, fmt.Errorf("Check the status of mysql service and dbuser/dbpass in horsebase.toml")
	}
	return hbdb, err
}

func (hbdb HBDB) InitDB() error {
	db, err := sql.Open("mysql", hbdb.DbUser+":"+hbdb.DbPass+"@/")
	if err != nil {
		return fmt.Errorf("Check the status of mysql service and dbuser/dbpass in horsebase.toml")
	}
	err = createDB(db)
	if err != nil {
		return fmt.Errorf("Check the status of mysql service and dbuser/dbpass in horsebase.toml")
	}

	db, err = sql.Open("mysql", hbdb.DbUser+":"+hbdb.DbPass+"@/horsebase")
	if err != nil {
		return fmt.Errorf("Check the status of mysql service and dbuser/dbpass in horsebase.toml")
	}

	err = createTable(db)
	if err != nil {
		hbdb.DropDB()
		return fmt.Errorf("Check the status of mysql service and dbuser/dbpass in horsebase.toml")
	}

	err = createIDX(db)
	if err != nil {
		hbdb.DropDB()
		return err
	}

	defer db.Close()

	return err

}

func (hbdb HBDB) DropDB() error {
	query := "DROP DATABASE	horsebase"
	_, err := hbdb.db.Exec(query)
	return err
}

func (hbdb HBDB) DeleteBloodType() error {
	query := "DELETE FROM bloodtype"
	_, err := hbdb.db.Exec(query)
	return err
}

func (hbdb HBDB) InsertBloodType(btname string) error {

	query := "INSERT INTO bloodtype(name) SELECT * FROM (SELECT ?) AS TMP WHERE NOT EXISTS(SELECT * FROM bloodtype WHERE name=?)"
	_, err := hbdb.db.Exec(query, btname, btname)
	return err
}

func (hbdb HBDB) InsertHorse(horse Horse) error {
	query := "INSERT INTO horse(id,name) VALUES(?,?)"
	_, err := hbdb.tx.Exec(query, horse.HorseID, horse.Name)
	return err
}

func (hbdb HBDB) UpdateHorse(horse Horse) error {
	query := "UPDATE horse SET father_id=?,father_m_id=?,father_fm_id=?,father_mm_id=? WHERE id=?"
	_, err := hbdb.db.Exec(query, horse.Father, horse.FatherOfM, horse.FatherOfFM, horse.FatherOfMM, horse.HorseID)
	return err
}

func (hbdb HBDB) InsertJockey(jockey Jockey) error {
	query := "INSERT INTO jockey(name) SELECT * FROM (SELECT ?) AS TMP WHERE NOT EXISTS(SELECT * FROM jockey WHERE name=?)"
	_, err := hbdb.db.Exec(query, jockey.Name, jockey.Name)
	return err
}

func (hbdb HBDB) InsertRaceData(rd RaceData) error {
	query := "INSERT INTO racedata VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,null,?,?)"
	_, err := hbdb.tx.Exec(query, rd.RaceID, rd.Name, rd.Course, rd.Corner, rd.Distance, rd.Date.Format(layout),
		rd.Grade, rd.Turf, rd.RaceNumber, rd.Day, rd.Surface, rd.Weather, rd.TrackCond, rd.AgeGr, rd.SexGr)
	return err
}

func (hbdb HBDB) UpdateHorseNum(rd RaceData) error {
	query := "UPDATE racedata SET horsenum=? WHERE id=?"
	_, err := hbdb.tx.Exec(query, rd.Horsenum, rd.RaceID)
	return err
}

func (hbdb HBDB) InsertRaceresult(rrd RaceResultData) error {

	stime := rrd.Time.String()
	stime = strings.Split(stime, "+")[0]

	query := "INSERT INTO raceresult VALUES(?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)"
	_, err := hbdb.tx.Exec(query, rrd.RaceID, rrd.HorseID, rrd.Rank, rrd.JockeyID, rrd.Popularity, rrd.Odds, rrd.Age,
		rrd.Weight, rrd.Bweight, rrd.Hnumber, rrd.Wnumber, rrd.LastThreeFur, rrd.Sex, stime, rrd.DifTime,
		rrd.POrder[0], rrd.POrder[1], rrd.POrder[2], rrd.POrder[3], rrd.Belonging)

	return err
}

func (hbdb HBDB) InsertStallion(name string) error {
	query := "INSERT INTO stallion(name) SELECT * FROM (SELECT ?) AS TMP WHERE NOT EXISTS(SELECT * FROM stallion WHERE name=?)"
	_, err := hbdb.db.Exec(query, name, name)
	return err
}

func (hbdb HBDB) InsertWinData(rd RaceData, i int) error {
	query := "INSERT INTO win VALUES(?,?,?,?)"
	_, err := hbdb.tx.Exec(query, rd.RaceID, rd.Win.HorseNum[i], rd.Win.Dividend[i], rd.Win.Popularity[i])
	return err
}

func (hbdb HBDB) InsertPlaceData(rd RaceData, i int) error {
	query := "INSERT INTO place VALUES(?,?,?,?)"
	_, err := hbdb.tx.Exec(query, rd.RaceID, rd.Place.HorseNum[i], rd.Place.Dividend[i], rd.Place.Popularity[i])
	return err
}

func (hbdb HBDB) InsertQuinellaData(rd RaceData, i int) error {
	query := "INSERT INTO quinella VALUES(?,?,?,?,?)"
	_, err := hbdb.tx.Exec(query, rd.RaceID, rd.Quinella.HorseNum[i][0], rd.Quinella.HorseNum[i][1], rd.Quinella.Popularity[i])
	return err
}

func (hbdb HBDB) InsertExactaData(rd RaceData, i int) error {
	query := "INSERT INTO exacta VALUES(?,?,?,?,?)"
	_, err := hbdb.tx.Exec(query, rd.RaceID, rd.Exacta.HorseNum[i][0], rd.Exacta.HorseNum[i][1], rd.Exacta.Popularity[i])
	return err
}

func (hbdb HBDB) InsertQPData(rd RaceData, i int) error {
	query := "INSERT INTO qp VALUES(?,?,?,?,?)"
	_, err := hbdb.tx.Exec(query, rd.RaceID, rd.QP.HorseNum[i][0], rd.QP.HorseNum[i][1], rd.QP.Popularity[i])
	return err
}

func (hbdb HBDB) InsertTrioData(rd RaceData, i int) error {
	query := "INSERT INTO trio VALUES(?,?,?,?,?)"
	_, err := hbdb.tx.Exec(query, rd.RaceID, rd.Trio.HorseNum[i][0], rd.Trio.HorseNum[i][1],
		rd.Trio.HorseNum[i][2], rd.Trio.Popularity[i])
	return err
}

func (hbdb HBDB) InsertTrifectaData(rd RaceData, i int) error {
	query := "INSERT INTO trifecta VALUES(?,?,?,?,?)"
	_, err := hbdb.tx.Exec(query, rd.RaceID, rd.Trifecta.HorseNum[i][0], rd.Trifecta.HorseNum[i][1],
		rd.Trifecta.HorseNum[i][2], rd.Trifecta.Popularity[i])
	return err
}

func createDB(db *sql.DB) error {
	_, err := db.Exec("CREATE DATABASE IF NOT EXISTS horsebase")
	db.Close()
	return err
}

func createTable(db *sql.DB) error {
	var err error
	query := `CREATE TABLE IF NOT EXISTS horsebase.racedata
(id BIGINT NOT NULL PRIMARY KEY,
name VARCHAR(20),
course INT,
corner INT,
distance INT,
date DATE,
grade INT,
turf VARCHAR(20),
racenumber INT,
day INT,
surface INT,
weather INT,
track_cond INT,
horsenum INT,
age_gr INT,
sex_gr INT)`

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS horsebase.bloodtype
( id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
name VARCHAR(20)
)`

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS horsebase.stallion
(id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
name VARCHAR(32),
bloodtype_id INT,
subbloodtype_id INT,
FOREIGN KEY(bloodtype_id) REFERENCES bloodtype(id) ON DELETE SET NULL,
FOREIGN KEY(subbloodtype_id) REFERENCES bloodtype(id) ON DELETE SET NULL
)`

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS horsebase.horse
(id INT NOT NULL PRIMARY KEY,
name VARCHAR(16),
father_id INT,
father_m_id INT,
father_fm_id INT,
father_mm_id INT,
FOREIGN KEY(father_id) REFERENCES stallion(id) ON DELETE SET NULL,
FOREIGN KEY(father_m_id) REFERENCES stallion(id) ON DELETE SET NULL,
FOREIGN KEY(father_fm_id) REFERENCES stallion(id) ON DELETE SET NULL,
FOREIGN KEY(father_mm_id) REFERENCES stallion(id) ON DELETE SET NULL
)`

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS horsebase.jockey
(id INT AUTO_INCREMENT NOT NULL PRIMARY KEY,
name VARCHAR(20)
)`

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS horsebase.raceresult
(id BIGINT NOT NULL,
horse_id INT NOT NULL,
rank INT,
jockey_id INT,
popularity INT,
odds DOUBLE(5,1),
age INT,
weight INT,
bweight DOUBLE(3,1),
horse_num INT,
waku_num INT,
last3f DOUBLE(3,1),
sex INT,
time TIME(1),
diftime DOUBLE(3,1),
passing_order1 INT,
passing_order2 INT,
passing_order3 INT,
passing_order4 INT,
belonging INT,
PRIMARY KEY(id, horse_id),
FOREIGN KEY(id) REFERENCES racedata(id) ON DELETE RESTRICT,
FOREIGN KEY(horse_id) REFERENCES horse(id) ON DELETE RESTRICT,
FOREIGN KEY(jockey_id) REFERENCES jockey(id) ON DELETE SET NULL
)`

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS horsebase.win
(race_id BIGINT NOT NULL,
horse_number INT NOT NULL,
dividend INT,
popularity INT,
PRIMARY KEY(race_id, horse_number),
FOREIGN KEY(race_id) REFERENCES racedata(id) ON DELETE RESTRICT
)`

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS horsebase.place
(race_id BIGINT NOT NULL,
horse_number INT NOT NULL,
dividend INT,
popularity INT,
PRIMARY KEY(race_id, horse_number),
FOREIGN KEY(race_id) REFERENCES racedata(id) ON DELETE RESTRICT
)`

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS horsebase.quinella
(race_id BIGINT NOT NULL,
horse_number1 INT NOT NULL,
horse_number2 INT NOT NULL,
dividend INT,
popularity INT NOT NULL,
PRIMARY KEY(race_id, horse_number1, horse_number2),
FOREIGN KEY(race_id) REFERENCES racedata(id) ON DELETE RESTRICT
)`

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS horsebase.exacta
(race_id BIGINT NOT NULL,
horse_number1 INT NOT NULL,
horse_number2 INT NOT NULL,
dividend INT,
popularity INT NOT NULL,
PRIMARY KEY(race_id, horse_number1, horse_number2),
FOREIGN KEY(race_id) REFERENCES racedata(id) ON DELETE RESTRICT
)`

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS horsebase.qp
(race_id BIGINT NOT NULL,
horse_number1 INT NOT NULL,
horse_number2 INT NOT NULL,
dividend INT,
popularity INT NOT NULL,
PRIMARY KEY(race_id, horse_number1, horse_number2),
FOREIGN KEY(race_id) REFERENCES racedata(id) ON DELETE RESTRICT
)`

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS horsebase.trio
(race_id BIGINT NOT NULL,
horse_number1 INT NOT NULL,
horse_number2 INT NOT NULL,
horse_number3 INT NOT NULL,
dividend INT,
popularity INT NOT NULL,
PRIMARY KEY(race_id, horse_number1, horse_number2, horse_number3),
FOREIGN KEY(race_id) REFERENCES racedata(id) ON DELETE RESTRICT
)`

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = `CREATE TABLE IF NOT EXISTS horsebase.trifecta
(race_id BIGINT NOT NULL,
horse_number1 INT NOT NULL,
horse_number2 INT NOT NULL,
horse_number3 INT NOT NULL,
dividend INT,
popularity INT NOT NULL,
PRIMARY KEY(race_id, horse_number1, horse_number2, horse_number3),
FOREIGN KEY(race_id) REFERENCES racedata(id) ON DELETE RESTRICT
)`

	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return err

}

func createIDX(db *sql.DB) error {
	var err error

	query := "ALTER TABLE racedata ADD INDEX idx_racedata_id(id)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE racedata ADD INDEX idx_racedata_name(name)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE racedata ADD INDEX idx_racedata_course(course)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE racedata ADD INDEX idx_racedata_corner(corner)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE racedata ADD INDEX idx_racedata_distance(distance)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE racedata ADD INDEX idx_racedata_grade(grade)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE racedata ADD INDEX idx_racedata_turf(turf)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE racedata ADD INDEX idx_racedata_surface(surface)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE racedata ADD INDEX idx_racedata_weather(weather)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE racedata ADD INDEX idx_racedata_track_cond(track_cond)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE racedata ADD INDEX idx_racedata_horsenum(horsenum)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE racedata ADD INDEX idx_racedata_age_gr(age_gr)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE racedata ADD INDEX idx_racedata_sex_gr(sex_gr)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE bloodtype ADD INDEX idx_bloodtype_id(id)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE bloodtype ADD INDEX idx_bloodtype_name(name)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE stallion ADD INDEX idx_stallion_id(id)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE stallion ADD INDEX idx_stallion_name(name)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE stallion ADD INDEX idx_stallion_bloodtype(bloodtype_id)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE stallion ADD INDEX idx_stallion_subbloodtype(subbloodtype_id)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE horse ADD INDEX idx_horse_name(name)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE raceresult ADD INDEX idx_raceresult_rank(rank)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE raceresult ADD INDEX idx_raceresult_popularity(popularity)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE raceresult ADD INDEX idx_raceresult_age(age)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE raceresult ADD INDEX idx_raceresult_weight(weight)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE raceresult ADD INDEX idx_raceresult_bweight(bweight)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE raceresult ADD INDEX idx_raceresult_horse_num(horse_num)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE raceresult ADD INDEX idx_raceresult_waku_num(waku_num)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE raceresult ADD INDEX idx_raceresult_last3f(last3f)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE raceresult ADD INDEX idx_raceresult_sex(sex)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE raceresult ADD INDEX idx_raceresult_time(time)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE raceresult ADD INDEX idx_raceresult_diftime(diftime)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE raceresult ADD INDEX idx_raceresult_po1(passing_order1)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE raceresult ADD INDEX idx_raceresult_po2(passing_order2)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE raceresult ADD INDEX idx_raceresult_po3(passing_order3)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE raceresult ADD INDEX idx_raceresult_po4(passing_order4)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	query = "ALTER TABLE raceresult ADD INDEX idx_raceresult_belonging(belonging)"
	_, err = db.Exec(query)
	if err != nil {
		return err
	}

	return err

}

func (hbdb HBDB) UpdateMainBloodMatch(bt string, stallion string) error {
	query := "UPDATE stallion, bloodtype SET stallion.bloodtype_id=bloodtype.id WHERE stallion.name=? AND bloodtype.name=?"
	_, err := hbdb.db.Exec(query, stallion, bt)
	return err
}

func (hbdb HBDB) UpdateSubBloodMatch(bt string, stallion string) error {
	query := "UPDATE stallion, bloodtype SET stallion.subbloodtype_id=bloodtype.id WHERE stallion.name=? AND bloodtype.name=?"
	_, err := hbdb.db.Exec(query, stallion, bt)
	return err
}

/*
func execSQL(db *sql.DB, query string) error {
	_, err := db.Exec(query)
	return err
}
*/

func (hbdb HBDB) GetId(table string, name string) (int, error) {
	var id int
	query := "SELECT id FROM " + table + " WHERE name=\"" + name + "\""
	err := hbdb.db.QueryRow(query).Scan(&id)
	if err != nil && err != sql.ErrNoRows {
		return id, err
	}
	return id, nil
}

func (hbdb HBDB) GetLatestDate() (string, error) {
	var date string
	query := "SELECT MAX(date) FROM racedata"
	err := hbdb.db.QueryRow(query).Scan(&date)
	return date, err
}

func (hbdb HBDB) RaceExistenceCheck(raceID string) (bool, error) {
	query := "SELECT id FROM racedata WHERE id = " + raceID
	return hbdb.rowExists(query)
}

func (hbdb HBDB) HorseExistenceCheck(horseID string) (bool, error) {
	query := "SELECT id FROM horse WHERE id = " + horseID
	return hbdb.rowExists(query)
}

func (hbdb HBDB) rowExists(query string) (bool, error) {
	var exists bool
	query = fmt.Sprintf("SELECT EXISTS (%s)", query)
	err := hbdb.db.QueryRow(query).Scan(&exists)
	if err != nil && err != sql.ErrNoRows {
		return false, err
	}
	return exists, nil
}
