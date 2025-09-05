package postgres

import (
	"os"
	"testing"
)

func setupDB(t *testing.T, db *DB) {
	query, err := os.ReadFile("./database_setup/database.up.sql")
	if err != nil {
		t.Fatal(err)
	}
	if err = db.Exec(string(query)).Error; err != nil {
		t.Fatal(err)
	}
}

func teardownDB(t *testing.T, db *DB) {
	err := db.Exec("DROP TABLE exchanges; DROP TABLE valutes;").Error
	if err != nil {
		t.Fatal(err)
	}
}
