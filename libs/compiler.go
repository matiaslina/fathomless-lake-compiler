package main

import (
    "bytes"
    "log"
    "fmt"
    "os/exec"
    "strings"
)

type Compiler struct {
    Executable  string
    Gcc         *exec.Cmd
    Program     *exec.Cmd
}

func compile (source string, executable string) *Compiler {
    return &Compiler {
        Executable: executable,
        Gcc: exec.Command("gcc", "-Wall","-Werror","-Wextra","-o " + executable,source),
        Program: nil,
    }
}

func (c *Compiler) RunProgram (channelOut chan string, input string) {
    var out bytes.Buffer
    c.Program = exec.Command (c.Executable)
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
    c := compile ("algo", "chau")
    var out bytes.Buffer
    c.Gcc.Stdout = &out
    err := c.Gcc.Run()
    if err != nil {
        log.Fatal (err)
    }
    fmt.Printf ("[" + c.Executable + "]\n", out.String())
}
