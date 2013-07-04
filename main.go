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
    listenAddr  = ":" + os.Getenv("PORT") // Server address
    pwd, _      = os.Getwd()
    RootTemp    = template.Must (template.ParseFiles (pwd + "/index.html"))
    Inputs      = []string{"1","2","3"}
    Outputs     = []string{"2","3","4"}
    MyTester    = tester.NewTester(Inputs, Outputs)
    JsonAPI     = ""
)

func init () {
    http.HandleFunc ("/", RootHandler)
    http.HandleFunc ("/submitted", SubmittedHandler)
    http.HandleFunc ("/api", APIHandler)
}

func RootHandler (w http.ResponseWriter, req *http.Request) {
    err := RootTemp.Execute (w, listenAddr)
    if err != nil {
        http.Error (w, err.Error(), http.StatusInternalServerError)
    }
}

func getJsonTest (jtc chan *tester.JSONTest, name, source string) {
    jtc <- MyTester.Test(name, source, Inputs, Outputs)
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
    ch := make (chan *tester.JSONTest)
    go getJsonTest(ch, "prueba", "tmp/" + handler.Filename)
    
    log.Println("Running app")
    jt := <- ch
    http.Redirect(w, req, "http://localhost:3000", http.StatusFound)
    JsonAPI, _ = jt.Jsonify()
}

func APIHandler (w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf (w, JsonAPI)
}

func main () {
    log.Println("Server running at port " + listenAddr)
    err := http.ListenAndServe (listenAddr, nil)
    if err != nil {
        panic ("Listen and serve error: " + err.Error())
    }

}

