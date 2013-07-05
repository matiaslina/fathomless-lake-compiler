package tester

import (
    "os/exec"
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

// Main class that holds all the data of the test. Should be more
// extended in the future.
type JSONTest struct {
    // Unique ID for the test
    ID              string
    // The name of the program
    Program         string
    // The name of the file with the source code of the program
    Source          string
    // If the program can be compiled or not.
    CouldCompile    bool
    // Number of the tests to be runned
    Count           int
    // Array with the passed and unpassed test. The only purpose for
    // This field is check which test are passed.
    PassedTest      []bool
    // Some info about the test. Here are all the errors, comparsion
    // between the I/O, etc.
    Status          []string
}

// Structs with the Array of the I/O
type Tester struct {
    Inputs  []string
    Outputs []string
}
