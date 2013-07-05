// Tester package makes everything posible
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

// Creates a compiler for the program, the output it's the same
// for every program since the objective it's not to store the 
// exec files  
func Compiler (source string) *Program {
    return &Program {
        Label: "Compiler",
        Program: exec.Command("gcc", "-Wall","-Werror",
                              "-Wextra","-o" + FOLDER + EXECUTABLE, source),
    }
}

// Creates a new program with the name of the executable
func NewProgram(executable string) *Program {
    return &Program {
        Label: executable,
        Program: exec.Command ("./" + FOLDER + EXECUTABLE),
    }
}

// Runs te program with the input and the number of the test
// that are trying to check.
// The output is passed by a channel with the purpose of using
// gorutines in the main program.
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
