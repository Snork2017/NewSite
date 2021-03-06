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
    ID         uint64 `json:"id"`
    Firstname  string 
    Secondname string
    Phone      string
}

var (
    data  []Data
    count  uint64
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
func editData(w http.ResponseWriter, r *http.Request) {
    //var ww = "kant"
    var request Data
    if r.Method != http.MethodPost {
        return
}
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Println("server.go -> editData() -> json.ReadAll(): ", err)
        return
    }
    fmt.Println(string(body))
    if json.Unmarshal(body, &request);   err != nil {
        fmt.Println("server.go -> editData() -> json.Unmarshal(): ", err)
        return
    }
    fmt.Println("ИМЯ", request)
    for k := range data {
        if data[k].ID == request.ID{
            data[k].Firstname = request.Firstname
            data[k].Secondname = request.Secondname
            data[k].Phone = request.Phone

        }
    } 
    data = append(data)
}

func popUpUser(w http.ResponseWriter, r *http.Request) {

    if r.Method != http.MethodPost {
        return
    }
    body, err := ioutil.ReadAll(r.Body)
    if err != nil {
        fmt.Println("server.go -> popUpUser() -> json.ReadAll(): ", err)
        return
    }
    id := types.Uint64(string(body))
    fmt.Println("UserID:", id)
    for k := range data {
        if data[k].ID == id {
            a, _ := json.Marshal(data)
            w.Write(a)
            fmt.Println(string(a))
        }
    }
}    
func main() {
    http.HandleFunc("/receive", popUpUser)
    http.HandleFunc("/login", IndexPage)
    http.HandleFunc("/save", saveUsers)
    http.HandleFunc("/get/data", getData)
    http.HandleFunc("/delete/data", deleteData)
    http.HandleFunc("/edit/data", editData)
    fmt.Println("Server is listening...")
    if err := http.ListenAndServe(":8155", nil); err != nil {
        fmt.Println("main.go -> main() -> ListenAndServe(): ", err)
    }
}