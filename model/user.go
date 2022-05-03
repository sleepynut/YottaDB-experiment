package model

import "github.com/sleepynut/YottaDB-experiment/util"

type User struct {
	UserID      string
	Title       string
	FirstName   string
	LastName    string
	TitleEN     string
	FirstNameEN string
	LastNameEN  string
	Age         int
	Gender      string
}

func (u User) ColTransformation(i int, h string, values []string, m map[string]util.ColValue) {
	switch i {
	case 7:
		m[h] = util.ColValue{Value: values[i], Parser: util.Vanilla}
	default:
		m[h] = util.ColValue{Value: values[i], Parser: util.SingleQuote}
	}
}
