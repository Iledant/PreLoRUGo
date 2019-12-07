package actions

import (
	"database/sql"
	"net/http"

	"github.com/Iledant/PreLoRUGo/models"
	"github.com/kataras/iris"
)

// setDBMiddleware return a middleware to add db to context values
func setDBMiddleware(db *sql.DB, superAdminEmail string) func(iris.Context) {
	return func(ctx iris.Context) {
		ctx.Values().Set("db", db)
		ctx.Values().Set("superAdminEmail", superAdminEmail)
		ctx.Next()
	}
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
	Masks: []int64{
		models.SuperAdminBit,
		models.ActiveAdminMask,
		models.ActiveCoproMask,
	},
	Message: "Droits sur les copropriétés requis",
}

var coproPreProgHandler = RightHandler{
	Masks: []int64{
		models.SuperAdminBit,
		models.ActiveAdminMask,
		models.ActiveCoproPreProgMask,
	},
	Message: "Droits préprogrammation sur les copropriétés requis",
}

var rpHandler = RightHandler{
	Masks: []int64{
		models.SuperAdminBit,
		models.ActiveAdminMask,
		models.ActiveRenewProjectMask,
	},
	Message: "Droits sur les projets RU requis",
}

var rpPreProgHandler = RightHandler{
	Masks: []int64{
		models.SuperAdminBit,
		models.ActiveAdminMask,
		models.ActiveRenewProjectPreProgMask,
	},
	Message: "Droits préprogrammation sur les projets RU requis",
}

var housingHandler = RightHandler{
	Masks: []int64{
		models.SuperAdminBit,
		models.ActiveAdminMask,
		models.ActiveHousingMask,
	},
	Message: "Droits sur les projets logement requis",
}

var housingPreProgHandler = RightHandler{
	Masks: []int64{
		models.SuperAdminBit,
		models.ActiveAdminMask,
		models.ActiveHousingPreProgMask,
	},
	Message: "Droits préprogrammation sur les projets logement requis",
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
