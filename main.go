package main

import (
    "code.google.com/p/go.net/websocket"
    "html/template"
    "os"
    "log"
    "net/http"
)

var (
    listenAddr = ":" + os.Getenv("PORT") // Server address
    pwd, _      = os.Getwd()
    RootTemp    = template.Must (template.ParseFiles (pwd + "/index.html"))
    JSON        = websocket.JSON
    Message     = websocket.Message
    ActiveClients   = make (map[ClientConn] int)
)

type ClientConn struct {
    websocket *websocket.Conn
    clientIP string
}

func init () {
    http.HandleFunc ("/", RootHandler)
    http.Handle ("/sock", websocket.Handler (SockServer))
}

func SockServer (ws *websocket.Conn) {
    var err error
    var clientMessage string

    // Clean server side
    defer func () {
        if err = ws.Close(); err != nil {
            log.Println ("Websocket could not be closed", err.Error())
        }
    }()

    client := ws.Request().RemoteAddr
    log.Println ("Client connected:", client)
    sockCli := ClientConn{ws, client}
    ActiveClients[sockCli] = 0
    log.Println ("Number of clients connected ...",len (ActiveClients))

    for {
        if err = Message.Receive(ws, &clientMessage); err != nil {
            log.Println ("websocket Disconnected waiting", err.Error())
            delete (ActiveClients, sockCli)
            log.Println ("Number of clients still connected ...",len (ActiveClients))
            return
        }

        clientMessage = sockCli.clientIP + " Said: " + clientMessage
        for cs, _ := range ActiveClients {
            if err = Message.Send (cs.websocket, clientMessage); err != nil {
                log.Println ("Could not send message to", cs.clientIP, err.Error())
            }
        }
    }
}

func RootHandler (w http.ResponseWriter, req *http.Request) {
    err := RootTemp.Execute (w, listenAddr)
    if err != nil {
        http.Error (w, err.Error(), http.StatusInternalServerError)
    }
}

func main () {
    err := http.ListenAndServe (listenAddr, nil)
    if err != nil {
        panic ("Listen and serve error: " + err.Error())
    }
}

