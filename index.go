package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type User struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Location string `json:"location"`
}

const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "744833"
	dbname   = "test"
)

var conn = connectdb()

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/user", getUsers).Methods("GET")
	r.HandleFunc("/user", createUsers).Methods("POST")
	r.HandleFunc("/user/{id}", getUser).Methods("GET")
	r.HandleFunc("/user/{id}", updateUser).Methods("PUT")
	r.HandleFunc("/user/{id}", deleteUser).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", r))
}
func createUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)

	if err != nil {
		log.Fatal(err)
	}

	sqlStatement := `INSERT INTO "user" (id,name,location) VALUES ($1,$2,$3)`
	_, err = conn.Exec(sqlStatement, user.Id, user.Name, user.Location)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(user)
}
func getUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	rows, err := conn.Query("SELECT * FROM public.user LIMIT 1000")
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	users := []User{}

	for rows.Next() {
		user := User{}
		err = rows.Scan(&user.Id, &user.Name, &user.Location)
		if err != nil {
			log.Fatal(err)
		}
		users = append(users, user)
		fmt.Println(user)
	}
	// err = rows.Err()
	// if err != nil {
	// 	panic(err)
	// }

	json.NewEncoder(w).Encode(users)

}
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sqlStatement := `SELECT * FROM public.user WHERE id=$1;`
	user := User{}
	param := r.URL.Query().Get("id")
	row := conn.QueryRow(sqlStatement, param)
	err := row.Scan(&user.Id, &user.Name, &user.Location)
	switch err {
	case sql.ErrNoRows:
		fmt.Println("No rows were returned!")
	case nil:
		fmt.Println(user)
	default:
		panic(err)
	}
	json.NewEncoder(w).Encode(user)

}
func updateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	user := User{}
	param := r.URL.Query().Get("id")
	err := json.NewDecoder(r.Body).Decode(&user)
	sqlStatement := `UPDATE public.user SET name = $2, location = $3 WHERE id = $1;`
	_, err = conn.Exec(sqlStatement, param, user.Name, user.Location)
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully updated")
	json.NewEncoder(w).Encode(user)

}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	sqlStatement := `DELETE FROM public.user WHERE id = $1;`
	param := r.URL.Query().Get("id")
	result, err := conn.Exec(sqlStatement, param)
	if err != nil {
		panic(err)
	}
	fmt.Println("successfully deleted")
	json.NewEncoder(w).Encode(result)
}
func connectdb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	fmt.Println("Successfully connected!")
	return db
}
