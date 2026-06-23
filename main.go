package main

import (
	"fmt"
	"os"
)

const (
	buildDir = ".build"
	version  = "1.1"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "init":
		initProject(".")
	case "new":
		if len(os.Args) < 3 {
			fmt.Println("Usage: cleanpy new <project_name>")
			os.Exit(1)
		}
		initProject(os.Args[2])
	case "run":
		fmt.Println("Running checks before run...")
		if !checkProjectSilent() {
			fmt.Println("\nRun aborted due to critical errors!")
			os.Exit(1)
		}
		runProject()
	case "check":
		checkProject()
	case "clean":
		cleanProject()
	case "version":
		fmt.Printf("CleanPy v%s\n", version)
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println(`CleanPy - A simple Python project manager

Usage:
  clean-py init              Initialize project in current directory
  clean-py new <name>        Create new project with given name
  clean-py run               Check and run the project
  clean-py check             Run pylint on the project
  clean-py clean             Remove build artifacts
  clean-py version           Show version information`)
}
