package model

import t "github.com/sleepynut/YottaDB-experiment/transformer"

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

func (u User) ColTransformation(i int, h string, values []string, m map[string]t.ColValue) {
	switch i {
	// case 0, 1, 2:
	// 	m[h] = util.ColValue{Value: values[i], Parser: util.SingleQuote}
	case 7:
		m[h] = t.ColValue{Value: values[i], Parser: t.ToInt}
	default:
		m[h] = t.ColValue{Value: values[i], Parser: t.Vanilla}
	}
}
