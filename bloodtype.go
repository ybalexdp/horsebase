package main

import "github.com/BurntSushi/toml"

// BloodTypeToml :系統と種牡馬を定義したtomlファイルのデータ格納用
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

func (btt BloodTypeToml) New() BloodTypeToml {
	_, err := toml.DecodeFile("./file/bloodtype.toml", &btt)
	if err != nil {
		panic(err)
	}

	return btt
}

func (bt BloodTypeDefine) RegistBloodType(hbdb HBDB) error {
	var err error
	for _, btname := range bt.Bloodtypes {

		err = hbdb.InsertBloodType(btname)
		if err != nil {
			return err
		}
	}
	return err
}

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

func (btt BloodTypeToml) MatchBloodType(hbdb HBDB) error {
	var err error

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
