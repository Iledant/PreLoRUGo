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

// RightHandler is used by RightsMiddleWare to handle permissions
// If none of the Masks matches with the user's rights, the Messages is used
// to send en error back
type RightHandler struct {
	Masks   []int64
	Message string
}

var admHandler = RightHandler{
	Masks:   []int64{models.SuperAdminBit, models.ActiveAdminMask},
	Message: "Droits administrateur requis",
}

var coproHandler = RightHandler{
	Masks:   []int64{models.SuperAdminBit, models.ActiveAdminMask, models.ActiveCoproMask},
	Message: "Droits sur les copropriétés requis",
}

var rpHandler = RightHandler{
	Masks:   []int64{models.SuperAdminBit, models.ActiveAdminMask, models.ActiveRenewProjectMask},
	Message: "Droits sur les projets RU requis",
}

var housingHandler = RightHandler{
	Masks:   []int64{models.SuperAdminBit, models.ActiveAdminMask, models.ActiveHousingMask},
	Message: "Droits sur les projets logement requis",
}

var userHandler = RightHandler{
	Masks:   []int64{models.SuperAdminBit, models.ActiveBit},
	Message: "Connexion requise",
}

// RightsMiddleWare checks if the user attached to the token match with the bit
// rights sent
func RightsMiddleWare(r *RightHandler) func(iris.Context) {
	return func(ctx iris.Context) {
		u, err := bearerToUser(ctx)
		if err != nil {
			ctx.StatusCode(http.StatusInternalServerError)
			ctx.JSON(jsonError{err.Error()})
			ctx.StopExecution()
			return
		}
		rights := true
		for _, mask := range r.Masks {
			rights = u.Rights&mask == mask
			if rights {
				break
			}
		}
		if !rights {
			ctx.StatusCode(http.StatusUnauthorized)
			ctx.JSON(jsonError{r.Message})
			ctx.StopExecution()
			return
		}
		ctx.Next()
	}
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
