package misc

import (
	"crypto/rand"
	"crypto/sha512"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"
)

type User struct {
	ID           int
	Username     string
	PasswordHash string
	PasswordSalt string
	FirstName    string
	LastName     string
}

type Session struct {
	UserID int
	Expire time.Time
}

func GenerateRandomSalt(saltSize int) (salt []byte, err error) {
	salt = make([]byte, saltSize)

	_, err = rand.Read(salt[:])

	return salt, err
}

func HashPassword(password string, salt []byte) string {
	var pw_bytes = []byte(password)
	pw_bytes = append(pw_bytes, salt...)

	var hash = sha512.New()
	hash.Write(pw_bytes)

	return hex.EncodeToString(hash.Sum(nil))
}

func PasswordMatchHash(password string, hash string, salt []byte) bool {
	var hashed_password = HashPassword(password, salt)

	return hash == hashed_password
}

func CreateCookie(r *http.Request, salt []byte) (cookie string, err error) {
	var cookie_bytes []byte
	user_agent := r.Header.Get("User-Agent")
	log.Printf("User-agent: %s", user_agent)
	cookie_bytes = append(cookie_bytes, []byte(user_agent)...)

	var hash = sha512.New()
	hash.Write(cookie_bytes)

	cookie = hex.EncodeToString(hash.Sum(nil))

	return
}

func VerifySessionCookie(r *http.Request, db *sql.DB) (authorized bool, err error) {
	cookie, err := r.Cookie("session")
	if err != nil {
		return false, err
	}
	var session Session
	rows, err := db.Query("SELECT user_id, expire FROM sessions WHERE cookie_hash = ?", cookie.Value)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		if err = rows.Scan(&session.UserID, &session.Expire); err != nil {
			continue
		}
		if time.Now().Before(session.Expire) {
			return true, nil
		}
	}
	return false, rows.Err()
}

func AddSessionCookie(w http.ResponseWriter, r *http.Request, db *sql.DB) (err error) {
	var cookie_value string
	cookie_value, err = CreateCookie(r, []byte("12345")) // TODO: make salt random
	if err != nil {
		return
	}
	cookie := &http.Cookie{
		Name:  "session",
		Value: cookie_value,
		Path:  "/",
		//			Secure:   true,
		//			HttpOnly: true,
	}
	http.SetCookie(w, cookie)

	var user User
	query := "INSERT INTO sessions (user_id, cookie_hash, cookie_salt, expire) VALUES (?, ?, ?, ?)"
	id := 0
	err = db.QueryRow(query, user.ID, cookie_value, "12345", time.Now().Add(24*time.Hour)).Scan(&id)
	// Ignore 'error'
	if err != nil && err != sql.ErrNoRows {
		return
	}
	return nil
}

func VerifyPOSTLogin(r *http.Request, db *sql.DB) (authorized bool, err error) {
	// If the user already have a cookie, ignore login info and return true.
	authorized, err = VerifySessionCookie(r, db)
	if authorized {
		return true, err
	}
	if err != nil {
		log.Printf("%s", err.Error())
	}

	if err = r.ParseForm(); err != nil {
		return false, err
	}
	username := r.FormValue("username")
	if len(username) < 1 {
		return false, errors.New("username is missing")
	}
	password := r.FormValue("password")
	if len(password) < 1 {
		return false, errors.New("password is missing")
	}
	// Debug
	fmt.Printf("Username: %s, Password: %s\n", username, password)

	var user User
	row := db.QueryRow("SELECT id, password_hash, password_salt FROM users WHERE username = ?", username)
	err = row.Scan(&user.ID, &user.PasswordHash, &user.PasswordSalt)
	if err != nil {
		return false, err
	}
	if !PasswordMatchHash(password, user.PasswordHash, []byte(user.PasswordSalt)) {
		return false, errors.New("wrong password")
	}
	return true, nil
}

func VerifyLogin(r *http.Request, db *sql.DB) (authorized bool, err error) {
	authorized = false // failsafe
	if r.Method == "POST" {
		return VerifyPOSTLogin(r, db)
	} else {
		return VerifySessionCookie(r, db)
	}
}
