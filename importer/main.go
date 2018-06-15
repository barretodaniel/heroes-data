package main

import (
	"database/sql"
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	_ "github.com/lib/pq"
)

// Hero is the shape that will go into the database
type Hero struct {
	Name       string `json:"name"`
	Portrait   string `json:"icon"`
	AttackType string `json:"type"`
	Role       string `json:"role"`
}

func main() {
	heroesDir := flag.String("dir", "", "-dir is the directory containing the hero files")
	flag.Parse()

	if *heroesDir == "" {
		log.Fatal("Please specify a directory with -dir")
	}

	files, err := ioutil.ReadDir(*heroesDir)
	check(err)

	connStr := "dbname=hots sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	check(err)

	defer db.Close()

	checkForHeroStmt, err := db.Prepare("SELECT COUNT(*) FROM heroes WHERE name = $1")
	check(err)
	defer checkForHeroStmt.Close()

	getRoleIDStmt, err := db.Prepare("SELECT id FROM roles WHERE name = $1")
	check(err)
	defer getRoleIDStmt.Close()

	insertHeroStmt, err := db.Prepare("INSERT INTO heroes(name, portrait, tier, attack_type, role_id) VALUES($1,$2,$3,$4,$5)")
	check(err)
	defer insertHeroStmt.Close()

	for _, file := range files {
		log.Println("Reading " + file.Name())
		dat, err := ioutil.ReadFile(*heroesDir + file.Name())

		if err != nil {
			log.Printf("Skipping file `%s` because there was an error reading it: %s\n", file, err.Error())
			continue
		}

		var info Hero
		err = json.Unmarshal(dat, &info)
		if err != nil {
			log.Printf("Skipping file `%s` because there was an error processing it: %s\n", file, err.Error())
			continue
		}

		var count int
		err = checkForHeroStmt.QueryRow(info.Name).Scan(&count)
		check(err)

		if count > 0 {
			log.Printf("Hero %q has already been imported... Skipping\n", info.Name)
			continue
		}

		var role int
		err = getRoleIDStmt.QueryRow(info.Role).Scan(&role)
		check(err)

		_, err = insertHeroStmt.Exec(info.Name, info.Portrait, 0, info.AttackType, role)
		check(err)
		log.Printf("%q was successfully imported\n", info.Name)
	}

	log.Println("Done!")
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
