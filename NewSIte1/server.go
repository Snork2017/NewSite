package main

import (
    "fmt"
    "net/http"
    "html/template"
    "encoding/json"
)


type Data struct {
    Firstname string
    Secondname string
    Phone string
    Id int
}

var data []Data 

func getHtml(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        tmpl := template.Must(template.ParseFiles("create.html"))
        tmpl.Execute(w, nil)
    }  
}
func saveUsers(w http.ResponseWriter, r *http.Request) {
    var request Data
    if r.Method == "POST" {
        err := json.NewDecoder(r.Body).Decode(&request)
        if err != nil {
            fmt.Println("server.go -> json.NewDecoder(): ", err)
            return
        }
        for k, _ := range data {
            if data[k].Firstname == request.Firstname{          
                fmt.Println("The name already exists!\n")    
                w.WriteHeader(http.StatusInternalServerError)
                w.Write([]byte("500 - Something bad happened!"))
                return 
            }
        }
        for k, _ := range data {
            if data[k].Id == request.Id {
                request.Id = request.Id + 1 //Присваивание Id каждому новому пользователю!
            }
        }
        data = append(data, request)
    }
}



func getData(w http.ResponseWriter, r *http.Request) {
    body, err := json.Marshal(data)
    if err != nil {
        fmt.Println("server.go -> getData() -> json.Marshal(): ", err)
        return 
    }

    w.Header().Set("Content-Type", "application/json")
    w.Write(body)
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
    http.HandleFunc("/login", getHtml)
    http.HandleFunc("/save", saveUsers)
    http.HandleFunc("/get/data", getData)
    http.HandleFunc("/delete/data", deleteData)
    fmt.Println("Server is listening...")
    http.ListenAndServe(":1010", nil)
}