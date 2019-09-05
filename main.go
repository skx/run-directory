package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"syscall"
)

// IsExecutable returns true if the given path points to an executable file.
func IsExecutable(path string) bool {
	d, err := os.Stat(path)
	if err == nil {
		m := d.Mode()
		return !m.IsDir() && m&0111 != 0
	}
	return false
}

// RunCommand is a helper to run a command, showing the output
// and aborting if the exit-code was non-zero
func RunCommand(command string) (stdout string, stderr string, exitCode int) {
	var outbuf, errbuf bytes.Buffer
	cmd := exec.Command(command)
	cmd.Stdout = &outbuf
	cmd.Stderr = &errbuf

	err := cmd.Run()
	stdout = outbuf.String()
	stderr = errbuf.String()

	if err != nil {
		// try to get the exit code
		if exitError, ok := err.(*exec.ExitError); ok {
			ws := exitError.Sys().(syscall.WaitStatus)
			exitCode = ws.ExitStatus()
		} else {
			// This will happen (in OSX) if `name` is not available in $PATH,
			// in this situation, exit code could not be get, and stderr will be
			// empty string very likely, so we use the default fail code, and format err
			// to string and set to stderr
			exitCode = 0
			if stderr == "" {
				stderr = err.Error()
			}
		}
	} else {
		// success, exitCode should be 0 if go is ok
		ws := cmd.ProcessState.Sys().(syscall.WaitStatus)
		exitCode = ws.ExitStatus()
	}
	return stdout, stderr, exitCode
}

func main() {

	//
	// The command-line flags we accept.
	//
	verbose := flag.Bool("verbose", false, "Show details of what we're doing")
	flag.Parse()

	//
	// Ensure we have a single argument.
	//
	if len(flag.Args()) < 1 {
		fmt.Printf("Usage: rd <directory>")
		os.Exit(1)
	}

	//
	// Get the directory
	//
	input := flag.Args()[0]

	//
	// Find the files beneath that directory.
	//
	files, err := ioutil.ReadDir(input)
	if err != nil {
		fmt.Printf("Error reading directory: %s\n", err)
		os.Exit(1)
	}

	//
	// For each file.
	//
	for _, f := range files {

		//
		// Get the path
		//
		path := filepath.Join(input, f.Name())

		//
		// Skip dotfiles
		//
		if f.Name()[0] == '.' {
			if *verbose {
				fmt.Printf("Skipping dotfile: %s\n", path)
			}
			continue
		}

		//
		// Is it executable?
		//
		if IsExecutable(path) {

			if *verbose {
				fmt.Printf("Running %s\n", path)
			}

			stdout, stderr, exitCode := RunCommand(path)

			//
			// Show STDOUT
			//
			if len(stdout) > 0 {
				fmt.Print(stdout)
			}

			//
			// Show STDERR
			//
			if len(stderr) > 0 {
				fmt.Print(stderr)
			}

			if exitCode != 0 {
				fmt.Printf("%s returned non-zero exit-code\n", path)
				os.Exit(1)
			}

		} else {
			if *verbose {
				fmt.Printf("Skipping non-executable %s\n", path)
			}
		}
	}

}
