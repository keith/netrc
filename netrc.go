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

