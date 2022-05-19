package main

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/mattn/go-oci8"
	"github.com/sleepynut/YottaDB-experiment/model"
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

	// albums, err := albumsByArtist("John Coltrane")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Printf("Albums found: %v\n", albums)

	// // create receiving channel
	// c := make(chan string)

	// // many inserts into account
	// fname := "./rsc/account.csv"
	// var astruct any = model.Account{}
	// hs, colValues := util.ToOracleFormat(fname, (model.Account{}).ColTransformation)
	// go util.InsertMany("account", hs, colValues, db, astruct, c)

	// // many inserts into user
	// fname = "./rsc/user.csv"
	// astruct = model.User{}
	// hs, colValues = util.ToOracleFormat(fname, (model.User{}).ColTransformation)
	// go util.InsertMany("userInfo", hs, colValues, db, astruct, c)

	// // wait for the return of sub routine
	// fmt.Println(<-c)
	// fmt.Println(<-c)

	// {
	// 	fmt.Println("INSIDE TEMP")
	// 	acc := "1000000001"
	// 	bal := "2000000"

	// 	ctx, cancel := context.WithTimeout(context.Background(), 55*time.Second)

	// 	defer cancel()

	// 	time.Sleep(10 * time.Second)

	// 	tx, err := db.BeginTx(ctx, nil)
	// 	defer func() {
	// 		if err != nil {
	// 			fmt.Println("ANY ERROR: ", err.Error())
	// 		}
	// 		fmt.Println("getting here?")
	// 		tx.Rollback()
	// 	}()

	// 	stmt := "UPDATE account set balance=:bal where accountID=:acc"
	// 	_, err = tx.ExecContext(ctx, stmt, bal, acc)
	// }

	// ACID: atomic
	// err = model.UpdateBalance("1000000001", -50000, db)
	// err = model.MoveBalance("1000000001", "1000000007", 50000, db)
	// if err != nil {
	// 	fmt.Println("ERROR - main: ", err.Error())
	// }

	// ACID: consistency
	// model.MakePrimaryAccount("1000000002", "1", db)

	// ACID: Isolation
	dstAccounts := []model.DstAccount{
		{AccountID: "1000000002", Balance: 500},
		{AccountID: "1000000003", Balance: 1000},
		{AccountID: "1000000004", Balance: 500},
	}

	err = model.SplitBalances("1000000001", 2000, db, dstAccounts...)
	if err != nil {
		fmt.Println(err.Error())
	}

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
