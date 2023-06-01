package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Stu struct {
	Id   int
	Name string
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/studentInfo", getStudent).Methods("GET")
	r.HandleFunc("/addStudent", setStudent).Methods("POST")
	http.ListenAndServe(":8080", r)
}

func getStudent(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "root:875254Broj#@tcp(127.0.0.1:3406)/gorest")
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	res, err := db.Query("SELECT * FROM student")
	if err != nil {
		http.Error(w, "Failed to query data from the database", http.StatusInternalServerError)
		return
	}
	defer res.Close()

	for res.Next() {
		var stu Stu
		err := res.Scan(&stu.Id, &stu.Name)
		if err != nil {
			http.Error(w, "Failed to scan row from the database", http.StatusInternalServerError)
			return
		}

		str := "My name is " + stu.Name + " My ID is: " + strconv.Itoa(stu.Id)
		fmt.Fprintln(w, str)
	}
}

func setStudent(w http.ResponseWriter, r *http.Request) {
	data, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusBadRequest)
		return
	}

	var stu Stu
	err = json.Unmarshal(data, &stu)
	if err != nil {
		http.Error(w, "Failed to decode JSON data", http.StatusBadRequest)
		return
	}

	fmt.Fprintln(w, stu.Name)

	db, err := sql.Open("mysql", "root:875254Broj#@tcp(127.0.0.1:3406)/gorest")
	if err != nil {
		http.Error(w, "Failed to connect to the database", http.StatusInternalServerError)
		return
	}
	defer db.Close()

	_, err = db.Exec("INSERT INTO student(name) VALUES(?)", stu.Name)
	if err != nil {
		http.Error(w, "Failed to insert data into the database", http.StatusInternalServerError)
		return
	}

	fmt.Fprintln(w, stu.Name+" has been added to the table")
}
