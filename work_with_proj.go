package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func buildProject() {
	if _, err := os.Stat("src/main.py"); os.IsNotExist(err) {
		fmt.Println("Error: src/main.py not found. Run 'cleanpy init' first")
		os.Exit(1)
	}

	if err := os.MkdirAll(buildDir, 0755); err != nil {
		fmt.Printf("Error creating build directory: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Building project...")

	err := filepath.Walk("src", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(path, ".py") {
			relPath, _ := filepath.Rel("src", path)

			destDir := filepath.Join(buildDir, filepath.Dir(relPath))
			if err := os.MkdirAll(destDir, 0755); err != nil {
				return err
			}

			cmd := exec.Command("python3", "-m", "compileall", "-q", path)
			if output, err := cmd.CombinedOutput(); err != nil {
				fmt.Printf("Error compiling %s:\n%s\n", path, string(output))
				return err
			}

			dir := filepath.Dir(path)
			pycacheDir := filepath.Join(dir, "__pycache__")

			matches, err := filepath.Glob(filepath.Join(pycacheDir, "*.pyc"))
			if err != nil || len(matches) == 0 {
				return fmt.Errorf("compiled file not found for %s", path)
			}

			pycFile := matches[0]

			baseName := strings.TrimSuffix(filepath.Base(path), ".py")
			destFile := filepath.Join(buildDir, filepath.Dir(relPath), baseName+".pyc")

			if err := copyFile(pycFile, destFile); err != nil {
				return err
			}

			hash, _ := hashFile(path)
			os.WriteFile(destFile+".hash", []byte(hash), 0644)

			fmt.Printf("  Compiled: %s -> %s\n", path, destFile)
		}
		return nil
	})

	if err != nil {
		fmt.Printf("Build failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Build completed successfully!")
}

func copyFile(src, dst string) error {
	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer source.Close()

	dest, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer dest.Close()

	_, err = io.Copy(dest, source)
	return err
}

func runProject() {
	needsRebuild := false

	if _, err := os.Stat(buildDir); os.IsNotExist(err) {
		needsRebuild = true
	} else {
		hashes := make(map[string]string)

		err := filepath.Walk("src", func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(path, ".py") {
				hash, err := hashFile(path)
				if err != nil {
					return err
				}
				hashes[path] = hash
			}
			return nil
		})

		if err != nil {
			fmt.Printf("Error checking files: %v\n", err)
			os.Exit(1)
		}

		for pyFile, hash := range hashes {
			relPath, _ := filepath.Rel("src", pyFile)
			baseName := strings.TrimSuffix(filepath.Base(pyFile), ".py")
			pycFile := filepath.Join(buildDir, filepath.Dir(relPath), baseName+".pyc")

			if _, err := os.Stat(pycFile); os.IsNotExist(err) {
				needsRebuild = true
				break
			}

			hashFile := pycFile + ".hash"
			if _, err := os.Stat(hashFile); os.IsNotExist(err) {
				needsRebuild = true
				break
			}

			savedHash, err := os.ReadFile(hashFile)
			if err != nil {
				needsRebuild = true
				break
			}

			if string(savedHash) != hash {
				needsRebuild = true
				break
			}
		}
	}

	if needsRebuild {
		fmt.Println("Changes detected, rebuilding...")
		buildProject()
	} else {
		fmt.Println("No changes detected, using cached build")
	}

	fmt.Println("\nRunning project...\n")
	cmd := exec.Command("python3", filepath.Join(buildDir, "main.pyc"))
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Printf("\nRuntime error: %v\n", err)
		os.Exit(1)
	}
}

func hashFile(filepath string) (string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}

	return hex.EncodeToString(hash.Sum(nil)), nil
}

func cleanProject() {
	if err := os.RemoveAll(buildDir); err != nil {
		fmt.Printf("Error cleaning: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("Build artifacts cleaned!")
}
