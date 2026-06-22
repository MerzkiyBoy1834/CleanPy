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

	fmt.Println("\n✅ Project initialized successfully!")
	fmt.Println("\nNext steps:")
	fmt.Println("  1. cd " + projectPath)
	fmt.Println("  2. cleanpy check  # Check code quality")
	fmt.Println("  3. cleanpy build  # Compile the project")
	fmt.Println("  4. cleanpy run    # Run the project")
}
