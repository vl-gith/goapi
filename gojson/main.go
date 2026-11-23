package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone int    `json:"phone"`
}

func userAdd(response http.ResponseWriter, request *http.Request) {
	if request.Method != http.MethodPost {
		http.Error(response, "Метод не поддерживается, используйте метод POST", http.StatusMethodNotAllowed)
		return
	}

	var newUser User

	err := json.NewDecoder(request.Body).Decode(&newUser)

	if err != nil {
		http.Error(response, "Ошибка чтения JSON: "+err.Error(), http.StatusBadRequest)
		return
	}

	query := `INSERT INTO users (name, phone) VALUES ($1, $2)`
	_, err = dbPool.Exec(request.Context(), query, newUser.Name, newUser.Phone)

	if err != nil {
		http.Error(response, "Ошибка сохранения в БД: "+err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Print("Пользователь добавлен!\n")

	response.WriteHeader(http.StatusCreated)
	fmt.Fprint(response, "Пользователь успешно сохранен!\r\n")
}

func userList(response http.ResponseWriter, request *http.Request) {
	query := `SELECT id, name, phone FROM users`
	rows, err := dbPool.Query(request.Context(), query)

	if err != nil {
		http.Error(response, "Ошибка чтения из БД: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	users := make([]User, 0)

	for rows.Next() {
		var user User
		err := rows.Scan(&user.ID, &user.Name, &user.Phone)
		if err != nil {
			http.Error(response, "Ошибка сканирования строки: "+err.Error(), http.StatusInternalServerError)
			return
		}
		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		http.Error(response, "Ошибка при чтении потока строк: "+err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	json.NewEncoder(response).Encode(users)
}

func main() {

	connStr := "postgres://login:password@ip:port/testjson"
	var err error

	dbPool, err = pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Ошибка настройки пула: %v\n", err)
	}
	defer dbPool.Close()

	err = dbPool.Ping(context.Background())
	if err != nil {
		log.Fatalf("Не удалось подключиться к базе данных %v\n", err)
	}
	fmt.Println("Успешно подключились к PostgreSQL")

	http.HandleFunc("/useradd", userAdd)
	http.HandleFunc("/userlist", userList)

	http.ListenAndServe("localhost:8888", nil)
}
