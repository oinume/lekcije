package main

import (
	"encoding/csv"
	"flag"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/oinume/lekcije/server/bootstrap"
	"github.com/oinume/lekcije/server/model"
)

var (
	file = flag.String("file", "", "CSV file of country data")
)

func main() {
	flag.Parse()
	if *file == "" {
		log.Fatalf("Must specify -file")
	}

	f, err := os.Open(*file)
	if err != nil {
		log.Fatalf("File open error: path=%v, err=%v", *file, err)
	}

	bootstrap.CheckCLIEnvVars()
	db, err := model.OpenDB(bootstrap.CLIEnvVars.DBURL, 1, true)
	if err != nil {
		log.Fatal(err)
	}
	reader := csv.NewReader(f)
	for {
		columns, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v", err)
		}
		if columns[0] == "国・地域名" {
			continue
		}
		// "アイスランド","Iceland","352","ISL","IS","北ヨーロッパ","ISO 3166-2:IS"
		country := &model.MCountry{
			NameJA: columns[0],
			Name:   columns[1],
			ID:     parseCountryID(columns[2]),
		}
		if err := db.FirstOrCreate(country).Error; err != nil {
			log.Fatalf("%v", err)
		}
	}
}

func parseCountryID(s string) uint16 {
	if v, err := strconv.ParseUint(s, 10, 32); err == nil {
		return uint16(v)
	} else {
		panic(err)
	}
}
