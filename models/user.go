package models

import (
	"database/sql"
	"errors"

	"golang.org/x/crypto/bcrypt"
)

// User model
type User struct {
	ID       int64  `json:"ID"`
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Password string `json:"-"`
	Role     string `json:"Role"`
	Active   bool   `json:"Active"`
}

// Users embeddes an array of User for json export.
type Users struct {
	Users []User `json:"User"`
}

const (
	// AdminRole defines value of role row in users table for an admin
	AdminRole = "ADMIN"
	// ObserverRole defines value of role row in users table for an observer
	ObserverRole = "OBSERVER"
	// UserRole defines value of role row in users table for an usual user
	UserRole = "USER"
)

// GetByID fetches a user from database using ID.
func (u *User) GetByID(db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT id, name, email, role, password, active 
	FROM users WHERE id = $1 LIMIT 1`, u.ID).Scan(&u.ID,
		&u.Name, &u.Email, &u.Role, &u.Password, &u.Active)
	return err
}

// CryptPwd crypt not codded password field.
func (u *User) CryptPwd() (err error) {
	cryptPwd, err := bcrypt.GenerateFromPassword([]byte(u.Password), 10)
	if err != nil {
		return err
	}
	u.Password = string(cryptPwd)
	return nil
}

// ValidatePwd compared sent uncodded password with internal password.
func (u *User) ValidatePwd(pwd string) error {
	return bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(pwd))
}

// GetAll fetches all users from database.
func (users *Users) GetAll(db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id, name, email, role, active FROM users`)
	if err != nil {
		return err
	}
	var r User
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&r.ID, &r.Name, &r.Email, &r.Role, &r.Active); err != nil {
			return err
		}
		users.Users = append(users.Users, r)
	}
	err = rows.Err()
	if len(users.Users) == 0 {
		users.Users = []User{}
	}
	return err
}

//GetRole fetches all users according to a role.
func (users *Users) GetRole(role string, db *sql.DB) (err error) {
	rows, err := db.Query(`SELECT id, name, email, role, active 
	FROM users WHERE role = $1`, role)
	if err != nil {
		return err
	}
	var r User
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&r.ID, &r.Name, &r.Email, &r.Role, &r.Active); err != nil {
			return err
		}
		users.Users = append(users.Users, r)
	}
	err = rows.Err()
	if len(users.Users) == 0 {
		users.Users = []User{}
	}
	return err
}

// GetByEmail fetches an user by email.
func (u *User) GetByEmail(email string, db *sql.DB) (err error) {
	err = db.QueryRow(`SELECT id, name, email, role, 
	password, active FROM users WHERE email = $1 LIMIT 1`, email).Scan(&u.ID,
		&u.Name, &u.Email, &u.Role, &u.Password, &u.Active)
	return err
}

// Exists checks if name or email is already in database.
func (u *User) Exists(db *sql.DB) error {
	var count int64
	if err := db.QueryRow(`SELECT count(1) FROM users WHERE email=$1 OR name=$2`,
		u.Email, u.Name).Scan(&count); err != nil {
		return err
	}
	if count != 0 {
		return errors.New("Utilisateur existant")
	}
	return nil
}

// Create insert a new user into database updating time fields.
func (u *User) Create(db *sql.DB) (err error) {
	err = db.QueryRow(`INSERT INTO users (name, email, 
		password, role, active) VALUES($1,$2,$3,$4,$5) RETURNING id`,
		u.Name, u.Email, u.Password, u.Role, u.Active).Scan(&u.ID)
	return err
}

// Update modifies a user into database.
func (u *User) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE users SET name=$1, email=$2, 
	password=$3, role=$4, active=$5 WHERE id=$6 `, u.Name, u.Email, u.Password,
		u.Role, u.Active, u.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Utilisateur introuvable")
	}
	return err
}

// Delete removes a user from database.
func (u *User) Delete(db *sql.DB) (err error) {
	res, err := db.Exec("DELETE FROM users WHERE id = $1", u.ID)
	if err != nil {
		return err
	}
	count, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if count != 1 {
		return errors.New("Utilisateur introuvable")
	}
	return nil
}
