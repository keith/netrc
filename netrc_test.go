package netrc

import (
    "testing"
    "os"
    "log"
    "fmt"
    "path"
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

func TestRead(t *testing.T) {
    tests := [][]byte {
        []byte("machine foobar username bazqix password sekret"),
    }
    
    for _, test := range tests {
        wd, _ := os.Getwd()
        file, _ := os.Create("foo")
        file.Write(test)
        filepath := path.Join(wd, "foo")
        n := Open(filepath)

        if n.Path != filepath {
            t.Errorf("Expected path to be (%s) got (%s)", filepath, n.Path)
        }

        os.Remove(filepath)
    }
}

