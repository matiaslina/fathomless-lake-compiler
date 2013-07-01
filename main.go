package main

import (
    "html/template"
    "os"
    "net/http"
    "log"
    "sheltered-inlet/tester"
    "fmt"
    "io/ioutil" 
)

var (
    listenAddr = ":" + os.Getenv("PORT") // Server address
    pwd, _      = os.Getwd()
    RootTemp    = template.Must (template.ParseFiles (pwd + "/index.html"))
)

func init () {
    http.HandleFunc ("/", RootHandler)
    http.HandleFunc ("/submitted", SubmittedHandler)
}

func RootHandler (w http.ResponseWriter, req *http.Request) {
    err := RootTemp.Execute (w, listenAddr)
    if err != nil {
        http.Error (w, err.Error(), http.StatusInternalServerError)
    }
}

func SubmittedHandler (w http.ResponseWriter, req *http.Request) {
    file, handler, err := req.FormFile ("file")
    if err != nil {
        fmt.Println("[step 1]Oh.. we have an error with the file :/")
    }
    data, err := ioutil.ReadAll (file)
    if err != nil {
        fmt.Println("[step 2]Oh.. we have an error with the file :/")
    }
    err = ioutil.WriteFile ("tmp/" + handler.Filename, data, 0755)
    if err != nil {
        fmt.Println ("[step 3] Oh.. can't set the file in the disk :/")
    }
    
    http.Redirect(w, req, "http://localhost:3000", http.StatusFound)
}

func main () {
    log.Println("Server running at port " + listenAddr)
    a := tester.NewJSONTest ("hola","chau",1)
    s, _ := a.Jsonify ()
    fmt.Println ("JSON -> ", s)
    err := http.ListenAndServe (listenAddr, nil)
    if err != nil {
        panic ("Listen and serve error: " + err.Error())
    }

}

