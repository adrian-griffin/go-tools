package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"bufio"
	//"strings"
)

// prompt for input
func promptInput(prompt string) string {
	fmt.Print(prompt)
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

// configure Git credential helper
func configureGitCredentialCache() error {
	cmd := exec.Command("git", "config", "--global", "credential.helper", "cache --timeout=600")
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to configure Git credential cache: %w", err)
	}
	return nil
}

// get distro
func getDistro() string {
	releaseInfo, err := os.ReadFile("/etc/os-release")
	if err != nil {
		return "unknown"
	}
	info := string(releaseInfo)
	case contains(info, "Arch") || contains(info, "Manjaro"):
		return "arch"
	case contains(info, "Debian") || contains(info, "Ubuntu"):
		return "debian"
	default:
		return "unknown"
}

// clone or pull the Git repository
func gitCloneRepo(repoURL, destination string) error {
	// Check if destination directory exists
	if _, err := os.Stat(destination); !os.IsNotExist(err) {
		// Destination exists, perform a git pull to update the repo
		pullCmd := exec.Command("git", "-C", destination, "pull")
		output, err := pullCmd.CombinedOutput()
		if err != nil {
			fmt.Println("Failed to pull latest changes in", destination, "Error:", err)
			fmt.Println(string(output)) // Print output to understand what went wrong
			return err
		}
		fmt.Println("Git repo", repoURL, "pulled successfully to", destination)
		return nil
	}

	authRepoURL := repoURL

	// Destination doesn't exist, perform a git clone
	cmd := exec.Command("git", "clone", "--depth", "1", authRepoURL, destination)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Failed to clone git repo", repoURL, "to", destination, "Error:", err)
		fmt.Println(string(output)) // Print output to understand what went wrong
		return err
	}
	fmt.Println("Git repo", repoURL, "cloned successfully to", destination)
	return nil
}

func moveFile(srcPath string, dstPath string) error {
	// Check if destination path already exists
	//if _, err := os.Stat(dstPath); !os.IsNotExist(err) {
	//	fmt.Println(dstPath, "already exists, skipping")
	//	return nil
	//}
	
	// Moving file from srcPath to dstPath
	cmd := exec.Command("sudo", "mv", srcPath, dstPath)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error moving file to destination", dstPath)
		return err
	}

	fmt.Println("File successfully moved to", dstPath)
	return nil 
}

func aptUpdate() error {
	cmd := exec.Command("sudo", "apt-get", "update", "-y")
	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Println("Successful apt update . . .")
	return nil
}

func pacmanUpdate() error {
	cmd := exec.Command("sudo", "pacman", "-Syu")
	if err := cmd.Run(); err != nil {
		return err
	}
	fmt.Println("Succussful pacman update . . .")
	return nil
}



// install tool packages via cli
func installLinuxTools() error {
	distro := getLinuxDistro()

	switch distro {
	case "debian":
		if err := aptUpdate(); err != nil {
			return err
		}
		cmds := [][]string{
			{"sudo", "apt-get", "install", "-y", "zsh"},
			{"sudo", "apt-get", "install", "-y", "git"},
			{"sudo", "apt-get", "install", "-y", "neovim"},
			{"sudo", "apt-get", "install", "-y", "fzf"}
		}
		for _, cmdArgs := range cmds {
			if err := exec.Command(cmdArgs[0], cmdArgs[1:]...).Run(); err != nil {
				return fmt.Errorf("error installing %s: %w", cmdArgs[3], err)
			}
		}
	case "arch":
		if err := pacmanUpdate(); err != nil {
			return err
		}
		cmds := [][]string{
			{"sudo", "pacman", "-Sy", "--noconfirm", "zsh"},
			{"sudo", "pacman", "-Sy", "--noconfirm", "git"},
			{"sudo", "pacman", "-Sy", "--noconfirm", "neovim"},
			{"sudo", "pacman", "-Sy", "--noconfirm", "fzf"},
		}
		for _, cmdArgs := range cmds {
			if err := exec.Command(cmdArgs[0], cmdArgs[1:]...).Run(); err != nil {
				return fmt.Errorf("error installing %s: %w", cmdArgs[4], err)
			}
		}
	default:
		return fmt.Errorf("unsupported distro for package install")
	}

	fmt.Println("Core packages installed successfully.")
	return nil
}

////
// -- Git Clone Dotties
////
func gitCloneDotties() error {
	// Determine user and generate homedir path
	usr, _ := user.Current()
	homePath := usr.HomeDir

	err := gitCloneRepo("https://github.com/adrian-griffin/dotties.git", filepath.Join(homePath, "dotties"))
	if err != nil {
		return err
	}
	return nil
}

