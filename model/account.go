package model

import (
	"github.com/sleepynut/YottaDB-experiment/util"
)

type Account struct {
	// AccountID     string
	// UserID        string
	// Type          string
	// Balance       float64
	// CreateDt      time.Time
	// LastUpdatedDt time.Time
	// IsPrimary     bool
}

// func (a Account) toOracleFormat(fname string) ([]string, []map[string]util.ColValue) {
// 	var rows []map[string]util.ColValue

// 	f, err := os.Open(fname)
// 	if err != nil {
// 		log.Fatal("ERROR - openning file: ", err.Error())
// 	}
// 	defer f.Close()

// 	reader := bufio.NewReader(f)

// 	// skip header
// 	var header string
// 	var hs []string
// 	if header, err = reader.ReadString('\n'); err != nil {
// 		log.Fatal("ERROR - empty file")
// 	}

// 	for {
// 		m := make(map[string]util.ColValue)

// 		line, err := reader.ReadString('\n')
// 		if err != nil {
// 			if err == io.EOF && line == "" {
// 				break
// 			} else if err != io.EOF {
// 				log.Fatal("ERROR - reading file: ", err.Error())
// 			}
// 		}

// 		// remove trailling new line character
// 		line = strings.TrimSuffix(line, "\n")
// 		values := strings.Split(line, ",")

// 		// remove trailing new line character
// 		header = strings.TrimSuffix(header, "\n")
// 		hs = strings.Split(header, ",")
// 		for i, h := range hs {
// 			switch i {
// 			case 0, 1, 2:
// 				m[h] = util.ColValue{Value: values[i], Parser: util.SingleQuote}
// 			case 4, 5:
// 				m[h] = util.ColValue{Value: values[i], Parser: util.ToDateTime}
// 			default:
// 				m[h] = util.ColValue{Value: values[i], Parser: util.Vanilla}
// 			}
// 			// colTransformation(i, h, values, m)
// 		}
// 		rows = append(rows, m)
// 	}
// 	return hs, rows
// }

func (a Account) ColTransformation(i int, h string, values []string, m map[string]util.ColValue) {
	switch i {
	case 0, 1, 2:
		m[h] = util.ColValue{Value: values[i], Parser: util.SingleQuote}
	case 4, 5:
		m[h] = util.ColValue{Value: values[i], Parser: util.ToDateTime}
	default:
		m[h] = util.ColValue{Value: values[i], Parser: util.Vanilla}
	}
}
