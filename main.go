package main

import (
    "bufio"
    "errors"
    "fmt"
    "os"
    "os/exec"
    "strings"
    "os/user"
)


func main() {
    var username string
    var prefixPath string

    reader := bufio.NewReader(os.Stdin)
    user, err := user.Current()
    if err != nil {
        username = "admin"
    } else {
        username = user.Username
    }
    userList := strings.Split(username, "\\")
    username = userList[len(userList)-1]

    for {
        pwd, errPwd := os.Getwd()
        if errPwd != nil {
            prefixPath = "~"
        } else {
            prefixPath = pwd
        }
        
        fmt.Printf(blueItalicPattern + redPattern + yellowItalicPattern + redPattern + magentaPattern ,
                   username, "@", prefixPath, "|", "⇒ ")
        // Read the keyboad input.
        input, err := reader.ReadString('\n')
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
        }

        // Handle the execution of the input.
        if irr := execInput(input); irr != nil {
            fmt.Fprintln(os.Stderr, irr)
        }
    }
}

// ErrNoPath is returned when 'cd' was called without a second argument.
var ErrNoPath = errors.New("path required")

// ErrMuchArgument is returned when have too much arguments
var ErrMuchArguments = errors.New("too many arguments")

// implementation the ANSI escape code for clearing the screen
func clearScreen() {
    fmt.Print("\033[2J")  // Clear entire screen
    fmt.Print("\033[H")   // Move cursor to top-left corner
}

// "source" function
func executeSourceFile(filename string) error {
    f, err := os.Open(filename)
    if err != nil {
        return err
    }
    defer f.Close()

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        if err := execInput(line); err != nil {
            return err
        }
    }

    return scanner.Err()
}


func execInput(input string) error {
    // Remove the newline character.
    input = strings.TrimSuffix(input, "\n")
	input = strings.TrimRight(input, "\r")

    // Split the input separate the command and the arguments.
    args := strings.Split(input, " ")

    // Check for built-in commands.
    switch args[0] {
    case "cd":
        // 'cd' to home with empty path not yet supported.
        if len(args) < 2 {
            return ErrNoPath
        }
        if len(args) > 2 {
            return ErrMuchArguments
        }
        // Change the directory and return the error.
        return os.Chdir(args[1])
    
    case "pwd":
        pwd, err := os.Getwd()
        if err != nil {
            return err
        }
        fmt.Println(pwd)
        return nil
    
    case "ls":
        if len(args) > 2 {
            return ErrMuchArguments
        }
        // Pwd := getPwd()
        ListDir, err := os.ReadDir(".")
        if err != nil {
            return err
        }
        for _, file := range ListDir {
            fmt.Println(file.Name())
        }
        return nil

    case "clear", "cls":
        clearScreen()
        return nil
    
    case "source":
        if len(args) != 2 {
            return errors.New("Usage: source <filename>")
        }
        return executeSourceFile(args[1])
    

    case "exit":
        os.Exit(0)
    }

    // Prepare the command to execute.
    cmd := exec.Command(args[0], args[1:]...)

    // Set the correct output device.
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout

    // Execute the command and return the error.
    return cmd.Run()
}
