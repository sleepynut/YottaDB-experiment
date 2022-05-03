package util

import (
	"bufio"
	"database/sql"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

type ColValue struct {
	Value  string
	Parser func(string) string
}

func InsertMany(tbName string, colNames []string, colValues []map[string]ColValue, db *sql.DB) string {
	stmt := "INSERT ALL\n"

	for i := 0; i < len(colValues); i++ {
		cols := make([]string, len(colValues[i]))

		for j, name := range colNames {
			cols[j] = colValues[i][name].Parser(colValues[i][name].Value)
		}
		temp := fmt.Sprintf("\tINTO %s (%s) VALUES (%s)\n",
			tbName, strings.Join(colNames, ","), strings.Join(cols, ","))
		stmt += temp
	}

	stmt += "SELECT * from dual;\n"

	// temp
	stmt = `
INSERT INTO account (accountID,userID,accType,balance,createdDt,lastUpdatedDt,isPrimary)
SELECT '1000000001','1','SV',100000.00,TO_DATE('2020-04-30 14:04:05','yyyy-mm-dd hh24:mi:ss'),TO_DATE('2020-04-30 14:04:05','yyyy-mm-dd hh24:mi:ss'),1 from DUAL UNION ALL
SELECT '1000000002','1','SV',1100000.00,TO_DATE('2020-04-30 14:04:05','yyyy-mm-dd hh24:mi:ss'),TO_DATE('2020-04-30 14:04:05','yyyy-mm-dd hh24:mi:ss'),0 from DUAL;`
	// temp-end

	stmtSQL, err := db.Prepare(stmt)
	if err != nil {
		log.Fatal("ERROR - exec PREPARE MANY: ", err.Error())
	}

	result, err := stmtSQL.Exec()
	if err != nil {
		log.Fatal("ERROR - exec INSERT STMT: ", err.Error())
	}

	fmt.Println(result)
	return stmt
}

func Vanilla(s string) string     { return s }
func SingleQuote(s string) string { return fmt.Sprintf("'%s'", s) }
func ToDateTime(s string) string {
	return fmt.Sprintf("TO_DATE('%s','yyyy-mm-dd hh24:mi:ss')", s)
}

func ToOracleFormat(fname string,
	transform func(int, string, []string, map[string]ColValue)) ([]string, []map[string]ColValue) {

	var rows []map[string]ColValue

	f, err := os.Open(fname)
	if err != nil {
		log.Fatal("ERROR - openning file: ", err.Error())
	}
	defer f.Close()

	reader := bufio.NewReader(f)

	// skip header
	var header string
	var hs []string
	if header, err = reader.ReadString('\n'); err != nil {
		log.Fatal("ERROR - empty file")
	}

	for {
		m := make(map[string]ColValue)

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

		// remove trailing new line character
		header = strings.TrimSuffix(header, "\n")
		hs = strings.Split(header, ",")
		for i, h := range hs {
			transform(i, h, values, m)
		}
		rows = append(rows, m)
	}
	return hs, rows
}
