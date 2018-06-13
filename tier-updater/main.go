package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
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

	log.Println(dat.Top)
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
