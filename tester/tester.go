package tester

import (
    "encoding/json"
    "crypto/md5"
    "io"
    "fmt"
    "log"
)

func getNewID (str string) string{
    var ret string
    md5sum := md5.New()
    io.WriteString(md5sum, str)
    ret = fmt.Sprintf("%x",md5sum.Sum(nil))
    return ret
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
    b, err := json.MarshalIndent (jt, "", "  ")
    if err != nil {
        return "", err
    }
    return string(b), err
}

func (jt *JSONTest) SetPassedTest (n int, passed bool, status string) {
    jt.PassedTest[n] = passed
    jt.Status[n] = fmt.Sprintf ("Test %d: %s", n, status)
}

func NewTester (in,out []string) *Tester {
    if len(in) == len(out) {
        return &Tester {
            Inputs: in,
            Outputs: out,
        }
    } else {
        log.Fatal ("The test aren't equals")
        return nil
    }
}

func nonEqualIO (in,out string) string {
    return "[Error] Founded " + out + " Expected " + in
}

func (t *Tester) Test (name,source string, in, out []string) *JSONTest {
    var data Data
    compiler := Compiler (source)
    output := make (chan Data)
    jsonTest := NewJSONTest (name, source, len (in))
    go compiler.Run (output,"", 0)
    if data = <-output; data.Err != nil {
        jsonTest.CouldCompile = false
        return jsonTest
    }

    for i := 0; i < len (out);i++ {
        go NewProgram(EXECUTABLE).Run (output,in[i], i)
    }
    
    i := 1
    for i < len(in)+1 {
        data = <-output
        var status string
        passed := (data.Input == data.Output)
        
        if data.Err != nil {
            status = "[Error] -> " + data.Err.Error()
        } else {
            if passed {
                status = "OK"
            } else {
                status = nonEqualIO(data.Input,data.Output)
            }
        }
        jsonTest.SetPassedTest (data.Number, passed, status)
        i++
    }

    return jsonTest
}
