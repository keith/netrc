package netrc

import (
    "testing"
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
        per  int
    }{
        {"foo", 0777},
        {"bar", 0600},
        {"baz", 0606},
    }
    for _, test := range tests {
        res := CheckPermissions(test.name)
        if test.per == 0600 && res.Error() != "" {
            t.Errorf("Got error (%s) for name (%s)", res.Error(), test.name)
        } else if res.Error() == "" {
            t.Errorf("Expected error for name (%s)", test.name)
        }
    }
}

