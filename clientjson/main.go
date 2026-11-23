package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"strconv"

	"encoding/json"
	"net/http"
)

type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Phone int    `json:"phone"`
}

func showUsers() {
	response, err := http.Get("http://localhost:8888/userlist")
	if err != nil {
		fmt.Println("Ошибка получения ресурса", err)
		return
	}
	defer response.Body.Close()

	var users []User
	if err := json.NewDecoder(response.Body).Decode(&users); err != nil {
		fmt.Println("Ошибка распаковки JSON: ", err)
		return
	}

	fmt.Println("Список зарегистрированных пользователей: ")
	if len(users) == 0 {
		fmt.Println("Пока никого нет.")
		return
	}

	for _, user := range users {
		fmt.Printf("ID: %d | Имя: %s | Телефон: %d\n", user.ID, user.Name, user.Phone)
	}

	fmt.Println("Ответ от сервера: ", response.Status)
}

func addUser(scanner *bufio.Scanner) {
	fmt.Println("Добавьте данные пользователя!")

	fmt.Print("Введите имя: ")
	scanner.Scan()
	name := scanner.Text()

	fmt.Print("Введите телефон: ")
	scanner.Scan()
	phoneStr := scanner.Text()

	phone, err := strconv.Atoi(phoneStr)
	if err != nil {
		fmt.Println("Ошибка: телефон должен состоять только из цифр!")
		return
	}

	fmt.Printf("Вы добавили: %s, number: %d\n", name, phone)

	user := User{
		Name:  name,
		Phone: phone,
	}

	jsonStr, err := json.Marshal(user)

	if err != nil {
		fmt.Println("Ошибка упаковки в JSON: ", err)
		return
	}

	response, err := http.Post("http://localhost:8888/useradd", "application/json", bytes.NewBuffer(jsonStr))

	if err != nil {
		fmt.Println("Ошибка отправки на сервер: ", err)
		return
	}

	defer response.Body.Close()

	bodyBytes, _ := io.ReadAll(response.Body)
	fmt.Println("Ответ от сервера: ", response.Status)
	fmt.Println("Сообщение: ", string(bodyBytes))
}

func main() {
	fmt.Print("Введите \"1\", чтобы показать пользователей, \"2\" -- зарегистрировать нового: ")
	scanner := bufio.NewScanner(os.Stdin)

	scanner.Scan()
	userChoice := scanner.Text()

	switch userChoice {
	case "1":
		showUsers()
	case "2":
		addUser(scanner)
	default:
		fmt.Println("Внимание! Вы выбрали значение, для которого нет доступных функций, выбрано значение по умолчанию -- 1")
		showUsers()
	}
}
