package goPost05

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
)

func openConnection() (*sql.DB, error) {
	conn := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		Hostname,
		Port,
		Username,
		Password,
		Database)

	db, err := sql.Open("postgres", conn)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func exists(username string) int {
	username = strings.ToLower(username)

	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}

	defer db.Close()

	userID := -1
	statement := fmt.Sprintf(
		`SELECT "id" FROM "users" where username = '%s'`,
		username)

	rows, err := db.Query(statement)

	for rows.Next() {
		var id int
		err = rows.Scan(&id)
		if err != nil {
			fmt.Println("Scan", err)
			return -1
		}

		userID = id
	}

	defer rows.Close()
	return userID
}

func AddUser(d Userdata) int {
	d.Username = strings.ToLower(d.Username)
	db, err := openConnection()
	if err != nil {
		fmt.Println(err)
		return -1
	}
	defer db.Close()

	userID := exists(d.Username)
	if userID != -1 {
		fmt.Println("User already exists:", Username)
		return -1
	}

	insertStatement := `insert into "users" ("username") values ($1)`
	_, err = db.Exec(insertStatement, d.Username)

	if err != nil {
		fmt.Println(err)
		return -1
	}

	userID = exists(d.Username)
	if userID == -1 {
		return userID
	}

	insertStatement = `insert into "userdata" ("userid", "name", "surname", "description") values($1, $2, $3, $4)`
	_, err = db.Exec(insertStatement, userID, d.Name, d.Surname, d.Description)
	if err != nil {
		fmt.Println("db.Exec()", err)
		return -1
	}

	return userID
}

func DeleteUser(id int) error {
	db, err := openConnection()
	if err != nil {
		return err
	}

	defer db.Close()

	statement := fmt.Sprintf(
		`select "username" from "users" where id = %d`,
		id)
	rows, err := db.Query(statement)

	var username string
	for rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			return err
		}
	}

	defer rows.Close()

	if exists(username) != id {
		return fmt.Errorf("User with ID %d does not exist", id)
	}

	deleteQuery := `delete from "userdata" where userid=$1`
	_, err = db.Exec(deleteQuery, id)
	if err != nil {
		return err
	}

	deleteQuery = `delete from "users" where id=$1`
	_, err = db.Exec(deleteQuery, id)
	if err != nil {
		return err
	}

	return nil
}

func ListUsers() ([]Userdata, error) {
	Data := []Userdata{}
	db, err := openConnection()
	if err != nil {
		return Data, err
	}

	defer db.Close()

	rows, err := db.Query(`
		select "id","username","name","surname","description"
		from "users" 
		inner join "userdata" users.id = userdata.userid`)
	if err != nil {
		return Data, err
	}

	for rows.Next() {
		var id int
		var username string
		var name string
		var surname string
		var description string
		err = rows.Scan(&id, &username, &name, &surname, &description)
		temp := Userdata{
			ID:          id,
			Username:    username,
			Name:        name,
			Surname:     surname,
			Description: description,
		}

		Data = append(Data, temp)
	}

	defer rows.Close()
	return Data, nil
}

func UpdateUser(d Userdata) error {
	db, err := openConnection()
	if err != nil {
		return err
	}

	defer db.Close()

	userID := exists(d.Username)
	if userID == -1 {
		return errors.New("User does not exist")
	}

	d.ID = userID
	updateQuery := `update "userdata" set "name"=$1, "surname"=$2, "description"=$3 where "userid"=$4`
	_, err = db.Exec(updateQuery, d.Name, d.Surname, d.Description, d.ID)
	if err != nil {
		return err
	}

	return nil
}
