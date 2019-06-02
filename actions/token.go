package actions

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Iledant/PreLoRUGo/models"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/kataras/iris"
)

// UserClaims hold token user fields to avoid fetching database
type UserClaims struct {
	Rights int64
}

// customClaims add role and active to token to avoid fetching database
type customClaims struct {
	Rights int64 `json:"rig"`
	jwt.StandardClaims
}

var (
	signingKey   = []byte(os.Getenv("JWT_SIGNING_KEY"))
	expireDelay  = time.Second * 30
	refreshDelay = int64((time.Hour * 15 * 24).Seconds())
	iss          = "https://www.propera.net"
	tokens       = map[int]bool{}
	// ErrNoToken happens when header have no or bad authorization bearer
	ErrNoToken = errors.New("Token absent")
	// ErrBadToken happends when bearer token can't be verified
	// or isn't already stored after login or refresh
	ErrBadToken = errors.New("Token invalide")
)

// getTokenString store claims and return JWT token string
func getTokenString(claims *customClaims) (tokenString string, err error) {
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	if tokenString, err = token.SignedString(signingKey); err != nil {
		return "", err
	}
	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return "", err
	}
	tokens[userID] = true
	return tokenString, nil
}

// setToken creates or update a token for a given user
func setToken(u *models.User) (string, error) {
	t := time.Now()
	claims := customClaims{
		Rights: u.Rights,
		StandardClaims: jwt.StandardClaims{
			Subject:   strconv.FormatInt(u.ID, 10),
			ExpiresAt: t.Add(expireDelay).Unix(),
			IssuedAt:  t.Unix(),
			Issuer:    iss}}
	return getTokenString(&claims)
}

// delToken remove user ID from list of stored tokens
func delToken(userID int) {
	mutex := &sync.Mutex{}
	mutex.Lock()
	defer mutex.Unlock()
	delete(tokens, userID)
}

// refreshToken replace an existing expired token and add it to the response header
func refreshToken(ctx iris.Context, u *customClaims) error {
	t := time.Now()
	u.ExpiresAt = t.Add(expireDelay).Unix()
	u.IssuedAt = t.Unix()
	tokenString, err := getTokenString(u)
	if err != nil {
		return err
	}
	ctx.Header("Authorization", "Bearer "+tokenString)
	ctx.Header("Access-Control-Expose-Headers", "Authorization")
	return nil
}

// bearerToUser gets user claims (ID, role, active) from token in request header
// and send refreshed token if first time expired
func bearerToUser(ctx iris.Context) (claims *customClaims, err error) {
	bearer := ctx.GetHeader("Authorization")
	if len(bearer) < 8 {
		return nil, ErrNoToken
	}
	tokenString := strings.TrimPrefix(bearer, "Bearer ")
	if tokenString == "" {
		return nil, ErrNoToken
	}
	parser := jwt.Parser{ValidMethods: nil, UseJSONNumber: true,
		SkipClaimsValidation: true}
	token, err := parser.ParseWithClaims(tokenString, &customClaims{},
		func(token *jwt.Token) (interface{}, error) { return []byte(signingKey), nil })
	if err != nil || !token.Valid {
		return nil, ErrBadToken
	}
	claims = token.Claims.(*customClaims)
	// Check if previously connected
	userID, _ := strconv.Atoi(claims.Subject)
	mutex := &sync.Mutex{}
	mutex.Lock()
	_, ok := tokens[userID]
	mutex.Unlock()
	if !ok {
		return nil, ErrBadToken
	}
	// Refresh if expired
	t := time.Now().Unix()
	if t > claims.IssuedAt+refreshDelay {
		return claims, errors.New("Token expiré")
	}
	if t > claims.ExpiresAt {
		err = refreshToken(ctx, claims)
	}
	ctx.Values().Set("userID", userID)
	ctx.Values().Set("rights", claims.Rights)
	return claims, err
}

// isActive check an existing token in header and, if succeed,
// parse returning user active field
func isActive(ctx iris.Context) (bool, error) {
	u, err := bearerToUser(ctx)
	if err != nil {
		return false, err
	}
	return u.Rights&models.ActiveBit != 0 || u.Rights&models.SuperAdminBit != 0, nil
}

// isCopro check an existing token in header and, if succeed,
// check if user is admin, superadmin or has copro rights
func isCopro(ctx iris.Context) (bool, error) {
	u, err := bearerToUser(ctx)
	if err != nil {
		return false, err
	}
	return u.Rights&models.ActiveCoproMask == models.ActiveCoproMask ||
		u.Rights&models.ActiveAdminMask == models.ActiveAdminMask ||
		u.Rights&models.SuperAdminBit != 0, nil
}

