package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

var (
	config struct {
		path string
	}
)

func init() {
	flag.StringVar(&config.path, "p", "/etc/secret-volume", "a secret volume's mount path")
}

func main() {
	flag.Parse()

	// Find the expanded path to the given executable.
	command, err := exec.LookPath(flag.Arg(0))
	if err != nil {
		log.Fatal(err)
	}

	// We want to use the current environment as a starting point and
	// overwrite values included as part of the secret object.
	env := make(map[string]string)
	for _, variable := range os.Environ() {
		components := strings.Split(variable, `=`)
		env[components[0]] = components[1]
	}

	// List the provided directory.
	files, err := ioutil.ReadDir(config.path)
	if err != nil {
		log.Fatal(err)
	}

	for _, file := range files {

		// Verify that the file is a file, not a directory.
		filename := file.Name()
		fileInfo, err := os.Stat(filepath.Join(config.path, filename))
		if err != nil {
			log.Fatal(err)
		}
		if fileInfo.IsDir() {
			continue
		}

		// Read the file.
		value, err := ioutil.ReadFile(filepath.Join(config.path, filename))
		if err != nil {
			log.Fatal(err)
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

	log.Fatal(
		syscall.Exec(command, flag.Args(), commandEnv),
	)
}
