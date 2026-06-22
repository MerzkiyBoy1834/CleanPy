package main

import (
	"fmt"
	"os"
)

const (
	buildDir = ".build"
	version  = "1.0"
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
	case "build":
		fmt.Println("Running checks before build...")
		if !checkProjectSilent() {
			fmt.Println("\nBuild aborted due to critical errors!")
			os.Exit(1)
		}
		buildProject()
	case "run":
		fmt.Println("Running checks before build...")
		if !checkProjectSilent() {
			fmt.Println("\nBuild aborted due to critical errors!")
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
  cleanpy init              Initialize project in current directory
  cleanpy new <name>        Create new project with given name
  cleanpy build             Check and compile Python files to .pyc
  cleanpy run               Check, build and run the project
  cleanpy check             Run pylint on the project
  cleanpy clean             Remove build artifacts
  cleanpy version           Show version information`)
}
