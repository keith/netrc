package netrc

import (
    "os"
    "os/user"
    "fmt"
    "io/ioutil"
    "strings"
    "log"
    "errors"
)

type Cred struct {
    User string
    Pass string
}

type Netrc struct {
    Path string
    Creds []map[string]Cred
}

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

func (n *Netrc) Read() {
    path := n.Path
    b, err := ioutil.ReadFile(path)
    if err != nil {
        log.Fatal(err)
    }

    lines := strings.Fields(string(b))
    for i := 0; i < len(lines); i++ {
        line := strings.TrimSpace(lines[i])
        fmt.Printf("'%s'\n", line)
        // if line == "" {
        //     continue
        // }
    }
}

func FilePath(path string) string {
    totalPath := ""
    usr, err := user.Current()
    if err != nil {
        log.Fatal(err)
    }

    home := usr.HomeDir
    components := []string{home, "/.netrc"}

    if (strings.HasPrefix(path, "/")) {
        return path
    } else if (path != "") {
        components[1] = fmt.Sprintf("/%s", path)
    }

    totalPath = strings.Join(components, "")
    return totalPath
}

func FileExists(path string) bool {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        return false
    }

    return true
}

func CheckPermissions(path string) error {
    info, err := os.Stat(path)
    if err != nil {
        return err
    }

    per := fmt.Sprintf("%o", os.FileMode.Perm(info.Mode()))
    if per == "600" {
        return nil
    }

    return errors.New("Incorrect file permissions")
}


func Read() (string, error) {
    // home := os.Getenv("HOME")
    home := ""
    components := []string{home, "/.netrc"}
    path := strings.Join(components, "")

    fmt.Printf("%s %s\n", home, path)
    b, err := ioutil.ReadFile("/Users/ksmiley/.netrc")
    if err != nil { return "", err }
    lines := strings.Split(string(b), "\n")
    // fmt.Printf(lines)
    return lines[0], nil
    // file, err := os.Open("$HOME/.netrc")
    // if err != nil {
    //     return "Err"
    // }

    // buf := make([]byte, 1024)
    // for {
    //     n, err := file.Read(buf)
    //     if err != nil && err != io.EOF {
    //         return "Err2"
    //     }

    //     if n == 0 {
    //         return "0 read"
    //     }
    // }

    // fmt.Printf(buf)
    // return ""
}

