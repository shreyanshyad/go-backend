package middlewares

import (
	"backend/utils"
	"fmt"
	"net/http"
	"os"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/mux"
)

func validateJWT(token *jwt.Token) (interface{}, error) {
	if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
		fmt.Println("Unexpected signing method: ", token.Header["alg"])
		return nil, fmt.Errorf("invalid signing method")
	}

	aud := "samudai-dash"
	iss := "samudai-auth"

	checkAudience := token.Claims.(jwt.MapClaims).VerifyAudience(aud, false)
	if !checkAudience {
		fmt.Println("Invalid audience")
		return nil, fmt.Errorf("invalid aud")
	}

	checkIss := token.Claims.(jwt.MapClaims).VerifyIssuer(iss, false)
	if !checkIss {
		fmt.Println("Invalid issuer")
		return nil, fmt.Errorf("invalid iss")
	}

	key := os.Getenv("JWT_SECRET_KEY")
	fmt.Println("Key: ", key)
	return []byte(key), nil
}

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header["Authorization"] != nil {
			token, err := jwt.Parse(r.Header["Authorization"][0], validateJWT)
			if err != nil || !token.Valid {
				fmt.Println(err, token.Valid)
				utils.WriteFailureResponse(w, http.StatusUnauthorized, "Invalid Token")
				return
			}

			userId := token.Claims.(jwt.MapClaims)["sub"]
			mux.Vars(r)[KeyUser] = userId.(string)
			fmt.Println("User ID: ", userId)

			next.ServeHTTP(w, r)
		} else {
			utils.WriteFailureResponse(w, http.StatusUnauthorized, "No Authorization Token provided")
		}
	})
}
