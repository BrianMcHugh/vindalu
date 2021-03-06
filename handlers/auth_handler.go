package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/context"

	"github.com/vindalu/vindalu/auth"
)

func (ir *VindaluApiHandler) AuthTokenHandler(w http.ResponseWriter, r *http.Request) {

	user := context.Get(r, Username).(string)
	isAdmin := context.Get(r, IsAdmin).(bool)
	cfg := ir.Config()

	token := auth.GetNewToken(user, cfg.Auth.Token.TTL)
	token.Claims["admin"] = isAdmin

	tokenStr, err := token.SignedString([]byte(cfg.Auth.Token.SigningKey))
	if err != nil {
		ir.writeAndLogResponse(w, r, 500, map[string]string{"Content-Type": "text/plain"}, []byte(err.Error()))
		return
	}

	ir.writeAndLogResponse(w, r, 200, map[string]string{"Content-Type": "application/json"},
		[]byte(fmt.Sprintf(`{"token": "%s", "admin": %v}`, tokenStr, isAdmin)))
}
