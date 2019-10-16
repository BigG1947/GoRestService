package db

import (
	"database/sql"
	"unicode"
)
import _ "github.com/mattn/go-sqlite3"

type User struct {
	Id       int64  `json:"id"`
	FullName string `json:"full_name"`
	About    string `json:"about"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

func ConnectionToDb() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "userInfo.db")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec(createUserTable)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func (u *User) Add(db *sql.DB) (int64, error) {
	res, err := db.Exec(addUserScript, u.FullName, u.About, u.Login, u.Password)
	if err != nil {
		return 0, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (u *User) Update(db *sql.DB) error {
	_, err := db.Exec(updateUserScript, u.FullName, u.About, u.Login, u.Password, u.Id)
	return err
}

func (u *User) Delete(db *sql.DB) error {
	_, err := db.Exec(deleteUserScript, u.Id)
	return err
}

func GetAllUsers(db *sql.DB) ([]User, error) {
	var users []User

	rows, err := db.Query(getAllUserScript)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.FullName, &user.About, &user.Login, &user.Password)
		if err != nil {
			return nil, err
		}

		users = append(users, user)
	}

	return users, nil
}

func (u *User) GetUserById(db *sql.DB, id int64) (bool, error) {
	err := db.QueryRow(getUserByIdScript, id).Scan(&u.Id, &u.FullName, &u.About, &u.Login, &u.Password)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func CheckUserExist(db *sql.DB, login string) (int64, bool, error) {
	var id int64
	err := db.QueryRow(checkUserExistScript, login).Scan(&id)

	if err == sql.ErrNoRows {
		return 0, false, nil
	} else if err != nil {
		return 0, false, err
	}

	return id, true, nil
}

func (u *User) CheckAuthData(login string, password string) bool {
	if login == u.Login && password == u.Password {
		return true
	}
	return false
}

func (u *User) ValidateData() bool {
	if u.Login == "" || u.Password == "" || u.About == "" || u.FullName == "" {
		return false
	}

	for i := range u.Login {
		if unicode.IsUpper(rune(u.Login[i])) ||
			!(unicode.IsDigit(rune(u.Login[i])) || unicode.IsLetter(rune(u.Login[i])) || u.Login[i] == '_') {
			return false
		}
	}
	return true
}
