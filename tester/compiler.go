// Tester package makes everything posible
package tester

import (
    "os/exec"
    "strings"
)

const (
    EXECUTABLE  = "temporal"
    FOLDER      = "tmp/"
)

// An abstract layer of a program.
type Program struct {
    // A label for the program.
    Label  string
    // This is the command that is executed for the program to run.
    Program     *exec.Cmd
}

// The data struct represents the number of the test and the I/O
type Data struct {
    // Number of the test.
    Number  int
    // This field is setted if the program runs with a problem
    // (Pretty useful)
    Err     error
    // The parameters passed by stdin
    Input   string
    // Output of the program.
    Output  string
}

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
    if input != "" {
        c.Program.Stdin = strings.NewReader (input)
    }

    out, err := c.Program.Output()
    channelOut <- Data {
        Number: n,
        Output: string(out),
        Err:    err,
        Input:  input,
    }
}
