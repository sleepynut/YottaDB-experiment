package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-oci8"
	"github.com/sleepynut/YottaDB-experiment/model"
	"github.com/sleepynut/YottaDB-experiment/util"
)

const dsn = `hr/oracle@127.0.0.1:1521/orcl`

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("oci8", dsn)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// ping database
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	// else ping success
	fmt.Println("Connected!!")

	albums, err := albumsByArtist("John Coltrane")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Albums found: %v\n", albums)

	// Experiment
	fname := "./rsc/account.csv"
	// hs, colValues := model.Translate(fname, model.Account{})
	hs, colValues := util.ToOracleFormat(fname, (model.Account{}).ColTransformation)
	fmt.Println(util.InsertMany("account", hs, colValues, db))

	// fname = "./rsc/user.csv"
	// hs, colValues = util.ToOracleFormat(fname, (model.User{}).ColTransformation)
	// fmt.Println(util.InsertMany("user", hs, colValues, db))

}

type Album struct {
	ID     int64
	Title  string
	Artist string
	Price  float32
}

func albumsByArtist(name string) ([]Album, error) {
	var albums []Album

	rows, err := db.Query("SELECT * FROM album WHERE artist = :name", name)
	if err != nil {
		return nil, fmt.Errorf("albumsByArtist %q: %v", name, err)
	}

	for rows.Next() {
		var alb Album
		if err := rows.Scan(&alb.ID, &alb.Title, &alb.Artist, &alb.Price); err != nil {
			return nil, fmt.Errorf("albumsByArtist-rowscan %q: %v", name, err)
		}
		albums = append(albums, alb)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("albumsByArtist-rowErr %q: %v", name, err)
	}

	return albums, nil
}

func testParser(value string, parser func(string) interface{}) interface{} {
	return parser(value)
}
