package main

import (
    "html/template"
    "os"
    "net/http"
    "log"
    "sheltered-inlet/tester"
    "fmt"
)

var (
    listenAddr = ":" + os.Getenv("PORT") // Server address
    pwd, _      = os.Getwd()
    RootTemp    = template.Must (template.ParseFiles (pwd + "/index.html"))
)

func init () {
    http.HandleFunc ("/", RootHandler)
}

func RootHandler (w http.ResponseWriter, req *http.Request) {
    err := RootTemp.Execute (w, listenAddr)
    if err != nil {
        http.Error (w, err.Error(), http.StatusInternalServerError)
    }
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

