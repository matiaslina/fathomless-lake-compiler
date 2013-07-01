package tester

import (
    "os/exec"
)

type Program struct {
    Label  string
    Program     *exec.Cmd
}

type Data struct {
    Err     error
    Input   string
    Output  string
}

type JSONTest struct {
    ID              string
    Program         string
    Source          string
    CouldCompile    bool
    Count           int
    PassedTest      []bool
    Status          []string
}

type Tester struct {
    Inputs  []string
    Outputs []string
}
