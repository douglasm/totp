package users

import (
	"errors"

	"github.com/gorilla/schema"

	"totp/types"
)

var (
	logins = map[string]types.LoginRec{}

	decoder = schema.NewDecoder()

	errPassUsed = errors.New("That password has been revealed in a data breach")
)

func init() {
	theLogin := types.LoginRec{Mail: "fred@fr4edd", Pass: "all hands on deck", Cookie: "OqOkNkwoMzeINVkZ", Secret: "FKB245SRZQYWG66NOOWAICIGEHSHIJJT"}
	logins["OqOkNkwoMzeINVkZ"] = theLogin
}

func GetLogin(cookie string) types.LoginRec {
	theUser := logins[cookie]
	return theUser
}
