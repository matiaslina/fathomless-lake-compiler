package tester

import (
    "bytes"
    "os/exec"
    "strings"
)

const (
    EXECUTABLE = "temporal"
)

func Compiler (source string) *Program {

    return &Program {
        Label: "Compiler",
        Program: exec.Command("gcc", "-Wall","-Werror",
                              "-Wextra","-o" + EXECUTABLE, source),
    }
}

func NewProgram(executable string) *Program {
    return &Program {
        Label: executable,
        Program: nil,
    }
}

func (c *Program) Run(channelOut chan Data, input string) {
    var out bytes.Buffer
    c.Program.Stdout = &out
    if input != "" {
        c.Program.Stdin = strings.NewReader (input)
    }

    err := c.Program.Run()
    channelOut <- Data {
        Output: out.String(),
        Err:    err,
        Input:  input,
    }
}
