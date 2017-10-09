package main

import (
	"os"
	"path"

	"github.com/BurntSushi/toml"
)

type BloodTypeToml struct {
	Btd BloodTypeDefine          `toml:"bloodtype"`
	Mbt map[string]MainBloodType `toml:"mainbloodtypes"`
	Sbt map[string]SubBloodType  `toml:"subbloodtypes"`
}

type BloodTypeDefine struct {
	Bloodtypes []string `toml:"bloodtypes"`
}

type MainBloodType struct {
	Stallions []string `toml:"stallions"`
}

type SubBloodType struct {
	Stallions []string `toml:"stallions"`
}

// New generates a BloodTypeToml object
func (btt BloodTypeToml) New() BloodTypeToml {
	_, err := toml.DecodeFile(path.Dir(os.Args[0])+"/file/bloodtype.toml", &btt)
	if err != nil {
		panic(err)
	}

	return btt
}

// RegistBloodType store blood-type information according to bloodtype.toml
func (btd BloodTypeDefine) RegistBloodType(hbdb HBDB) error {
	var err error
	for _, btname := range btd.Bloodtypes {

		err = hbdb.InsertBloodType(btname)
		if err != nil {
			return err
		}
	}
	return err
}

// matchBloodType map stallion data with mainbloodtype
// according to the mainbloodtype defined in bloodtype.toml
func (btt BloodTypeToml) matchBloodType(hbdb HBDB) error {
	var err error

	for _, btname := range btt.Btd.Bloodtypes {
		for _, name := range btt.Mbt[btname].Stallions {
			err = hbdb.UpdateMainBloodMatch(btname, name)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// matchSubBloodType map stallion data with sugbloodtype
// according to the subbloodtype defined in bloodtype.toml
func (btt BloodTypeToml) matchSubBloodType(hbdb HBDB) error {
	var err error

	for _, btname := range btt.Btd.Bloodtypes {
		for _, name := range btt.Sbt[btname].Stallions {
			err = hbdb.UpdateSubBloodMatch(btname, name)
			if err != nil {
				return err
			}
		}
	}

	return err
}

// MatchBloodType map stallion data with bloodtype
// according to the bloodtype defined in bloodtype.toml
func (btt BloodTypeToml) MatchBloodType(hbdb HBDB) error {
	var err error

	hbdb, err = hbdb.New()
	if err != nil {
		return err
	}
	defer hbdb.db.Close()

	err = hbdb.DeleteBloodType()
	if err != nil {
		return err
	}

	err = btt.Btd.RegistBloodType(hbdb)
	if err != nil {
		return err
	}

	err = btt.matchBloodType(hbdb)
	if err != nil {
		return err
	}

	err = btt.matchSubBloodType(hbdb)
	if err != nil {
		return err
	}

	return err
}
