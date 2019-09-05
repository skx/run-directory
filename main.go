// run-directory
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

// Arguments set by the command-line arguments, along with our version-string.
var (
	verbose     *bool
	exitOnError *bool
	version     = "unreleased"
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

// RunCommand is a helper to run a command, returning output and the exit-code.
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
			// This will happen (in OSX) if `name` is not
			// available in $PATH, in this situation, exit
			// code could not be get, and stderr will be
			// empty string very likely, so we use the default
			// fail code, and format err to string and set to stderr
			exitCode = 1
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

// RunParts runs all the executables in the given directory.
func RunParts(directory string) {

	//
	// Find the files beneath the named directory.
	//
	files, err := ioutil.ReadDir(directory)
	if err != nil {
		fmt.Printf("error reading directory contents %s - %s\n", directory, err)
		os.Exit(1)
	}

	//
	// For each file we found.
	//
	for _, f := range files {

		//
		// Get the absolute path to the file.
		//
		path := filepath.Join(directory, f.Name())

		//
		// We'll skip any dotfiles.
		//
		if f.Name()[0] == '.' {
			if *verbose {
				fmt.Printf("Skipping dotfile: %s\n", path)
			}
			continue
		}

		//
		// We'll skip any non-executable files.
		//
		if !IsExecutable(path) {
			if *verbose {
				fmt.Printf("Skipping non-executable %s\n", path)
			}
			continue
		}

		//
		// Show what we're doing.
		//
		if *verbose {
			fmt.Printf("Running %s\n", path)
		}

		//
		// Run the command, capturing output and exit-code
		//
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

		//
		// If the exit-code was non-zero then we have to
		// terminate.
		//
		if exitCode != 0 {
			if *verbose {
				fmt.Printf("%s returned non-zero exit-code\n", path)
			}
			if *exitOnError {
				os.Exit(1)
			}
		}
	}
}

// main is our entry-point
func main() {

	//
	// The command-line flags we accept.
	//
	verbose = flag.Bool("verbose", false, "Show details of what we're doing")
	ver := flag.Bool("version", false, "Show our version")
	exitOnError = flag.Bool("exit-on-error", true, "Exit when the first script fails")
	flag.Parse()

	//
	// Show our version?
	//
	if *ver {
		fmt.Printf("run-directory %s\n", version)
		os.Exit(0)
	}

	//
	// Ensure we have at least one argument.
	//
	if len(flag.Args()) < 1 {
		fmt.Printf("Usage: rd <directory1> [directory2] .. [directoryN]\n")
		os.Exit(1)
	}

	//
	// Process each named directory
	//
	for _, entry := range flag.Args() {
		RunParts(entry)
	}

	//
	// Exit with a successfull status-code
	//
	os.Exit(0)
}
