package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"flag"
	"log"
)

// runCommand function, requires command name (docker); accepts multiple arguments
func runCommand(commandName string, args ...string) error {
	cmd := exec.Command(commandName, args...)
	cmd.Stdout = nil
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func main () {
	//// IMPORTANT: In order to change the default docker backup & docker compose container root paths, these variables much be changed manually and the script rebuilt!
	//// Don't forget the trailing `/`

	// Defines the root directories for all docker compose container directories, as well as the directory to store compressed backups 
	targetRootPath := "/opt/docker/"
	backupRootPath := "/opt/docker-backups/"

	// Flag definitions
	targetName := flag.String("target", "", "Name of target directory, not including rooth path")
	remoteSend := flag.Bool("remote-send", false, "Enable sending backup file to remote machine. Additional flags needed.")
	remoteUser := flag.String("remote-user", "", "Remote machine username. SSH key required.")
	remoteHost := flag.String("remote-host", "", "Remote machine IP address. SSH key required.")
	remoteFile := flag.String("remote-file", "", "Remote filepath. Defaults to /home/$USER/$TARGETNAME.bak.tar.gz")
	dockerBool := flag.Bool("docker", true, "Docker target? Default: true")

	flag.Parse()

	sourceDir := filepath.Join(targetRootPath, *targetName)
	backupFile := filepath.Join(backupRootPath, *targetName+".bak.tar.gz")

	if *targetName == "" {
		fmt.Println("Target must be specified!")
		fmt.Println("Exiting..")
		os.Exit(1)
	}

	if *dockerBool {
		// Stop docker container
		fmt.Println("-------------------------------------------------------------------------")
		fmt.Println("Stopping Docker container . . .")
		fmt.Println("Issuing docker compose down on ", filepath.Join(sourceDir, "docker-compose.yml"))
		fmt.Println("-------------------------------------------------------------------------")
		err := runCommand("docker", "compose", "-f", filepath.Join(sourceDir, "docker-compose.yml"), "down")
		if err != nil {
			log.Fatalf("Error stopping Docker container: %v", err)
		}
	}

	// Compress target directory
	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("Compressing container directory . . .")
	fmt.Println("-------------------------------------------------------------------------")
	err := runCommand("tar", "-cvzf", backupFile, "-C", targetRootPath, *targetName)
	if err != nil {
		log.Fatalf("Error compressing directory: %v", err)
	}
	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("Backup file saved at:", backupFile)
	fmt.Println("-------------------------------------------------------------------------")

	// Optional: Rsync to remote destination
	if *remoteSend {
		if *remoteUser == "" || *remoteHost == "" {
			log.Fatalf("Remote user and host must be specified when sending to a remote machine.")
		}
		fmt.Println("Copying to remote machine . . .")
		// Checksum forced
		rsyncArgs := []string{
			"-avz", "--checksum", "-e", "ssh", backupFile, fmt.Sprintf("%s@%s:%s", *remoteUser, *remoteHost, *remoteFile),
		}
		err = runCommand("rsync", rsyncArgs...)
		if err != nil {
			log.Fatalf("Error sending file to remote server: %v", err)
		}
	}

	if *dockerBool {
		// Restart docker container
		fmt.Println("-------------------------------------------------------------------------")
		fmt.Println("Starting Docker container . . .")
		fmt.Println("-------------------------------------------------------------------------")
		err = runCommand("docker", "compose", "-f", filepath.Join(sourceDir, "docker-compose.yml"), "up", "-d")
		if err != nil {
			log.Fatalf("Error starting Docker container: %v", err)
		}
	}
}
