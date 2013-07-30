package netrc

import (
    "testing"
    "os"
    "log"
    "fmt"
)

func TestFilePath(t *testing.T) {
    tests := []struct {
        in string
        out string
    }{
        {"", "/Users/ksmiley/.netrc"},
        {".netlc", "/Users/ksmiley/.netlc"},
        {"/.netrc", "/.netrc"},
    }
    for _, test := range tests {
        res := FilePath(test.in)
        if res != test.out {
            t.Errorf("Response for (%s) got (%s) expected (%s)", test.in, res, test.out)
        }
    }
}

func TestPermissions(t *testing.T) {
    tests := []struct {
        name string
        per  uint32
    }{
        {"foo", 0777},
        {"bar", 0600},
        {"baz", 0606},
    }
    for _, test := range tests {
        file, err := os.Create(test.name)
        if err != nil {
            log.Fatal(err)
        }

        file.Chmod(os.FileMode(test.per))

        res := CheckPermissions(test.name)
        if fmt.Sprintf("%o", test.per) == "600" {
            if res != nil {
                t.Errorf("Got error (%s) for name (%s) per (%o)", res.Error(), test.name, test.per)
            }
        } else if res == nil {
            t.Errorf("Expected error for name (%s)", test.name)
        }

        err = os.Remove(test.name)
        if err != nil {
            log.Fatal(err)
        }
    }
}

