package netrc

import (
	"fmt"
	"os"
	"os/user"
	"path"
	"testing"
)

func TestFilePath(t *testing.T) {
	usr, _ := user.Current()
	tests := []struct {
		in  string
		out string
	}{
		{"", usr.HomeDir + "/.netrc"},
		{".netlc", usr.HomeDir + "/.netlc"},
		{"/.netrc", "/.netrc"},
	}
	for _, test := range tests {
		res := FilePath(test.in)
		if res != test.out {
			t.Errorf("Response for (%s) got (%s) expected (%s)", test.in, res, test.out)
		}
	}
}

func TestFileExists(t *testing.T) {
	filename := "foo.txt"
	os.Create(filename)
	if !FileExists(filename) {
		t.Error("File should exist")
	}

	os.Remove(filename)
	if FileExists(filename) {
		t.Error("File should not exist")
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
		file, _ := os.Create(test.name)
		file.Chmod(os.FileMode(test.per))

		res := ValidPermissions(test.name)
		if fmt.Sprintf("%o", test.per) == "600" {
			if !res {
				t.Errorf("Got invalid permissions for name (%s) per (%o)", test.name, test.per)
			}
		} else if res {
			t.Errorf("Expected error for name (%s)", test.name)
		}

		os.Remove(test.name)
	}
}

func TestRead(t *testing.T) {
	tests := [][]byte{
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

		n.Read()
		if len(n.Creds) != 1 {
			t.Errorf("Expected (%d) creds got (%d)", 1, len(n.Creds))
		}

		for key, val := range n.Creds {
			if key != "foobar" {
				t.Errorf("Expected key (%s) got (%s)", "foobar", key)
			}

			if val.User != "bazqix" {
				t.Errorf("Expected value (%s) got (%s)", "bazqix", val.User)
			}

			if val.Pass != "sekret" {
				t.Errorf("Expected value (%s) got (%s)", "sekret", val.Pass)
			}
		}

		os.Remove(filepath)
	}
}

func TestGetValue(te *testing.T) {
	tests := []struct {
		key string
		u   string
		p   string
		def bool
	}{
		{"qux", "du", "dp", true},      // Default usage
		{"foo", "foou", "foop", false}, // Specific usage
		{"bar", "baru", "barp", true},  // Specific default fallback
		{"baz", "", "", false},         // Specific no key
	}

	wd, _ := os.Getwd()
	file, _ := os.Create("foo")
	content := []byte("machine foo username foou password foop default username du password dp machine bar username baru password barp")
	file.Write(content)
	filepath := path.Join(wd, "foo")
	n := Open(filepath)

	for _, t := range tests {
		ru, rp, _ := n.GetValue(t.key, t.def)
		if ru != t.u || rp != t.p {
			te.Errorf("Expected u: (%s) p: (%s) got u: (%s) p: (%s)", t.u, t.p, ru, rp)
		}
	}

	os.Remove(filepath)
}

func TestCreateCred(t *testing.T) {
	tests := []struct {
		u    string
		p    string
		pass bool
	}{
		{"foo", "bar", true},
		{"", "baz", false},
		{"qux", "", false},
	}
	for _, test := range tests {
		res, err := CreateCred(test.u, test.p)
		if test.pass {
			if err != nil {
				t.Errorf("Got error (%s) for user: (%s) pass: (%s)", err, test.u, test.p)
			}

			if res.User != test.u || res.Pass != test.p {
				t.Errorf("Expected u: (%s) p: (%s) got u: (%s) p: (%s)", test.u, test.p, res.User, res.Pass)
			}
		} else {
			if err == nil {
				t.Errorf("Expected error u: (%s) p: (%s)", test.u, test.p)
			}

			if res != nil {
				t.Errorf("Expected nil res got: (%+v)", res)
			}
		}
	}
}
