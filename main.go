package main

import (
    "fmt"
    "html/template"
    "io/ioutil" 
    "log"
    "net/http"
    "os"
    "os/signal"
    "fathomless-lake-compiler/tester"
    "fathomless-lake-compiler/firebase"
    "syscall"
)

var (
    listenAddr  = ":" + os.Getenv("PORT") // Server address
    pwd, _      = os.Getwd()
    RootTemp    = template.Must (template.ParseFiles (pwd + "/views/index.html"))
    DocTemp     = template.Must (template.ParseFiles (pwd + "/views/docs.html"))
    
    Firebase    = firebase.New ("https://fathomless-lake.firebaseio.com/")
    // This should be changed for something more cute :3
    Inputs      = []string{"","1","2","3"}
    Outputs     = []string{"","2","3","4"}
)

const (
    FIREBASE_COMPILER_PATH = "/compiler/"
)

func init () {
    http.HandleFunc ("/", RootHandler)
    http.HandleFunc ("/submitted", SubmittedHandler)
    http.HandleFunc ("/api", APIHandler)
    http.HandleFunc ("/docs", DocsHandler)
    http.Handle ("/styles/", http.StripPrefix("/styles/", http.FileServer(http.Dir("views/styles"))))
    
    // Just a test handler
    http.HandleFunc ("/test", TestHandler)
}

func TestHandler (w http.ResponseWriter, req *http.Request) {
    err := req.ParseMultipartForm(2048)
    if err != nil {
        return
    }
    log.Println(req.Form)
}

func DocsHandler (w http.ResponseWriter, req *http.Request) {
    err := DocTemp.Execute (w, listenAddr)
    if err != nil {
        http.Error (w, err.Error(), http.StatusInternalServerError)
    }
}

func RootHandler (w http.ResponseWriter, req *http.Request) {
    err := RootTemp.Execute (w, listenAddr)
    if err != nil {
        http.Error (w, err.Error(), http.StatusInternalServerError)
    }
}

func StartTest (tc chan *tester.Tester, name, source string) {
    MyTester := tester.NewTester (name, source, Inputs, Outputs)
    tc <- MyTester.RunTest()
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
    ch := make (chan *tester.Tester)
    go StartTest(ch, "prueba", tester.FOLDER + handler.Filename)
    
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

