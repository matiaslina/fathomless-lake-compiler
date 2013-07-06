package main

import (
    "fmt"
    "html/template"
    "io/ioutil" 
    "log"
    "net/http"
    "os"
    "os/signal"
    "sheltered-inlet/tester"
    "sheltered-inlet/firebase"
    "syscall"
)

var (
    listenAddr  = ":" + os.Getenv("PORT") // Server address
    pwd, _      = os.Getwd()
    RootTemp    = template.Must (template.ParseFiles (pwd + "/views/index.html"))
    
    Firebase    = firebase.New ("https://fathomless-lake.firebaseio.com/")
    // This should be changed for something more cute :3
    Inputs      = []string{"1","2","3"}
    Outputs     = []string{"2","3","4"}
    MyTester    = tester.NewTester(Inputs, Outputs)
)

const (
    FIREBASE_COMPILER_PATH = "/compiler/"
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
    // Upload the file.
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
    go getJsonTest(ch, "prueba", tester.FOLDER + handler.Filename)
    
    log.Println("Running app")
    
    jt := <- ch
    // Store in firebase.
    Firebase.BuildURL (FIREBASE_COMPILER_PATH + jt.ID)
    Firebase.Set (FIREBASE_COMPILER_PATH + jt.ID, jt)
    log.Println ("Show json @ " + FIREBASE_COMPILER_PATH + jt.ID)
    
    // And redirect far away from here.
    http.Redirect(w, req, "http://localhost:3000", http.StatusFound)
}


// API handler, needed to fetch data from fathomless-lake
func APIHandler (w http.ResponseWriter, req *http.Request) {
    w.Header().Set("Content-Type", "application/json")
    // Get the data from firebase
    data, err := Firebase.Get(FIREBASE_COMPILER_PATH)
    if err != nil {
        fmt.Fprintf (w,"Ooops, something happen \n", err.Error())
        return
    }
    fmt.Fprintf (w, string(data))
}

func SignalHandlers (signals chan os.Signal) {
    go func () {
        sig := <- signals
        switch sig {
            case syscall.SIGINT:
                fmt.Printf("\nClosing the server...\nBye bye!\n")
                os.Exit(0)
        }
    }()
}

func main () {
    log.Println("Server running at port " + listenAddr)
    signals := make (chan os.Signal, 1)

    signal.Notify(signals, syscall.SIGINT, syscall.SIGUSR1)
    SignalHandlers (signals)
    err := http.ListenAndServe (listenAddr, nil)
    if err != nil {
        panic ("Listen and serve error: " + err.Error())
    }

}

