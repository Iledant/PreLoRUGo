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
	Rights   int64  `json:"Rights"`
}

// Users embeddes an array of User for json export.
type Users struct {
	Users []User `json:"User"`
}

const (
	// ActiveBit of user's rights field specifies the user is active
	ActiveBit = 1
	// SuperAdminBit of user's rights field specifies super user rights, used
	// at the beginning of the program and for admin purposes
	SuperAdminBit = 1 << 1
	// AdminBit of user's rights field specifies the user has access to all functions
	AdminBit = 1 << 2
	// CoproBit of user's rights field specifies the user has access to copro functions
	CoproBit = 1 << 3
	// RenewProjectBit of user's rights field specifies the user has access to renew project functions
	RenewProjectBit = 1 << 4
	// ObserverBit of user's rights field specifies the user has read access to functions
	ObserverBit = 1 << 5
	// HousingBit of user's rights field specifies the user has access to housing functions
	HousingBit = 1 << 6
	// PreProgBit of user's rights fields specifies the user can modify the preprogramming
	PreProgBit = 1 << 7
	// RightsMask is used to check if user's rights field are correctly filled
	RightsMask = ActiveBit | SuperAdminBit | AdminBit | CoproBit | RenewProjectBit | ObserverBit
	// ActiveAdminMask is used to check if a user is an active admin
	ActiveAdminMask = ActiveBit | AdminBit
	// ActiveObserverMask is used to check if a user is an active observer
	ActiveObserverMask = ActiveBit | ObserverBit
	// ActiveCoproMask is used to check is a user is active and has copro rights
	ActiveCoproMask = ActiveBit | CoproBit
	// ActiveCoproPreProgMask is used to check is a user is active and has copro and pre prog rights
	ActiveCoproPreProgMask = ActiveBit | CoproBit | PreProgBit
	// ActiveRenewProjectMask is used to check is a user is active and has renew project rights
	ActiveRenewProjectMask = ActiveBit | RenewProjectBit
	// ActiveRenewProjectPreProgMask is used to check is a user is active and has renew project and pre prog rights
	ActiveRenewProjectPreProgMask = ActiveBit | RenewProjectBit | PreProgBit
	// ActiveHousingMask is used to check is a user is active and has housing rights
	ActiveHousingMask = ActiveBit | HousingBit
	// ActiveHousingPreProgMask is used to check is a user is active and has housing and pre prog rights
	ActiveHousingPreProgMask = ActiveBit | HousingBit | PreProgBit
)

// Validate checks if field are correctly filled for database constraints
func (u *User) Validate() error {
	if u.Name == "" {
		return errors.New("Champ name vide")
	}
	if u.Email == "" {
		return errors.New("Champ email vide")
	}
	if u.Password == "" {
		return errors.New("Champ password vide")
	}
	return nil
}

// GetByID fetches a user from database using ID.
func (u *User) GetByID(db *sql.DB) error {
	return db.QueryRow(`SELECT id,name,email,password,rights FROM users
		WHERE id=$1`, u.ID).Scan(&u.ID, &u.Name, &u.Email, &u.Password, &u.Rights)
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
func (users *Users) GetAll(db *sql.DB, superAdminMail string) (err error) {
	rows, err := db.Query(`SELECT id, name, email, rights FROM users 
	WHERE email!=$1`, superAdminMail)
	if err != nil {
		return err
	}
	var r User
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&r.ID, &r.Name, &r.Email, &r.Rights); err != nil {
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
	err = db.QueryRow(`SELECT id, name, email, password, rights 
	FROM users WHERE email = $1 LIMIT 1`, email).Scan(&u.ID,
		&u.Name, &u.Email, &u.Password, &u.Rights)
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
func (u *User) Create(db *sql.DB) error {
	return db.QueryRow(`INSERT INTO users (name, email, 
		password, rights) VALUES($1,$2,$3,$4) RETURNING id`,
		u.Name, u.Email, u.Password, u.Rights).Scan(&u.ID)
}

// Update modifies a user into database.
func (u *User) Update(db *sql.DB) (err error) {
	res, err := db.Exec(`UPDATE users SET name=$1, email=$2, password=$3, 
	rights=$4 WHERE id=$5 `, u.Name, u.Email, u.Password, u.Rights, u.ID)
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
