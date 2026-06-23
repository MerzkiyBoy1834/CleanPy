package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func checkProject() {
	if _, err := os.Stat("src"); os.IsNotExist(err) {
		fmt.Println("Error: src directory not found. Run 'cleanpy init' first.")
		os.Exit(1)
	}

	fmt.Println("Running pylint checks...")
	fmt.Println()

	var pyFiles []string
	err := filepath.Walk("src", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".py") {
			pyFiles = append(pyFiles, path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error scanning files: %v\n", err)
		os.Exit(1)
	}

	if len(pyFiles) == 0 {
		fmt.Println("No Python files found in src/")
		os.Exit(1)
	}

	args := append([]string{"-m", "pylint", "--output-format=colorized"}, pyFiles...)
	cmd := exec.Command("python3", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := configurePythonCmd(cmd); err != nil {
		fmt.Printf("Error preparing Python environment: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode := exitErr.ExitCode()
			if exitCode&1 == 1 || exitCode&2 == 2 {
				fmt.Println()
				fmt.Println("Found critical errors. Fix them before running.")
			} else if exitCode > 0 {
				fmt.Printf("\nFound non-critical issues (exit code: %d).\n", exitCode)
			}
			os.Exit(exitCode)
		}
		fmt.Printf("\nError running pylint: %v\n", err)
		os.Exit(1)
	}

	fmt.Println()
	fmt.Println("All checks passed.")
}

func checkProjectSilent() bool {
	if _, err := os.Stat("src"); os.IsNotExist(err) {
		fmt.Println("Error: src directory not found. Run 'cleanpy init' first.")
		return false
	}

	var pyFiles []string
	err := filepath.Walk("src", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(path, ".py") {
			pyFiles = append(pyFiles, path)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Error scanning files: %v\n", err)
		return false
	}

	if len(pyFiles) == 0 {
		fmt.Println("No Python files found in src/")
		return false
	}

	args := append([]string{"-m", "pylint", "--output-format=colorized"}, pyFiles...)
	cmd := exec.Command("python3", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := configurePythonCmd(cmd); err != nil {
		fmt.Printf("Error preparing Python environment: %v\n", err)
		return false
	}

	if err := cmd.Run(); err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode := exitErr.ExitCode()
			if exitCode&1 == 1 || exitCode&2 == 2 {
				fmt.Println()
				fmt.Println("Found critical errors. Run aborted.")
				return false
			}
			if exitCode > 0 {
				fmt.Printf("\nFound non-critical issues (exit code: %d). Continuing...\n", exitCode)
				return true
			}
		}
		return false
	}

	fmt.Println()
	fmt.Println("All checks passed.")
	return true
}
