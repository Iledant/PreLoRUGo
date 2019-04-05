package actions

import (
	"database/sql"
	"errors"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// signInResp is used to send a JSON response after signing in
type signInResp struct {
	Token string      `json:"Token"`
	User  models.User `json:"User"`
}

type userResp struct {
	User models.User `json:"User"`
}

// sentUser is used to create or update an user taking into account that
// models.User doesn't JSON export the password
type sentUser struct {
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
	Rights   int64  `json:"Rights"`
}

// credentials is used to decode user login payload
type credentials struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

// Login handles user login using credentials and return token if success.
func Login(ctx iris.Context) {
	var c credentials
	if err := ctx.ReadJSON(&c); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Décodage login : " + err.Error()})
		return
	}
	if c.Email == "" || c.Password == "" {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Champ manquant ou incorrect"})
		return
	}
	var user models.User
	db := ctx.Values().Get("db").(*sql.DB)
	if err := user.GetByEmail(c.Email, db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{err.Error()})
		return
	}
	if err := user.ValidatePwd(c.Password); err != nil {
		ctx.StatusCode(http.StatusNotFound)
		ctx.JSON(jsonError{"Erreur de login ou mot de passe"})
		return
	}
	token, err := setToken(&user)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(signInResp{token, user})
}

// Logout handles users logout and destroy his token.
func Logout(ctx iris.Context) {
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Logout, impossible de récupérer l'ID"})
	}
	delToken(int(userID))
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonError{"Utilisateur déconnecté"})
}

// GetUsers handles the GET request for all users and send back only secure fields.
func GetUsers(ctx iris.Context) {
	var users models.Users
	db := ctx.Values().Get("db").(*sql.DB)
	if err := users.GetAll(db); err != nil {
		ctx.JSON(jsonError{"Liste des utilisateurs : " + err.Error()})
		ctx.StatusCode(http.StatusInternalServerError)
		return
	}
	ctx.JSON(users)
	ctx.StatusCode(http.StatusOK)
}

// CreateUser handles the creation by admin of a new user and returns the created user.
func CreateUser(ctx iris.Context) {
	var req sentUser
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création d'utilisateur, décodage : " + err.Error()})
		return
	}
	user := models.User{Name: req.Name, Email: req.Email, Password: req.Password,
		Rights: req.Rights}
	if err := user.Validate(); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création d'utilisateur : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	if err := user.Exists(db); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Création d'utilisateur : " + err.Error()})
		return
	}
	if err := user.CryptPwd(); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création d'utilisateur, cryptage : " + err.Error()})
		return
	}
	if err := user.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Création d'utilisateur, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(userResp{user})
}

// UpdateUser handles the updating by admin of an existing user and sent back modified user.
func UpdateUser(ctx iris.Context) {
	userID, err := ctx.Params().GetInt64("userID")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Modification d'utilisateur, paramètre : " + err.Error()})
		return
	}
	db, user := ctx.Values().Get("db").(*sql.DB), models.User{ID: userID}
	if err = user.GetByID(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification d'utilisateur, get : " + err.Error()})
		return
	}
	var req sentUser
	if err = ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification d'utilisateur, décodage : " + err.Error()})
		return
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Name != "" {
		user.Name = req.Name
	}
	if req.Password != "" {
		user.Password = req.Password
		if err = user.CryptPwd(); err != nil {
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.JSON(jsonError{"Modification d'utilisateur, mot de passe : " + err.Error()})
			return
		}
	}
	user.Rights = req.Rights
	if err = user.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Modification d'utilisateur, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(userResp{user})
}

// DeleteUser handles the deleting by admin of an existing user.
func DeleteUser(ctx iris.Context) {
	userID, err := ctx.Params().GetInt64("userID")
	if err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Suppression d'utilisateur, paramètre : " + err.Error()})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	user := models.User{ID: userID}
	if err = user.Delete(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Suppression d'utilisateur, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonError{"Utilisateur supprimé"})
}

type signUpReq struct {
	Name     string `json:"Name"`
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

// SignUp handles the request of a new user and creates an inactive account.
func SignUp(ctx iris.Context) {
	var req signUpReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Inscription d'utilisateur, décodage : " + err.Error()})
		return
	}
	if req.Name == "" || req.Email == "" || req.Password == "" {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Inscription d'utilisateur : Champ manquant ou incorrect"})
		return
	}
	db := ctx.Values().Get("db").(*sql.DB)
	user := models.User{Name: req.Name,
		Email:    req.Email,
		Password: req.Password}
	if err := user.Exists(db); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Inscription d'utilisateur, exists : " + err.Error()})
		return
	}
	if err := user.CryptPwd(); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Inscription d'utilisateur, password : " + err.Error()})
		return
	}
	if err := user.Create(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Inscription d'utilisateur, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusCreated)
	ctx.JSON(jsonError{"Utilisateur créé, en attente d'activation"})
}

type chgPwdReq struct {
	CurrentPassword string `json:"CurrentPassword"`
	NewPassword     string `json:"NewPassword"`
}

// ChangeUserPwd handles the request of a user to change his password.
func ChangeUserPwd(ctx iris.Context) {
	var req chgPwdReq
	if err := ctx.ReadJSON(&req); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Changement de mot de passe, décodage : " + err.Error()})
		return

	}
	if req.CurrentPassword == "" || req.NewPassword == "" {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Changement de mot de passe : ancien et nouveau mots de passe requis"})
		return
	}
	userID, err := getUserID(ctx)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Changement de mot de passe, user ID : " + err.Error()})
		return
	}
	db, user := ctx.Values().Get("db").(*sql.DB), models.User{ID: userID}
	if err = user.GetByID(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Changement de mot de passe, get : " + err.Error()})
		return
	}
	if err = user.ValidatePwd(req.CurrentPassword); err != nil {
		ctx.StatusCode(http.StatusBadRequest)
		ctx.JSON(jsonError{"Changement de mot de passe : erreur de mot de passe"})
		return
	}
	user.Password = req.NewPassword
	if err = user.CryptPwd(); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Changement de mot de passe, password : " + err.Error()})
		return
	}
	if err = user.Update(db); err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{"Changement de mot de passe, requête : " + err.Error()})
		return
	}
	ctx.StatusCode(http.StatusOK)
	ctx.JSON(jsonError{"Mot de passe changé"})
}

// getUserRoleAndID fetch user role and ID with the token
func getUserID(ctx iris.Context) (uID int64, err error) {
	u := ctx.Values().Get("userID")
	r := ctx.Values().Get("rights")
	if u == nil || r == nil {
		return 0, errors.New("Utilisateur non enregistré")
	}
	uID = int64(u.(int))
	rights := r.(int64)
	if rights&models.SuperAdminBit != 0 || rights&models.AdminBit != 0 {
		uID = 0
	}
	return uID, nil
}

type setPwdReq struct {
	Password string `json:"Password"`
}
