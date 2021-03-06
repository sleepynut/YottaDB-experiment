package model

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	c "github.com/sleepynut/YottaDB-experiment/custom"
	t "github.com/sleepynut/YottaDB-experiment/transformer"
)

type Account struct {
	AccountID     string
	UserID        string
	AccType       string
	Balance       float64
	CreatedDt     time.Time
	LastUpdatedDt time.Time
	IsPrimary     bool
}

type DstAccount struct {
	AccountID string
	Balance   float64
}

func (a Account) ColTransformation(i int, h string, values []string, m map[string]t.ColValue) {
	switch i {
	// case 0, 1, 2:
	// 	m[h] = util.ColValue{Value: values[i], Parser: util.SingleQuote}
	case 3:
		m[h] = t.ColValue{Value: values[i], Parser: t.ToFloat64}
	case 4, 5:
		m[h] = t.ColValue{Value: values[i], Parser: t.ToDateTime}
	case 6:
		m[h] = t.ColValue{Value: values[i], Parser: t.ToBool}
	default:
		m[h] = t.ColValue{Value: values[i], Parser: t.Vanilla}
	}
}

func UpdateBalance(id string, balance float64, tx *sql.Tx) error {
	stmt := `SELECT * from account where accountID=:id`
	row := tx.QueryRow(stmt, id)

	var acc Account
	if err := row.Scan(&acc.AccountID, &acc.UserID, &acc.AccType,
		&acc.Balance, &acc.CreatedDt, &acc.LastUpdatedDt, &acc.IsPrimary); err != nil {

		if err == sql.ErrNoRows {
			return c.RecordNotFound{TbName: "account", Id: id}
		}

		return fmt.Errorf("ERROR - updatedBalance(query): %s", err.Error())
	}

	if acc.Balance+balance < 0 {
		return fmt.Errorf(
			"operation could not be done: insufficient balance. Perform: %.2f Actual: %.2f",
			balance, acc.Balance)

	}

	balance += acc.Balance
	now := time.Now()
	stmt = `UPDATE account set balance=:balance,lastUpdatedDt=:now where accountId=:id`
	result, err := tx.Exec(stmt, balance, now, id)
	if err != nil {
		return fmt.Errorf("ERROR - updatedBalance(update): %s", err.Error())
	}

	nrow, _ := result.RowsAffected()
	fmt.Println("Effected row: ", nrow)

	return nil
}

func MoveBalance(src, dest string, balance float64, db *sql.DB) error {
	if balance < 0 {
		return errors.New("prohibit negative balance")
	}

	tx, err := db.Begin()
	defer func() {
		if err != nil {
			// rollback all previous moving balance operation
			tx.Rollback()

			// tx logging failed movement
			TxLogging("CANCEL", src, dest, balance, db)
		} else {

			// tx logging success movement
			TxLogging("SUCCESS", src, dest, balance, db)

		}

		// commit all dirty transaction
		tx.Commit()

	}()

	if err = UpdateBalance(src, -balance, tx); err != nil {
		if _, ok := err.(c.RecordNotFound); !ok {
			return err
		}
	}

	if err = UpdateBalance(dest, balance, tx); err != nil {
		return err
	}

	return nil
}

func MoveBalanceTx(src, dest string, balance float64, db *sql.DB, tx *sql.Tx) error {
	if balance < 0 {
		return errors.New("prohibit negative balance")
	}

	var err error
	// defer func() {
	// 	if err != nil {
	// 		// tx logging failed movement
	// 		TxLogging("CANCEL", src, dest, balance, db)
	// 	} else {

	// 		// tx logging success movement
	// 		TxLogging("SUCCESS", src, dest, balance, db)

	// 	}
	// }()

	if err = UpdateBalance(src, -balance, tx); err != nil {
		if _, ok := err.(c.RecordNotFound); !ok {
			return err
		}
	}

	if err = UpdateBalance(dest, balance, tx); err != nil {
		return err
	}

	return nil
}

func MoveBalanceThenCancel(src, dst string, balance float64, db *sql.DB) error {
	if balance < 0 {
		return errors.New("prohibit negative balance")
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("ERROR - MoveBalanceThenCancel: Begin tx : %s", err.Error())
	}

	defer func() {
		tx.Rollback()
	}()

	err = MoveBalanceTx(src, dst, balance, db, tx)

	if err != nil {
		return fmt.Errorf("ERROR - MoveBalanceThenCancel: moving: %s", err.Error())
	}

	return nil
}

