package model

import (
	"database/sql"
	"fmt"
	"log"
	"testing"

	_ "github.com/mattn/go-oci8"
	"github.com/stretchr/testify/assert"
)

const dsn = `hr/oracle@127.0.0.1:1521/orcl`

func TestSplitBalances(t *testing.T) {
	db := openConn()
	defer db.Close()

	dstAccounts := []DstAccount{
		{AccountID: "1000000002", Balance: 500},
		{AccountID: "1000000003", Balance: 1000},
		{AccountID: "1000000004", Balance: 500},
	}

	err := SplitBalances("1000000001", 2000, db, dstAccounts...)
	assert.Empty(t, err, fmt.Sprintf("Expect nil error, but got %s", err.Error()))
}

func openConn() *sql.DB {

	db, err := sql.Open("oci8", dsn)
	if err != nil {
		log.Fatal(err)
	}

	// ping database
	pingErr := db.Ping()
	if pingErr != nil {
		log.Fatal(pingErr)
	}

	// else ping success
	fmt.Println("Connected!!")
	return db
}
