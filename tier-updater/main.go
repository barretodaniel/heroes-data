package main

import (
	"database/sql"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

type tierData struct {
	Top       []hero `json:"s"`
	TierOne   []hero `json:"t1"`
	TierTwo   []hero `json:"t2"`
	TierThree []hero `json:"t3"`
}

type hero struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func main() {
	res, err := http.Get("http://www.robogrub.com/tierlist_api")
	check(err)

	defer res.Body.Close()

	if res.StatusCode != 200 {
		log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
	}

	body, err := ioutil.ReadAll(res.Body)

	var dat tierData
	err = json.Unmarshal(body, &dat)
	check(err)

	connStr := "dbname=hots sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	check(err)

	defer db.Close()

	processTier(db, dat.Top, 0)
	processTier(db, dat.TierOne, 1)
	processTier(db, dat.TierTwo, 2)
	processTier(db, dat.TierThree, 3)

	log.Printf("Done!\n")
}

func processTier(db *sql.DB, heroes []hero, tier int) {
	findHeroStmt, err := db.Prepare("SELECT COUNT(*) FROM heroes WHERE LOWER(name) = $1")
	check(err)
	defer findHeroStmt.Close()

	updateHeroStmt, err := db.Prepare("update heroes set tier = $1 where LOWER(name) = $2")
	check(err)
	defer updateHeroStmt.Close()

	for _, h := range heroes {
		var count int
		var heroName string
		// Handle some cases where the ID doesn't match up with the DB name
		if h.ID == "Butcher" {
			heroName = "The Butcher"
		} else if h.ID == "Lost Vikings" {
			heroName = "The Lost Vikings"
		} else {
			heroName = h.ID
		}
		err = findHeroStmt.QueryRow(strings.ToLower(heroName)).Scan(&count)

		if count < 1 {
			log.Printf("%q not found\n", heroName)
			continue
		}

		_, err = updateHeroStmt.Exec(tier, strings.ToLower(heroName))
		if err != nil {
			log.Printf("There was an error updating %q... Skipping\n", heroName)
			continue
		}

		log.Printf("%q was updated successfully!\n", heroName)
	}
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
