package main

import (
    "bytes"
    "fmt"
    "os/exec"
    "strings"
)

type Program struct {
    Label  string
    Program     *exec.Cmd
}

const (
    EXECUTABLE = "temporal"
)

func Compiler (source string) *Program {

    return &Program {
        Label: "Compiler",
        Program: exec.Command("gcc", "-Wall","-Werror",
                              "-Wextra","-o " + EXECUTABLE, source),
    }
}

func NewProgram(executable string) *Program {
    return &Program {
        Label: executable,
        Program: nil,
    }
}

func (c *Program) Run(channelOut chan string, input string) {
    var out bytes.Buffer
    c.Program.Stdout = &out
    if input != "" {
        c.Program.Stdin = strings.NewReader (input)
    }

    err := c.Program.Run()
    if err != nil {
        channelOut <- err.Error()
        return
    }
    channelOut <- out.String()
}

func main () {
    var out bytes.Buffer
    ch := make (chan string)
    var output string
    c := Compiler("chau")
    c.Program.Stdout = &out
    go c.Run(ch, "")
    output =  <- ch
    fmt.Printf ("[" + c.Label + "] " + output + "\n")
}
