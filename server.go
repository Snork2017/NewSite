package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

type Data struct {
	ID         uint64
	Firstname  string
	Secondname string
	Phone      string
}

var (
	data  []Data
	count uint64
)

func IndexPage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		return
	}

	tmpl := template.Must(template.ParseFiles("create.html"))
	if err := tmpl.Execute(w, nil); err != nil {
		fmt.Println("main.go -> IndexPage() -> Execute(): ", err)
	}
}

func saveUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		return
	}

	var request Data
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		fmt.Println("server.go -> json.NewDecoder(): ", err)
		return
	}

	for k := range data {
		if data[k].Firstname == request.Firstname {
			fmt.Println("The name already exists!\n")
			w.WriteHeader(http.StatusInternalServerError)

			if _, err := w.Write([]byte("500 - Something bad happened!")); err != nil {
				fmt.Println("main.go -> saveUsers() -> Write(): ", err)
			}
			return
		}
	}

	count++
	request.ID = count

	data = append(data, request)
}

func getData(w http.ResponseWriter, r *http.Request) {
	body, err := json.Marshal(data)
	if err != nil {
		fmt.Println("server.go -> getData() -> json.Marshal(): ", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if _, err := w.Write(body); err != nil {
		fmt.Println("main.go -> getData() -> Write(): ", err)
	}
}

func deleteData(w http.ResponseWriter, r *http.Request) {
	//var request Data
	//if r.Method == "DELETE" {
	//for k, _ := range data {
	//if request.Id == data[k].Id {
	//data = append(data[:request.Id], data[request.Id+1:]...)
	//}
	//}
	//}
}

func main() {
	http.HandleFunc("/login", IndexPage)
	http.HandleFunc("/save", saveUsers)
	http.HandleFunc("/get/data", getData)
	http.HandleFunc("/delete/data", deleteData)

	fmt.Println("Server is listening...")
	if err := http.ListenAndServe(":8003", nil); err != nil {
		fmt.Println("main.go -> main() -> ListenAndServe(): ", err)
	}
}
