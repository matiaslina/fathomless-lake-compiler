package main

import (
    "bytes"
    "log"
    "fmt"
    "os/exec"
    //"strings"
)

type Compiler struct {
    Executable  string
    Gcc         *exec.Cmd
}

func compile (source string, executable string) *Compiler {
    return &Compiler {
        Executable: executable,
        Gcc: exec.Command("gcc", "-Wall","-Werror","-Wextra","-o " + executable,source),
    }
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
