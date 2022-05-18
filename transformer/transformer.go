package transformer

import (
	"fmt"
	"log"
	"strconv"
	"time"
)

type ColValue struct {
	Value  string
	Parser func(string) any
}

func Vanilla(s string) any     { return s }
func SingleQuote(s string) any { return fmt.Sprintf("'%s'", s) }
func ToDateTime(s string) any {
	t, err := time.Parse(time.RFC3339, s)
	if err != nil {
		log.Fatal("ERROR - parse time: ", err.Error())
	}
	return t
}

func ToBool(s string) any {
	b, err := strconv.ParseBool(s)
	if err != nil {
		log.Fatal("ERROR - parse bool: ", err.Error())
	}
	return b
}

func ToFloat64(s string) any {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Fatal("ERROR - parse float64: ", err.Error())
	}
	return f
}

func ToInt(s string) any {
	n, err := strconv.Atoi(s)
	if err != nil {
		log.Fatal("ERROR - parse int: ", err.Error())
	}
	return n
}
