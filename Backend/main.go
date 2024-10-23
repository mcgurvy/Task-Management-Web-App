package main

import (
	"github.com/gorilla/handlers"
)

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

type Task struct {
	ID    int    `json:"id"`
	Title string `json:"title"`
	Done  bool   `json:"done"`
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", "./tasks.db")
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()
	createTable()

	r := mux.NewRouter()
	r.HandleFunc("/tasks", getTasks).Methods("GET")
	r.HandleFunc("/task", createTask).Methods("POST")
	r.HandleFunc("/task/{id}", updateTask).Methods("PUT")
	r.HandleFunc("/task/{id}", deleteTask).Methods("DELETE")

	http.Handle("/", r)
	log.Println("Server running on port 8080")
	http.ListenAndServe(":8080", nil)
	http.ListenAndServe(":8080", handlers.CORS(handlers.AllowedOrigins([]string{"http://localhost:3000"}))(r))
}

func createTable() {
	query := `
	CREATE TABLE IF NOT EXISTS tasks (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		done BOOLEAN
	);
	`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func getTasks(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT id, title, done FROM tasks")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var tasks []Task
	for rows.Next() {
		var task Task
		err := rows.Scan(&task.ID, &task.Title, &task.Done)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		tasks = append(tasks, task)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(tasks)
}

func createTask(w http.ResponseWriter, r *http.Request) {
	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("INSERT INTO tasks (title, done) VALUES (?, ?)", task.Title, task.Done)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func updateTask(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	var task Task
	err := json.NewDecoder(r.Body).Decode(&task)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, err = db.Exec("UPDATE tasks SET title = ?, done = ? WHERE id = ?", task.Title, task.Done, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func deleteTask(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["id"]

	_, err := db.Exec("DELETE FROM tasks WHERE id = ?", id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
