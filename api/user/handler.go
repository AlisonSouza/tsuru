package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

func CreateUser(w http.ResponseWriter, r *http.Request) error {
	var u User
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, &u)
	if err != nil {
		return err
	}
	err = u.Create()
	if err == nil {
		w.WriteHeader(http.StatusCreated)
		return nil
	}

	if u.Get() == nil {
		err = errors.New("This email is already registered")
	}

	return err
}

func Login(w http.ResponseWriter, r *http.Request) error {
	var pass map[string]string
	b, _ := ioutil.ReadAll(r.Body)
	err := json.Unmarshal(b, &pass)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("Invalid JSON")
	}
	password, ok := pass["password"]
	if !ok {
		w.WriteHeader(http.StatusBadRequest)
		return errors.New("You must provide a password to login")
	}
	u := User{Email: r.URL.Query().Get(":email")}
	err = u.Get()
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return errors.New("User not found")
	}
	if u.Login(password) {
		t, _ := NewToken(&u)
		t.Create()
		fmt.Fprintf(w, `{"token":"%x"}`, t.Token)
		return nil
	}
	w.WriteHeader(http.StatusUnauthorized)
	return errors.New("Authentication failed, wrong password")
}
