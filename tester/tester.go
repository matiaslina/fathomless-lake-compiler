package tester

import (
    "encoding/json"
    "crypto/md5"
    "io"
    "fmt"
    "log"
)

// Structs with the Array of the I/O
type Tester struct {
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
    
    // Inputs for the tests passed from stdin
    Inputs          []string
    // Outputs of the program
    Outputs         []string
}

// Create a unique hash to store the program with an ID
func getNewID (str string) string{
    var ret string
    md5sum := md5.New()
    io.WriteString(md5sum, str)
    ret = fmt.Sprintf("%x",md5sum.Sum([]byte(str)))
    return ret
}

// Creates a new Tester class. This is the main struct in the program.
func NewTester (programName, source string, 
                in,out []string) *Tester {
    if len(in) == len(out) {
        return &Tester {
            ID: getNewID (programName),
            Program: programName,
            Source: source,
            CouldCompile: false,
            Count: len(in),
            PassedTest: make ([]bool, len(in)),
            Status: make ([]string, len(in)),
            Inputs: in,
            Outputs: out,
        }
    } else {
        log.Fatal ("The test aren't equals")
        return nil
    }
}

// Jsonify converts the struct into an json string well formatted
func (t *Tester) Jsonify () (string, error) {
    b, err := json.MarshalIndent (t, "", "  ")
    if err != nil {
        return "", err
    }
    return string(b), err
}

// Set a new test with the value if the test pass or not.
func (t *Tester) SetPassedTest (n int, passed bool, status string) {
    if n < len(t.PassedTest) {
        t.PassedTest[n] = passed
        t.Status[n] = fmt.Sprintf ("Test %d: %s", n, status)
    } else {
        log.Fatal ("N:",n,"length:", len(t.PassedTest))
    }
}

func nonEqualIO (in,out string) string {
    return "[Error] Founded " + out + " Expected " + in
}

// Main function in the program. Runs a single gorutine for every test.
func (t *Tester) RunTest () *Tester {
    var data Data

    // Compile the app.
    compiler := Compiler (t.Source)
    output := make (chan Data)
    go compiler.Run (output,"", 0)
    if data = <-output; data.Err != nil {
        t.CouldCompile = false
        log.Println("[Error] " + data.Err.Error())
        return t
    }
    t.CouldCompile = true

    // Runs every test.
    for i := 1; i < t.Count;i++ {
        log.Println ("Running app, test", i)
        go NewProgram(EXECUTABLE).Run (output,t.Inputs[i], i)
    }
    
    i := 1
    for i < t.Count {
        log.Println ("Fetching test n", i)
        data = <-output
        var status string
        log.Println("data.Output = "+ data.Output)
        passed := (t.Outputs[data.Number] == data.Output)
        
        if data.Err != nil {
            status = "[Error] -> " + data.Err.Error()
            log.Println("[Error] " + data.Err.Error())
        } else {
            if passed {
                status = "OK"
            } else {
                status = nonEqualIO(data.Input,data.Output)
            }
        }
        t.SetPassedTest (data.Number, passed, status)
        i++
    }
    return t
}