// isRenewProject check an existing token in header and, if succeed,
// check if user is admin, superadmin or has renew project rights
func isRenewProject(ctx iris.Context) (bool, error) {
	u, err := bearerToUser(ctx)
	if err != nil {
		return false, err
	}
	return u.Rights&models.ActiveRenewProjectMask == models.ActiveRenewProjectMask ||
		u.Rights&models.ActiveAdminMask == models.ActiveAdminMask ||
		u.Rights&models.SuperAdminBit != 0, nil
}

// isHousing check an existing token in header and, if succeed,
// check if user is admin, superadmin or has housing rights
func isHousing(ctx iris.Context) (bool, error) {
	u, err := bearerToUser(ctx)
	if err != nil {
		return false, err
	}
	return u.Rights&models.ActiveHousingMask == models.ActiveHousingMask ||
		u.Rights&models.ActiveAdminMask == models.ActiveAdminMask ||
		u.Rights&models.SuperAdminBit != 0, nil
}

// isAdmin check an existing token in header and, if succeed,
// parse check if user active and admin
func isAdmin(ctx iris.Context) (bool, error) {
	u, err := bearerToUser(ctx)
	if err != nil {
		return false, err
	}
	return (u.Rights&models.ActiveAdminMask == models.ActiveAdminMask) ||
		u.Rights&models.SuperAdminBit != 0, nil
}

// isObserver check an existing token in header and, if succeed,
// parse check if user active and observer
func isObserver(ctx iris.Context) (bool, error) {
	u, err := bearerToUser(ctx)
	if err != nil {
		return false, err
	}
	return u.Rights&models.ActiveObserverMask == models.ActiveObserverMask, nil
}

// AdminMiddleware checks if there's a token and if it belongs to admin user
//  otherwise prompt error
func AdminMiddleware(ctx iris.Context) {
	admin, err := isAdmin(ctx)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{err.Error()})
		ctx.StopExecution()
		return
	}
	if !admin {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(jsonError{"Droits administrateur requis"})
		ctx.StopExecution()
		return
	}
	ctx.Next()
}

// ActiveMiddleware checks if there's a valid token and user is active otherwise prompt error
func ActiveMiddleware(ctx iris.Context) {
	active, err := isActive(ctx)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{err.Error()})
		ctx.StopExecution()
		return
	}
	if !active {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(jsonError{"Connexion requise"})
		ctx.StopExecution()
		return
	}
	ctx.Next()
}

// CoproMiddleware checks if there's a valid token and user is active and has
// copro rights otherwise prompt error
func CoproMiddleware(ctx iris.Context) {
	copro, err := isCopro(ctx)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{err.Error()})
		ctx.StopExecution()
		return
	}
	if !copro {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(jsonError{"Droits sur les copropriétés requis"})
		ctx.StopExecution()
		return
	}
	ctx.Next()
}

// RenewProjectMiddleware checks if there's a valid token and user is active and
// has rights on renew projects otherwise prompt error
func RenewProjectMiddleware(ctx iris.Context) {
	renewProject, err := isRenewProject(ctx)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{err.Error()})
		ctx.StopExecution()
		return
	}
	if !renewProject {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(jsonError{"Droits sur les projets RU requis"})
		ctx.StopExecution()
		return
	}
	ctx.Next()
}

// HousingMiddleware checks if there's a valid token and user is active and
// has rights on renew projects otherwise prompt error
func HousingMiddleware(ctx iris.Context) {
	renewProject, err := isHousing(ctx)
	if err != nil {
		ctx.StatusCode(http.StatusInternalServerError)
		ctx.JSON(jsonError{err.Error()})
		ctx.StopExecution()
		return
	}
	if !renewProject {
		ctx.StatusCode(http.StatusUnauthorized)
		ctx.JSON(jsonError{"Droits sur les projets logement requis"})
		ctx.StopExecution()
		return
	}
	ctx.Next()
}

// TokenRecover tries to load a previously saved file with tokens history.
// Used to allow users keep beeing logged in even after a relaunch of server.
func TokenRecover(fileName string) {
	fileContent, err := ioutil.ReadFile(fileName)
	if err != nil {
		return
	}
	if err = json.Unmarshal(fileContent, &tokens); err != nil {
		return
	}
}

// TokenSave saves the current map of tokens to a file in order to persist them
// for the next call of TokenRecover
func TokenSave(fileName string) {
	jsonTokens, err := json.Marshal(tokens)
	if err != nil {
		return
	}
	err = ioutil.WriteFile(fileName, jsonTokens, os.ModePerm)
}
