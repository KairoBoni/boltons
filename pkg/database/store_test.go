package database

import (
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
)

func TestInsertNfeAmount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	mock.ExpectExec(`INSERT INTO nfe_amount`).
		WithArgs("test-acc", "client").
		WillReturnResult(sqlmock.NewResult(1, 1))

	s := &Store{sqlx.NewDb(db, "pq")}
	assert.Nil(t, s.InsertNfeAmount("test-acc", "client"))
	assert.Nil(t, mock.ExpectationsWereMet())
}

func TestGetAccount(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"amount"}).AddRow("190.00")

	q := `SELECT (.+) FROM nfe_amount WHERE (.+);`

	mock.ExpectQuery(q).
		WithArgs("1").
		WillReturnRows(rows)

	s := &Store{sqlx.NewDb(db, "pq")}

	amount, err := s.GetNfeAmount("1")

	assert.Equal(t, "190.00", amount)
	assert.Nil(t, err)
	assert.Nil(t, mock.ExpectationsWereMet())
}
