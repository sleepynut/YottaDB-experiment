package util

import (
	"log"
	"strconv"
	"testing"
	"time"

	"github.com/sleepynut/YottaDB-experiment/model"
	t "github.com/sleepynut/YottaDB-experiment/transformer"
	"github.com/stretchr/testify/assert"
)

func TestRowToStruct(tt *testing.T) {

	// fname = "./rsc/user.csv"
	// hs, colValues = util.ToOracleFormat(fname, (model.User{}).ColTransformation)
	// fmt.Println(util.InsertMany("user", hs, colValues, db))

	// //reflect experiment
	// a := model.Account{}
	// fmt.Println("reflect: ", reflect.TypeOf(a).NumField())
	// var s = struct{ Foo int }{654}
	// var i = 200

	// rs := reflect.ValueOf(&s).Elem()
	// rf := rs.Field(0)
	// ri := reflect.ValueOf(i)

	// // rf = reflect.NewAt(rf.Type(), unsafe.Pointer(rf.UnsafeAddr())).Elem()
	// // ri.Set(rf)
	// rf.Set(ri)
	// fmt.Println("reflect - set: ", rf)
	// // fmt.Println(ri)

	// // Title experiment
	// fmt.Println(cases.Title(language.AmericanEnglish, cases.NoLower).String("sample"))

	// // any parser experiment
	// var f = func(s string) any {
	// 	v, err := strconv.ParseFloat(s, 64)
	// 	if err != nil {
	// 		log.Fatal("ERROR - parse float64", err.Error())
	// 	}
	// 	return v
	// }
	// fmt.Printf("%T\n", testParser("123.41", f))

	// // another reflect experiment
	m := make(map[string]t.ColValue)
	m["AccountID"] = t.ColValue{
		Value:  "1000000001",
		Parser: func(s string) any { return s },
	}

	m["UserID"] = t.ColValue{
		Value:  "1",
		Parser: func(s string) any { return s },
	}

	m["AccType"] = t.ColValue{
		Value:  "SA",
		Parser: func(s string) any { return s },
	}

	m["Balance"] = t.ColValue{
		Value: "123.45",
		Parser: func(s string) any {
			r, err := strconv.ParseFloat(s, 64)
			if err != nil {
				log.Fatal("ERROR - Float64 parser error", err.Error())
			}
			return r
		},
	}

	m["CreatedDt"] = t.ColValue{
		Value: "2021-05-04T15:53:00Z",
		Parser: func(s string) any {
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				log.Fatal("ERROR - time parser: ", err.Error())
			}
			return t
		},
	}

	m["LastUpdatedDt"] = t.ColValue{
		Value: "2021-05-04T15:53:00Z",
		Parser: func(s string) any {
			t, err := time.Parse(time.RFC3339, s)
			if err != nil {
				log.Fatal("ERROR - time parser: ", err.Error())
			}
			return t
		},
	}

	m["IsPrimary"] = t.ColValue{
		Value: "1",
		Parser: func(s string) any {
			b, err := strconv.ParseBool(s)
			if err != nil {
				log.Fatal("ERROR - bool parser: ", err.Error())
			}
			return b
		},
	}
	var anAccount any = model.Account{}

	dt, _ := time.Parse(time.RFC3339, "2021-05-04T15:53:00Z")
	expected := model.Account{
		AccountID:     "1000000001",
		UserID:        "1",
		AccType:       "SA",
		Balance:       123.45,
		CreatedDt:     dt,
		LastUpdatedDt: dt,
		IsPrimary:     true,
	}
	actual := rowToStruct(&anAccount, m).Interface()

	assert.EqualValues(tt, expected, actual)
}
