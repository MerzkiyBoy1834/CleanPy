package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

func initProject(projectPath string) {
	if projectPath != "." {
		if err := os.MkdirAll(projectPath, 0755); err != nil {
			fmt.Printf("Error creating project directory: %v\n", err)
			os.Exit(1)
		}
		fmt.Printf("Created directory: %s\n", projectPath)
	}

	oldDir, _ := os.Getwd()
	if projectPath != "." {
		if err := os.Chdir(projectPath); err != nil {
			fmt.Printf("Error changing directory: %v\n", err)
			os.Exit(1)
		}
		defer os.Chdir(oldDir)
	}

	srcDir := "src"
	if err := os.MkdirAll(srcDir, 0755); err != nil {
		fmt.Printf("Error creating src directory: %v\n", err)
		os.Exit(1)
	}

	initContent := `"""Source package for the application."""
# This file makes Python treat the src directory as a package
`
	if err := os.WriteFile(filepath.Join(srcDir, "__init__.py"), []byte(initContent), 0644); err != nil {
		fmt.Printf("Error creating __init__.py: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Created: src/__init__.py")

	mainContent := `"""Main entry point for the application."""
import sys


def main():
    """Main function."""
    print("Hello from CleanPy!")
    print(f"Python version: {sys.version}")
    return 0


if __name__ == "__main__":
    sys.exit(main())
`

	if err := os.WriteFile(filepath.Join(srcDir, "main.py"), []byte(mainContent), 0644); err != nil {
		fmt.Printf("Error creating main.py: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Created: src/main.py")

	requirements := `# Add your dependencies here
# Example:
# requests>=2.28.0
# numpy>=1.24.0
`
	if err := os.WriteFile("requirements.txt", []byte(requirements), 0644); err != nil {
		fmt.Printf("Error creating requirements.txt: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Created: requirements.txt")

	gitignore := `# Python
__pycache__/
*.py[cod]
*$py.class
*.so
.Python
.env
.venv
venv/
ENV/
env/
*.egg-info/
dist/
build/
*.egg

# CleanPy
.build/

# IDE
.vscode/
.idea/
*.swp
*.swo
`
	if err := os.WriteFile(".gitignore", []byte(gitignore), 0644); err != nil {
		fmt.Printf("Error creating .gitignore: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Created: .gitignore")

	if err := exec.Command("git", "init").Run(); err != nil {
		fmt.Printf("Warning: git init failed: %v\n", err)
	} else {
		fmt.Println("Initialized: git repository")
	}

	fmt.Println("\nProject initialized successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. cd " + projectPath)
	fmt.Println("  2. clean-py check  # Check code quality")
	fmt.Println("  3. clean-py run    # Run the project")
}

func projectRoot() (string, error) {
	return os.Getwd()
}

func pythonEnv() ([]string, error) {
	root, err := projectRoot()
	if err != nil {
		return nil, err
	}

	srcPath := filepath.Join(root, "src")
	existing := os.Getenv("PYTHONPATH")
	sep := string(os.PathListSeparator)
	pythonPath := srcPath
	if existing != "" {
		pythonPath = srcPath + sep + existing
	}

	return append(os.Environ(), "PYTHONPATH="+pythonPath), nil
}

func configurePythonCmd(cmd *exec.Cmd) error {
	root, err := projectRoot()
	if err != nil {
		return err
	}

	env, err := pythonEnv()
	if err != nil {
		return err
	}

	cmd.Dir = root
	cmd.Env = env
	return nil
}

func runProject() {
	if _, err := os.Stat("src/main.py"); os.IsNotExist(err) {
		fmt.Println("Error: src/main.py not found. Run 'cleanpy init' first")
		os.Exit(1)
	}

	fmt.Println("\nRunning project...\n")
	cmd := exec.Command("python3", filepath.Join("src", "main.py"))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := configurePythonCmd(cmd); err != nil {
		fmt.Printf("Error preparing Python environment: %v\n", err)
		os.Exit(1)
	}

	if err := cmd.Run(); err != nil {
		fmt.Printf("\nRuntime error: %v\n", err)
		os.Exit(1)
	}
}

func cleanProject() {
	if err := os.RemoveAll(buildDir); err != nil {
		fmt.Printf("Error cleaning: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Build artifacts cleaned!")
}