func gitCloneGoTools() error {
	// Determine user and generate homedir path
	usr, _ := user.Current()
	homePath := usr.HomeDir

	err := gitCloneRepo("https://github.com/adrian-griffin/go-tools.git", filepath.Join(homePath, "go-tools"))
	if err != nil {
		return err
	}
	return nil
}

/////
// -- Oh My Zsh Install & Customization
////

func installOhMyZsh() error {
	// Determine current user and OMZ path
	usr, _ := user.Current()
	ohMyZshPath := filepath.Join(usr.HomeDir, ".oh-my-zsh")

	// Check if OMZ is already installed
	if _, err := os.Stat(ohMyZshPath); !os.IsNotExist(err) {
		fmt.Println("Oh My Zsh is already installed at", ohMyZshPath, ", skipping")
		return nil
	}

	// Install Oh My Zsh if not installed
	cmd := exec.Command("sh", "-c", "curl -fsSL https://raw.githubusercontent.com/ohmyzsh/ohmyzsh/master/tools/install.sh | sh")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error generating command for installOhMyZsh:", err)
		fmt.Println(string(output)) // This will print the output and error messages from the command
		return err
	}
	fmt.Println("OhMyZsh Installed . . .")
	return nil
}

func ohMyZshCustomizations() error {
	// Determine user and generate OMZ custom path
	usr, _ := user.Current()
	customPath := filepath.Join(usr.HomeDir, ".oh-my-zsh/custom")

	// Install Zsh Autosuggestions
	err := gitCloneRepo("https://github.com/zsh-users/zsh-autosuggestions", filepath.Join(customPath, "plugins/zsh-autosuggestions"))
	if err != nil {
		return err
	}

	// Install Zsh Syntax Highlighting
	err = gitCloneRepo("https://github.com/zsh-users/zsh-syntax-highlighting.git", filepath.Join(customPath, "plugins/zsh-syntax-highlighting"))
	if err != nil {
		return err
	}

	// Install fzf-zsh-plugin
	err = gitCloneRepo("https://github.com/unixorn/fzf-zsh-plugin.git", filepath.Join(customPath, "plugins/fzf-zsh-plugin"))
	if err != nil {
		return err
	}

	// Mv .zshrc file from dotties dir to homedir

	dotFileZshrc := filepath.Join(usr.HomeDir, "dotties/.zshrc")
	homeDirZshrc := filepath.Join(usr.HomeDir, ".zshrc")

	moveFile(dotFileZshrc, homeDirZshrc)

	return nil
}
/////
// -x
////

/////
// -- Main
////

func main() {
	fmt.Println(" ---------------------------------------------------- ")

	// Install linux packages & tools
	fmt.Println(" --- Installing Linux Packages & Tools")
	if err := installLinuxTools(); err != nil {
		fmt.Println("Failed to install Zsh:", err)
		return
	}
	fmt.Println(" --- Completed Linux Packages & Tools Install")
	fmt.Println(" ---------------------------------------------------- ")

	// Set up git credential cache (takes keyboard input) for private git dotfile repository
	fmt.Println(" --- Configuring Git Credential Cache")
	if err := configureGitCredentialCache(); err != nil {
		fmt.Println("Error configuring Git credential cache:", err)
		return
	}
	fmt.Println(" --- Completed Git Credential Cache Configuration")
	fmt.Println(" ---------------------------------------------------- ")

	// Git clone remote git repo for dotfiles, stores at $HOME
	fmt.Println(" --- Git cloning dotfile repository")
	if err := gitCloneDotties(); err != nil {
		fmt.Println("Failed to git clone dotfiles:", err)
	}
	fmt.Println(" --- Completed dotfile repo clone")
	fmt.Println(" ---------------------------------------------------- ")

	// Install OMZ & Customizations
	fmt.Println(" --- Installing OhMyZsh & Customizations")
	if err := installOhMyZsh(); err != nil {
		fmt.Println("Failed to install Oh My Zsh:", err)
		return
	}
	if err := ohMyZshCustomizations(); err != nil {
		fmt.Println("Failed to install Oh My Zsh customizations:", err)
		return
	}
	fmt.Println(" --- Completed OhMyZsh & Customizations Install")
	fmt.Println(" ---------------------------------------------------- ")

	// Git clone/pull remote git repo for `adrian-griffin/go-tools`, stores at $HOME/go-tools
	fmt.Println(" --- Git cloning go-tools repository")
	if err := gitCloneGoTools(); err != nil {
		fmt.Println("Failed to git clone go-tools:", err)
	}
	fmt.Println(" --- Completed go-tools repo clone")
	fmt.Println(" ---------------------------------------------------- ")

	fmt.Println(" -- All work is complete, please source the shell with `source ~/.zshrc`")
}

/////
// -x
////
