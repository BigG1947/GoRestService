package main

import (
	"GoRestService/db"
	"database/sql"
	"encoding/json"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
)

var connection *sql.DB
var cash map[int64]db.User

type responseData map[string]interface{}

func main() {
	var err error
	connection, err = db.ConnectionToDb()
	if err != nil {
		log.Fatal(err)
		return
	}

	cash, err = refreshCash()
	if err != nil {
		log.Fatal(err)
		return
	}

	router := mux.NewRouter()
	router.HandleFunc("/api/users", addUsers).Methods("POST")
	router.HandleFunc("/api/users/{id:[0-9]+}", getUsersInfo).Methods("GET")
	router.HandleFunc("/api/users/{id:[0-9]+}", editUsers).Methods("PUT")
	router.HandleFunc("/api/users/{id:[0-9]+}", deleteUsers).Methods("DELETE")
	router.HandleFunc("/api/test/cash", viewCash).Methods("GET")

	log.Fatal(http.ListenAndServe(":8080", router))
}

func deleteUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	login, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("getUsersInfo: %s\n", err)
		return
	}

	user, ok := cash[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !user.CheckAuthData(login, password) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	err = user.Delete(connection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("%s\n", err)
		return
	}

	cash, err = refreshCash()
	if err != nil {
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}

func editUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	login, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("getUsersInfo: %s\n", err)
		return
	}

	user, ok := cash[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !user.CheckAuthData(login, password) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	var newUserInfo db.User

	err = json.NewDecoder(r.Body).Decode(&newUserInfo)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if !newUserInfo.ValidateData() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	newUserInfo.Id = user.Id
	err = newUserInfo.Update(connection)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("UpdateUser: %s\n", err)
		return
	}

	cash, err = refreshCash()
	if err != nil {
		log.Fatal(err)
		return
	}

	w.WriteHeader(http.StatusOK)
	return
}
func addUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	var user db.User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response := make(responseData)
		response["message"] = "Incorrect user info!"
		json.NewEncoder(w).Encode(response)
		return
	}

	if !user.ValidateData() {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	_, exist, err := db.CheckUserExist(connection, user.Login)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("Add user: %s\n", err)
		return
	}

	if !exist {
		user.Id, err = user.Add(connection)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			log.Printf("Add user: %s\n", err)
			return
		}

		cash, err = refreshCash()
		if err != nil {
			log.Fatal(err)
			return
		}
		w.WriteHeader(http.StatusOK)
		return
	} else {
		w.WriteHeader(http.StatusBadRequest)
		response := make(responseData)
		response["message"] = "This user is already exist!"
		json.NewEncoder(w).Encode(response)
	}
}

func getUsersInfo(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	login, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	vars := mux.Vars(r)
	id, err := strconv.ParseInt(vars["id"], 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		log.Printf("getUsersInfo: %s\n", err)
		return
	}

	user, ok := cash[id]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if !user.CheckAuthData(login, password) {
		w.WriteHeader(http.StatusForbidden)
		return
	}

	data := make(responseData)
	data["full_name"] = user.FullName
	data["about"] = user.About
	data["login"] = user.Login
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
	return
}

func viewCash(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cash)
	return
}
