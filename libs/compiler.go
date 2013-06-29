package main

import (
    "bytes"
    "log"
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
    c.Program = exec.Command (EXECUTABLE)
    c.Program.Stdout = &out
    c.Program.Stdin = strings.NewReader (input)

    err := c.Program.Run()
    if err != nil {
        channelOut <- err.Error()
        return
    }
    channelOut <- out.String()
}

func main () {
    c := Compiler("chau")
    var out bytes.Buffer
    c.Program.Stdout = &out
    err := c.Program.Run()
    if err != nil {
        log.Fatal (err)
    }
    fmt.Printf ("[" + c.Label + "]\n", out.String())
}
