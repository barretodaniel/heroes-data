package main

import (
	"database/sql"
	"encoding/json"
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

const heroesDir = "/Users/danielbarreto/dev/heroes-talents/hero/"

func main() {
	files, err := ioutil.ReadDir(heroesDir)
	check(err)

	connStr := "dbname=hots sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	check(err)

	defer db.Close()

	for _, file := range files {
		log.Println("Reading " + file.Name())
		dat, err := ioutil.ReadFile(heroesDir + file.Name())

		if err != nil {
			log.Println("There was an error reading the file... Skipping")
			continue
		}

		var info Hero
		err = json.Unmarshal(dat, &info)
		if err != nil {
			log.Println("There was an error processing the file... Skipping")
			continue
		}

		stmt, err := db.Prepare("select name from heroes where name = $1")
		check(err)

		var name string
		err = stmt.QueryRow(info.Name).Scan(&name)

		if err == nil {
			log.Println("This has already been imported... Skipping")
			continue
		}

		var role int
		stmt, err = db.Prepare("select id from roles where name = $1")
		check(err)

		err = stmt.QueryRow(info.Role).Scan(&role)

		stmt, err = db.Prepare("insert into heroes(name, portrait, tier, attack_type, role_id) Values($1,$2,$3,$4,$5)")
		check(err)

		_, err = stmt.Exec(info.Name, info.Portrait, 0, info.AttackType, role)
		check(err)
	}
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
