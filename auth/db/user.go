package db

import (
	"backend/utils"
	"database/sql"
	"encoding/json"
	"io"

	"github.com/go-playground/validator"
	"github.com/google/uuid"
)

type User struct {
	Id        uuid.UUID `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt string    `json:"-"`
}

func (p *User) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(p)
}

func (user *User) Validate() error {
	validate := validator.New()
	return validate.Struct(user)
}

func (user *User) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(user)
}

func (db *UsersDb) AddUser(user *User) error {
	hash, err := utils.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = hash
	var id uuid.UUID
	query := `INSERT INTO users (username, email, password) VALUES ($1, $2, $3) RETURNING id`
	err = db.Conn.QueryRow(query, user.Username, user.Email, user.Password).Scan(&id)
	if err != nil {
		return err
	}
	user.Id = id
	return nil
}

func (db *UsersDb) GetUserById(userId int) (User, error) {
	user := User{}
	query := `SELECT * FROM users WHERE id = $1;`
	row := db.Conn.QueryRow(query, userId)
	switch err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password); err {
	case sql.ErrNoRows:
		return user, ErrNoMatch
	default:
		return user, err
	}
}

func (db *UsersDb) DeleteUser(userId int) error {
	query := `DELETE FROM users WHERE id = $1;`
	_, err := db.Conn.Exec(query, userId)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}

func (db *UsersDb) UpdateUser(user *User) error {
	hash, err := utils.Hash(user.Password)
	if err != nil {
		return err
	}
	user.Password = hash
	query := `UPDATE users SET username = $1, email = $2, password = $3 WHERE id = $4;`
	_, err = db.Conn.Exec(query, user.Username, user.Email, user.Password, user.Id)
	switch err {
	case sql.ErrNoRows:
		return ErrNoMatch
	default:
		return err
	}
}

func (db *UsersDb) GetUser(email, password string) (User, error) {
	user := User{}
	dbpass, err := utils.Hash(password)
	if err != nil {
		return user, err
	}
	query := `SELECT * FROM users WHERE email = $1 AND password = $2;`
	row := db.Conn.QueryRow(query, email, dbpass)
	switch err := row.Scan(&user.Id, &user.Username, &user.Email, &user.Password, &user.CreatedAt); err {
	case sql.ErrNoRows:
		return user, ErrNoMatch
	default:
		return user, err
	}
}
