package tester

import (
    "bytes"
    "os/exec"
    "strings"
)

const (
    EXECUTABLE  = "temporal"
    FOLDER      = "tmp/"
)

func Compiler (source string) *Program {

    return &Program {
        Label: "Compiler",
        Program: exec.Command("gcc", "-Wall","-Werror",
                              "-Wextra","-o" + FOLDER + EXECUTABLE, source),
    }
}

func NewProgram(executable string) *Program {
    return &Program {
        Label: executable,
        Program: exec.Command ("./" + FOLDER + EXECUTABLE),
    }
}

func (c *Program) Run(channelOut chan Data, input string, n int) {
    var out bytes.Buffer
    c.Program.Stdout = &out
    if input != "" {
        c.Program.Stdin = strings.NewReader (input)
    }

    err := c.Program.Run()
    channelOut <- Data {
        Number: n+1,
        Output: out.String(),
        Err:    err,
        Input:  input,
    }
}
