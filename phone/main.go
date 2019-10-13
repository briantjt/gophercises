package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const (
	host   = "localhost"
	port   = 5432
	dbname = "gopher_phone"
)

func main() {
	err := godotenv.Load()
	must(err)
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s sslmode=disable", host, port, user, password)
	db, err := sql.Open("postgres", psqlInfo)
	must(err)

	err = resetDB(db, dbname)
	must(err)
	db.Close()

	psqlInfo = fmt.Sprintf("%s dbname=%s", psqlInfo, dbname)
	db, err = sql.Open("postgres", psqlInfo)
	must(err)
	defer db.Close()
	must(createTable(db))
	err = seedDB(db)
	must(err)
	phones, err := getAllPhoneNumbers(db)
	must(err)
	fmt.Println(phones)
	for _, phone := range phones {
		number := normalize(phone.value)
		if number != phone.value {
			dup_phone, err := findPhoneByNumber(db, number)
			if dup_phone == nil && err == nil {
				err := updatePhoneById(db, phone.id, number)
				must(err)
			} else if dup_phone != nil && err == nil {
				err := deletePhoneById(db, phone.id)
				must(err)
			} else {
				log.Fatal(err)
			}
		}
	}
	phones, err = getAllPhoneNumbers(db)
	must(err)
	fmt.Println(phones)
}

type Phone struct {
	id int
	value string
}

func seedDB(db *sql.DB) error {
	_, err := db.Exec(`
	INSERT INTO phone_numbers(value) VALUES
('1234567890'),
('123 456 7891'),
('(123) 456 7892'),
('(123) 456-7893'),
('123-456-7894'),
('123-456-7890'),
('1234567892'),
('(123)456-7892');
	`)
	return err
}

func updatePhoneById(db *sql.DB, id int, number string) error {
	_, err := db.Exec("UPDATE phone_numbers SET value=$2 WHERE id=$1", id, number)
	return err
}

func findPhoneByNumber(db *sql.DB, number string) (*Phone, error) {
	var phone Phone
	row := db.QueryRow("SELECT id, value FROM phone_numbers WHERE value=$1", number)
	err := row.Scan(&phone.id, &phone.value)
	if err == sql.ErrNoRows {
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	return &phone, nil
}

func deletePhoneById(db *sql.DB, id int) error {
	_, err := db.Exec("DELETE FROM phone_numbers where id=$1", id)
	return err
}

func getAllPhoneNumbers(db *sql.DB) ([]Phone, error) {
	rows, err := db.Query("SELECT id, value FROM phone_numbers")
	if err != nil {
		return nil, err
	}
	var phones []Phone
	defer rows.Close()
	for rows.Next() {
		var phone Phone
		if err := rows.Scan(&phone.id, &phone.value); err != nil {
			return nil, err
		}
		phones = append(phones, phone)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return phones, nil
}

func insertPhoneNumber(db *sql.DB, phone string) (int, error) {
	statement := `
	INSERT INTO phone_numbers(value) VALUES($1)	RETURNING id
	`
	var id int
	err := db.QueryRow(statement, phone).Scan(&id)
	if err != nil {
		return -1, err
	}
	return id, err
}

func must(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func createTable(db *sql.DB) error {
	statement := `
	CREATE TABLE IF NOT EXISTS phone_numbers (
		id SERIAL,
		value VARCHAR(255)
	)	
	`
	_, err := db.Exec(statement)
	return err
}
func createDB(db *sql.DB, name string) error {
	_, err := db.Exec("CREATE DATABASE " + name)
	return err
}

func resetDB(db *sql.DB, name string) error {
	_, err := db.Exec("DROP DATABASE IF EXISTS " + name)
	if err != nil {
		return err
	}
	return createDB(db, name)
}

func normalize(phone string) string {
	var buff bytes.Buffer
	for _, char := range phone {
		if char >= '0' && char <= '9' {
			buff.WriteRune(char)
		}
	}
	return buff.String()
}
