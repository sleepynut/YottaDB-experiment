package util

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"
	"strings"
	"time"

	"github.com/sleepynut/YottaDB-experiment/model"
	t "github.com/sleepynut/YottaDB-experiment/transformer"
)

func InsertMany(tbName string, colNames []string,
	colValues []map[string]t.ColValue, db *sql.DB,
	target any, c chan string) {
	// prepare to measure the function execution time
	start := time.Now()

	stmt := `INSERT INTO %s (%s) VALUES (%s);`

	// begin transaction to insert multiple rows
	tx, err := db.Begin()
	if err != nil {
		log.Fatal("ERROR - begin transaction: ", err.Error())
	}

	// prepare column's name string
	names := strings.Join(colNames, ",")

	// prepare placeholder string
	placeHolders := make([]string, len(colNames))

	// formating placeholder in sql statement of the form :v.<fieldname>
	for i, n := range colNames {
		placeHolders[i] = fmt.Sprintf(":%s", n)
	}

	for i := 0; i < len(colValues); i++ {

		// transform map of column values into struct
		structVal := rowToStruct(&target, colValues[i])

		// prepare sql statement
		stmt = fmt.Sprintf(`INSERT INTO %s (%s) VALUES (%s)`,
			tbName, names, strings.Join(placeHolders, ","))
		stmSQL, err := db.Prepare(stmt)
		if err != nil {
			tx.Rollback()
			log.Fatal("ERROR - prepare INSERT: ", err.Error())
		}

		// execute sql statement
		var errExec error
		switch v := structVal.Interface().(type) {
		case model.Account:
			_, errExec = stmSQL.Exec(&v.AccountID, &v.UserID, &v.AccType, &v.Balance,
				&v.CreatedDt, &v.LastUpdatedDt, &v.IsPrimary)

		case model.User:
			_, errExec = stmSQL.Exec(&v.UserID, &v.Title, &v.FirstName, &v.LastName,
				&v.TitleEN, &v.FirstNameEN, &v.LastNameEN, &v.Age, &v.Gender)

		default:
			log.Fatalf("ERROR - Unknow type of %T\n", v)
		}

		if errExec != nil {
			tx.Rollback()
			log.Fatal("ERROR - Exec INSERT: ", errExec.Error())
		}

	}
	tx.Commit()

	// measure the total function execution time
	fmt.Printf("INSERTMANY(%s) - time taken: %s\n", tbName, time.Since(start).String())

	// push result into channel
	c <- stmt
}

func ToOracleFormat(fname string,
	transform func(int, string, []string, map[string]t.ColValue)) ([]string, []map[string]t.ColValue) {

	var rows []map[string]t.ColValue

	f, err := os.Open(fname)
	if err != nil {
		log.Fatal("ERROR - openning file: ", err.Error())
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	// skip header
	var header string
	if header, err = reader.ReadString('\n'); err != nil {
		log.Fatal("ERROR - empty file")
	}

	// remove trailing new line character
	header = strings.TrimSuffix(header, "\n")
	hs := strings.Split(header, ",")

	for {
		m := make(map[string]t.ColValue)

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF && line == "" {
				break
			} else if err != io.EOF {
				log.Fatal("ERROR - reading file: ", err.Error())
			}
		}

		// remove trailling new line character
		line = strings.TrimSuffix(line, "\n")
		values := strings.Split(line, ",")

		for i, h := range hs {
			transform(i, h, values, m)
		}
		rows = append(rows, m)
	}

	// return (header, row)
	// header is of format <name>, <name>, ...
	return hs, rows
}

func rowToStruct(aStruct *any, colValues map[string]t.ColValue) reflect.Value {
	structType := reflect.TypeOf(*aStruct)
	structValue := reflect.New(reflect.ValueOf(aStruct).Elem().Elem().Type()).Elem()

	for i := 0; i < structType.NumField(); i++ {
		name := structType.Field(i).Name

		if _, ok := colValues[name]; !ok {
			log.Fatalf("ERROR - reflect mismatch between target struct & given column value: %s NOT FOUND!", name)
		}

		value := colValues[name].Parser(colValues[name].Value)

		// set value to structValue
		structValue.Field(i).Set(reflect.ValueOf(value))
	}

	return structValue
}
