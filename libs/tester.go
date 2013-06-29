package main

import (
    "encoding/json"
    "crypto/md5"
    "io"
    "fmt"
)

func getNewID (str string) string{
    var ret string
    md5sum := md5.New()
    io.WriteString(md5sum, str)
    ret = fmt.Sprintf("%x",md5sum.Sum(nil))
    return ret
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

func NewJSONTest (programName,source string, test int) *JSONTest {
    return &JSONTest {
        ID: getNewID (programName),
        Program: programName,
        Source: source,
        CouldCompile: false,
        Count: test,
        PassedTest: make ([]bool, test+1),
        Status: make ([]string, test+1),
    }
}

func (jt *JSONTest) Jsonify () (string, error) {
    b, err := json.Marshal (jt)
    if err != nil {
        return "", err
    }
    return string(b), err
}

func (jt *JSONTest) SetPassedTest (n int, passed bool, status string) {
    jt.PassedTest[n] = passed
    jt.Status[n] = fmt.Sprintf ("Test %d: %s", n, status)
}


/* Don't mind about this, should be eliminated in
 * a short time 
 */
func main () {
    /* Checking the jsonify */
    a := NewJSONTest ("lucky","lucky.c",3)
    a.SetPassedTest(1, true, "OK")
    a.SetPassedTest(2, false, "Incorrect answer")

    b, err := a.Jsonify ()
    if err != nil {
        fmt.Println(getNewID ("Holas"))
    } else {
        fmt.Println(b)
    }
}
