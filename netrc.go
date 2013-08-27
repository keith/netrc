package netrc

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"path/filepath"
	"strings"
)

// A simple credential struct for a username and password
type Cred struct {
	User string
	Pass string
}

// A struct for the Netrc
//	Path is the netrc filepath
//	Creds is a map of string keys that have pointers to Cred structs
//	The keys relate to the service name from the netrc or 'default'
//	ex: 'machine github.com' Then the key is github.com
type Netrc struct {
	Path  string
	Creds map[string]*Cred
}

// Opens the given path if it exists
// If the path is invalid a fatal error is thrown
// A new Netrc is created and the file is read into the Creds map
// A pointer to the Netrc will always be returned
func Open(path string) *Netrc {
	filePath := FilePath(path)
	if !FileExists(filePath) {
		log.Fatalf("Netrc file does not exist in '%s'", filePath)
	}

	n := new(Netrc)
	n.Path = filePath
	n.Read()

	return n
}

// Creates a new Cred struct with a user and pass
//	u: a username
//	p: a password
//	if either u or p are the empty string an error is returned
//	otherwise a pointer to the struct is returned
func CreateCred(u string, p string) (*Cred, error) {
	if u == "" || p == "" {
		return nil, errors.New("Invalid arguments")
	}
	cred := new(Cred)
	cred.User = u
	cred.Pass = p

	return cred, nil
}

// Creates a file path based on the passed string
//	if the empty string is passed $HOME/.netrc is used
//	if a filename is passed ex: .mynetrc it is assumed to be in your
//	  $HOME directory, ex: $HOME/.mynetrc
//	if a path starting with a / is passed it is used without modification
//
//	this uses path/filepath's Join(...) function so you can also pass
//	  relative paths ex: "../.netrc" would return /Users/.netrc on OS X
func FilePath(path string) string {
	if strings.HasPrefix(path, "/") {
		return path
	}

	totalPath := ""
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
	}

	home := usr.HomeDir
	file := ".netrc"

	if path != "" {
		file = path
	}

	totalPath = filepath.Join(home, file)
	return totalPath
}

// Check to see if the file exists at the given path
func FileExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false
	}

	return true
}

// Check the file permissions at the given path
//	600 is valid, all others are invalid
//	returns false if there is an issue with the path
func ValidPermissions(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	per := fmt.Sprintf("%o", os.FileMode.Perm(info.Mode()))
	if per == "600" {
		return true
	}

	return false
}

// Netrc instance methods

// Read the netrc from the file
//	This must be called from a valid netrc instance that has a Path
func (n *Netrc) Read() {
	path := n.Path
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatal(err)
	}

	cs := make(map[string]*Cred)

	words := strings.Fields(string(b))
	for i := 0; i < len(words); i++ {
		word := strings.TrimSpace(words[i])
		key, user, pass := "", "", ""
		if word == "machine" {
			key = strings.TrimSpace(words[i+1])
			user = strings.TrimSpace(words[i+3])
			pass = strings.TrimSpace(words[i+5])
			i += 5
		} else if word == "default" {
			key = word
			user = strings.TrimSpace(words[i+2])
			pass = strings.TrimSpace(words[i+4])
			i += 4
		}

		cred, err := CreateCred(user, pass)
		if err != nil {
			log.Fatal(err)
		}

		cs[key] = cred
	}

	n.Creds = cs
}

// Get a username and password for a service
//	This must be called from a valid netrc instace that has been Read
// Params
//	key: the string following machine in the netrc that identifies the service
//	def: pass true if you want the default credentials (if they exist)
//		 if there aren't any credentials for the given key
// Return
//	string: the username
//	string: the password
//	error:	an error if one occurs
//			 either because of a no credentials or an invalid key
func (n *Netrc) GetValue(key string, def bool) (string, string, error) {
	if n.Creds == nil {
		return "", "", fmt.Errorf("No credentails loaded from %s", n.Path)
	}

	val := n.Creds[key]
	if val == nil && def {
		val = n.Creds["default"]
	}

	if val == nil {
		return "", "", errors.New("No matching credentails")
	} else {
		return val.User, val.Pass, nil
	}
}
