package endpoints

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/sahib/config"
	"github.com/sahib/wishlist/cache"
	"github.com/sahib/wishlist/db"
)

type TokenHandler struct {
	db    *db.Database
	cache *cache.SessionCache
	cfg   *config.Config
}

func NewTokenHandler(db *db.Database, cache *cache.SessionCache, cfg *config.Config) *TokenHandler {
	return &TokenHandler{db: db, cache: cache, cfg: cfg}
}

func (th *TokenHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	token, ok := vars["token"]
	if !ok {
		jsonifyErrf(w, http.StatusBadRequest, "no token received")
		return
	}

	userID, err := th.cache.Confirm(token)
	if err != nil {
		jsonifyErrf(w, http.StatusInternalServerError, "failed to confirm token: %v", err)
		return
	}

	expireTime := time.Now().Add(th.cfg.Duration("auth.expire_time"))
	http.SetCookie(w, &http.Cookie{
		Name:    "session_id",
		Value:   token,
		Path:    "/",
		Expires: expireTime,
	})

	user, err := th.db.GetUserByID(userID)
	if err == nil {
		http.SetCookie(w, &http.Cookie{
			Name:    "user_name",
			Value:   user.Name,
			Path:    "/",
			Expires: expireTime,
		})
		http.SetCookie(w, &http.Cookie{
			Name:    "user_email",
			Value:   user.EMail,
			Path:    "/",
			Expires: expireTime,
		})
	}

	log.Printf("cookies set for %s (%s)", user.Name, user.EMail)
	http.Redirect(w, r, "/list.html", http.StatusSeeOther)
}

func (th *TokenHandler) NeedsAuthentication() bool {
	return false
}
