package main

import (
    "encoding/json"
    "fmt"
    "gopkg.in/night-codes/types.v1"
    "html/template"
    "io/ioutil"
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

/*func init() {
    data = []Data{
        {ID: 1, Firstname: "Кантемир", Secondname: "Задорожный", Phone: "+380"},
        {ID: 2, Firstname: "Анна", Secondname: "Задорожная", Phone: "+380"},
        {ID: 3, Firstname: "Виктор", Secondname: "Кондратюк", Phone: "+380"},
    }
}*/

func deleteData(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodDelete {
        return
    }

    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Println("server.go -> deleteData() -> json.ReadAll(): ", err)
        return
    }

    id := types.Uint64(string(body))
    for k := range data {
        if data[k].ID == id {
            data[k] = data[len(data)-1]
            data = data[:len(data)-1]
            break
        }
    }
}

func main() {
    http.HandleFunc("/login", IndexPage)
    http.HandleFunc("/save", saveUsers)
    http.HandleFunc("/get/data", getData)
    http.HandleFunc("/delete/data", deleteData)

    fmt.Println("Server is listening...")
    if err := http.ListenAndServe(":8001", nil); err != nil {
        fmt.Println("main.go -> main() -> ListenAndServe(): ", err)
    }
}