func SplitBalances(src string, srcBalance float64, db *sql.DB, dstAccounts ...DstAccount) error {
	if srcBalance < 0 {
		return errors.New("prohibit negative balance")
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("ERROR - SplitBalances: Begin tx : %s", err.Error())
	}

	defer func() {
		msg := "SUCCESS"
		if err != nil {
			tx.Rollback()
			msg = "CANCEL"
		}
		tx.Commit()

		// logging tx
		// all or nothing
		for _, dst := range dstAccounts {
			TxLogging(msg, src, dst.AccountID, dst.Balance, db)
		}

	}()

	var movedBalance float64
	for _, dstAccount := range dstAccounts {
		err = MoveBalanceTx(src, dstAccount.AccountID, dstAccount.Balance, db, tx)

		if err != nil {
			return fmt.Errorf("ERROR - SplitBalances: splitting: %s", err.Error())
		}

		movedBalance += dstAccount.Balance
	}

	if srcBalance != movedBalance {
		err = fmt.Errorf("ERROR - SplitBalances: Mismatch balance movement: Expected %f, but got %f",
			srcBalance, movedBalance)
		return err
	}

	return nil
}

func TxLogging(status string, src string, dst string, amount float64, db *sql.DB) {
	stmt := `INSERT INTO transInfo (uuid,status,src,dst,amount,createdDt,lastUpdatedDt) `
	stmt += `VALUES (:UUID,:Status,:Src,:Dst,:Amount,:CreatedDt,:LastUpdatedDt)`

	now := time.Now()
	t := Transaction{
		UUID:          uuid.NewString(),
		Status:        status,
		Src:           src,
		Dst:           dst,
		Amount:        amount,
		CreatedDt:     now,
		LastUpdatedDt: now,
	}

	_, err := db.Exec(stmt, &t.UUID, t.Status,
		t.Src, t.Dst, t.Amount, t.CreatedDt, t.LastUpdatedDt)

	if err != nil {
		log.Fatal("ERROR - transaction log: ", err.Error(), "\n", t)
	}
}

func MakePrimaryAccount(accID string, userID string, db *sql.DB) {
	stmt := `SELECT * from account WHERE accountID=:accID AND userID=:userID`
	row := db.QueryRow(stmt, accID, userID)

	tx, err := db.Begin()
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// check if designated account soon-to-become-primary exists
	var acc Account
	if err = row.Scan(&acc.AccountID, &acc.UserID, &acc.AccType,
		&acc.Balance, &acc.CreatedDt, &acc.LastUpdatedDt, &acc.IsPrimary); err != nil {
		if err != sql.ErrNoRows {
			// error other than account not found
			log.Fatal("ERROR - make primary account(Query): ", err.Error())
		}

		// if not found: operation is NOT allow
		fmt.Println("Make Primary Account: Operation NOT allowed: ",
			c.RecordNotFound{Id: accID, TbName: "account"}.Error())
		return
	}

	// select all marked primary account
	// SHOULD be ONLY 1 primary account
	stmt = `SELECT * from account WHERE userID=:userID AND isPrimary=1`
	rows, err := db.Query(stmt, userID)
	if err != nil {
		if err != sql.ErrNoRows {
			log.Fatal("ERROR - make primary account(Query primary): ", err.Error())
		}

		log.Fatal("ERROR - make primary account(Query primary): NO Primary account found!!")
	}
	defer rows.Close()

	// reset all primary to non-primary account
	stmt = `UPDATE account set isPrimary=0, lastUpdatedDt=:now WHERE userID=:UserID AND accountID=:AccountID`
	for rows.Next() {
		var pAcc Account
		if err = rows.Scan(&pAcc.AccountID, &pAcc.UserID, &pAcc.AccType,
			&pAcc.Balance, &pAcc.CreatedDt, &pAcc.LastUpdatedDt, &pAcc.IsPrimary); err != nil {
			log.Fatal("ERROR - make primary account(Scan primary): ", err.Error())
		}

		now := time.Now()
		result, err := tx.Exec(stmt, now, pAcc.UserID, pAcc.AccountID)
		if err != nil {
			log.Fatal("ERROR - make primary account(Update NonPrimary): ", err.Error())
		}

		nrow, err := result.RowsAffected()
		if err != nil {
			log.Fatal("ERROR - make primary account(Affected Update Nonprimary): ", err.Error())
		}

		fmt.Println("Affected row(s): ", nrow)
		fmt.Printf("Make primary account(Update to Non primary): %s\n", pAcc.AccountID)
	}

	// update designated account to primary account
	now := time.Now()
	stmt = `UPDATE account set isPrimary=1, lastUpdatedDt=:now WHERE userID=:userID AND accountID=:accID`
	result, err := tx.Exec(stmt, now, userID, accID)
	if err != nil {
		log.Fatal("ERROR - make primary account(Update Primary): ", err.Error())
	}

	nrow, err := result.RowsAffected()
	if err != nil {
		log.Fatal("ERROR - make primary account(Affected Update primary): ", err.Error())
	}

	fmt.Println("Affected row(s): ", nrow)
	fmt.Printf("Make primary account(Update to primary): %s\n", accID)
}
