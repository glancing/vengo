package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"bufio"
	"strings"
)

const envDir = ".vengo"

func getEnvHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("❌ Error fetching home directory:", err)
		os.Exit(1)
	}
	return filepath.Join(home, envDir)
}

func ensureVengoDir() {
	dir := getEnvHome()
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Println("❌ Failed to create vengo directory:", err)
			os.Exit(1)
		}
	}
}

func usage() {
	fmt.Println(`🦥 Vengo — Python virtual environment manager
Usage:
  vengo create <env-name>     Create a new virtual environment
  vengo list                  List all virtual environments
  vengo activate <env-name>   Show activation command for environment
  vengo delete <env-name>     Delete an environment`)
}

func main() {
	if len(os.Args) < 2 {
		usage()
		return
	}

	// Ensure that the .vengo directory exists before proceeding.
	ensureVengoDir()
	
	// Add the shell function to the user's shell configuration file.
	addShellFunction()

	cmd := os.Args[1]

	switch cmd {
	case "create":
		if len(os.Args) != 3 {
			fmt.Println("Usage: vengo create <env-name>")
			return
		}
		createEnv(os.Args[2])
	case "list":
		listEnvs()
	case "activate":
		if len(os.Args) != 3 {
			fmt.Println("Usage: vengo activate <env-name>")
			return
		}
		activateEnv(os.Args[2])
	case "delete":
		if len(os.Args) != 3 {
			fmt.Println("Usage: vengo delete <env-name>")
			return
		}
		deleteEnv(os.Args[2])
	default:
		usage()
	}
}

func checkIfEnvExists(name string) bool {
	envPath := filepath.Join(getEnvHome(), name)
	if _, err := os.Stat(envPath); !os.IsNotExist(err) {
		return true
	}
	return false
}

func createEnv(name string) {
	envPath := filepath.Join(getEnvHome(), name)

	// Check if environment already exists.
	if checkIfEnvExists(name) {
		fmt.Println("⚠️ Environment already exists.")
		return
	}

	// Run the command to create a new virtual environment.
	cmd := exec.Command("python3", "-m", "venv", envPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println("❌ Failed to create virtual environment:", err)
		return
	}

	fmt.Printf("✅ Created virtual environment '%s'\n", name)
}

func activateEnv(name string) {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("❌ Error getting current user:", err)
		return
	}
	shell := os.Getenv("SHELL")
	var rcFile string
	if strings.Contains(shell, "zsh") {
		rcFile = filepath.Join(usr.HomeDir, ".zshrc")
	} else {
		rcFile = filepath.Join(usr.HomeDir, ".bashrc")
	}
	if !isShellFunctionAdded(rcFile, "vengo() {") {
		fmt.Println("⚠️ Shell function not added. Run 'vengo' to add it.")
		return
	}
	if !checkIfEnvExists(name) {
		fmt.Println("⚠️ Environment does not exist.")
		return
	}
}

func listEnvs() {
	envs, err := os.ReadDir(getEnvHome())
	if err != nil {
		fmt.Println("❌ Error reading environments:", err)
		return
	}

	if len(envs) == 0 {
		fmt.Println("⚠️ No virtual environments created yet.")
		return
	}

	fmt.Println("📌 Your virtual environments:")
	for _, env := range envs {
		if env.IsDir() {
			fmt.Println(" -", env.Name())
		}
	}
}

func deleteEnv(name string) {
	path := filepath.Join(getEnvHome(), name)

	if !checkIfEnvExists(name) {
		fmt.Println("⚠️ Environment does not exist.")
		return
	}

	fmt.Printf("⚠️ Are you sure you want to delete '%s'? [y/N]: ", name)
	var confirm string
	fmt.Scanln(&confirm)
	if confirm != "y" && confirm != "Y" {
		fmt.Println("❌ Aborted deletion.")
		return
	}

	if err := os.RemoveAll(path); err != nil {
		fmt.Println("❌ Failed to delete environment:", err)
		return
	}

	fmt.Printf("🗑  Deleted environment '%s'.\n", name)
}

func addShellFunction() {
	usr, err := user.Current()
	if err != nil {
		fmt.Println("❌ Error getting current user:", err)
		return
	}

	shell := os.Getenv("SHELL")
	var rcFile string
	if strings.Contains(shell, "zsh") {
		rcFile = filepath.Join(usr.HomeDir, ".zshrc")
	} else {
		rcFile = filepath.Join(usr.HomeDir, ".bashrc")
	}

	shellFunction := `
vengo() {
    if [ "$1" = "activate" ]; then
        if [ -z "$2" ]; then
            echo "Usage: vengo activate <env-name>"
            return 1
        fi
        source "$HOME/.vengo/$2/bin/activate"
    else
        command vengo "$@"
    fi
}
`

	// Check if the shell function is already added
	if !isShellFunctionAdded(rcFile, "vengo() {") {
		file, err := os.OpenFile(rcFile, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			fmt.Println("❌ Error opening file for writing:", err)
			return
		}
		defer file.Close()

		if _, err := file.WriteString(shellFunction); err != nil {
			fmt.Println("❌ Error writing to file:", err)
			return
		}

		fmt.Println("✅ Shell function added to", rcFile)
		fmt.Println("Please restart your terminal or run 'source", rcFile, "' to apply the changes.")
	}
}

func isShellFunctionAdded(rcFile, functionSignature string) bool {
	file, err := os.Open(rcFile)
	if err != nil {
		fmt.Println("❌ Error opening file:", err)
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), functionSignature) {
			return true
		}
	}

	if err := scanner.Err(); err != nil {
		fmt.Println("❌ Error reading file:", err)
		return false
	}

	return false
}