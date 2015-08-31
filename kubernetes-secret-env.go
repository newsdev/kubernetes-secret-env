package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

// fatal prints an error's string value to stderr and exits with a non-zero
// status.
func fatal(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}

var (
	NoCommandError = errors.New("a command must be provided")
	path           string
)

func init() {
	flag.StringVar(&path, "p", "/etc/secret-volume", "a secret volume's mount path")
}

func main() {
	flag.Parse()

	// Find the expanded path to the given executable.
	command, err := exec.LookPath(flag.Arg(0))
	if err != nil {
		fatal(err)
	}

	// We want to use the current environment as a starting point and
	// overwrite values included as part of the secret object.
	env := make(map[string]string)
	for _, variable := range os.Environ() {
		components := strings.Split(variable, `=`)
		env[components[0]] = components[1]
	}

	// List the provided directory.
	files, err := ioutil.ReadDir(path)
	if err != nil {
		fatal(err)
	}

	for _, file := range files {

		// Read the file.
		filename := files.Name()
		value, err := ioutil.ReadFile(filepath.Join(path, filename))
		if err != nil {
			fatal(err)
		}

		// Set the key/value pair in the environment, overwriting existing values.
		varname := strings.Replace(strings.ToUpper(filename), "-", "_", -1)
		env[varname] = string(value)
	}

	// Exec expects the environment to be specified as a slice rather than a
	// map. Flatten it by joining keys and values with `=`.
	commandEnv := make([]string, 0, len(env))
	for key, value := range env {
		commandEnv = append(commandEnv, fmt.Sprintf("%s=%s", key, value))
	}

	if err := syscall.Exec(command, flag.Args(), commandEnv); err != nil {
		fatal(err)
	}
}
