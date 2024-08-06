package asynqmonauth

import (
	"crypto/sha256"
	"crypto/subtle"
	"fmt"
	"net/http"
)

func redirectToRootHandler(rootPath string) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("Redirecting to %s\n", rootPath)
		http.Redirect(w, r, rootPath, http.StatusPermanentRedirect)
	})
}

func basicAuthHandler(next http.HandlerFunc, auth *AuthBasic) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		username, password, ok := r.BasicAuth()
		if ok {
			if verifyAuth(auth.Username, username) && verifyAuth(auth.Password, password) {
				next.ServeHTTP(w, r)
				return
			}
		}

		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	})
}

func verifyAuth(raw string, hashed string) bool {
	hashedHash := sha256.Sum256([]byte(hashed))
	rawHash := sha256.Sum256([]byte(raw))

	return subtle.ConstantTimeCompare(hashedHash[:], rawHash[:]) == 1
}
