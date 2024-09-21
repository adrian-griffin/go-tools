package main

import (
	"fmt"
	//"os"
	"os/exec"
	//"os/user"
	"path/filepath"
	"flag"
	"log"
)

// runCommand function, requires command name (docker); accepts multiple arguments
func runCommand(commandName string, args ...string) error {
	cmd := exec.Command(commandName, args...)
	cmd.Stdout = nil
	cmd.Stderr = nil
	return cmd.Run()
}


func main () {
	//// IMPORTANT: In order to change the default docker backup & docker compose container root paths, these variables much be changed manually and the script rebuilt!
	//// Don't forget the trailing `/`

	// Defines the root directories for all docker compose container directories, as well as the directory to store compressed backups 
	dockerComposeRootPath := "/opt/docker/"
	dockerBackupRootPath := "/opt/docker-backups/"

	// Flag definitions
	dockerContainerName := flag.String("docker-target", "vaultwarden", "Name of docker container's directory to back up.")
	remoteSend := flag.Bool("remote-send", false, "Enable sending backup .tar.gz file to remote server. Additional flags needed.")
	remoteUser := flag.String("remote-user", "agriffin", "Remote machine username. SSH key required.")
	remoteHost := flag.String("remote-host", "10.115.0.1", "Remote machine IP address.  SSH key required.")
	remoteFile := flag.String("remote-file", "", "Remote filepath. Defaults to /home/$USER/$DOCKERNAME.bak.tar.gz")

	flag.Parse()
	
	sourceDir := filepath.Join(dockerComposeRootPath, *dockerContainerName, "/")
	backupFile := filepath.Join(dockerBackupRootPath, *dockerContainerName+".bak.tar.gz")

	// If remoteFile not specified, defaults to /home/$USER/$DOCKERNAME.bak.tar.gz
	if *remoteFile == "" {
		*remoteFile = fmt.Sprintf("/home/%s/%s.bak.tar.gz", *remoteUser, *dockerContainerName)
	}

	// Step 1: Stop docker container
	fmt.Println("Stopping Docker container . . .")
	fmt.Println("Issuing docker compose down on ", filepath.Join(sourceDir, "docker-compose.yml"))
	err := runCommand("docker", "compose", "-f", filepath.Join(sourceDir, "docker-compose.yml"), "down")
	if err != nil {
		log.Fatalf("Error stopping Docker container: %v", err)
	}

	// Step 2: Compress target directory
	fmt.Println("Compressing container directory . . .")
	err = runCommand("tar", "-cvzf", backupFile, "-C", dockerComposeRootPath, *dockerContainerName)
	fmt.Println("Backup file saved at:", backupFile)
	if err != nil {
		log.Fatalf("Error compressing directory: %v", err)
	}

	// Optional: Rsync to remote destination
	if *remoteSend {
		fmt.Println("Copying to remote machine . . .")
		// Checksum forced
		rsyncArgs := []string{
			"-avz", "--checksum", backupFile, fmt.Sprintf("%s@%s:%s", *remoteUser, *remoteHost, *remoteFile),
		}
		err = runCommand("rsync",rsyncArgs...)
		if err != nil {
			log.Fatalf("Error sending file to remote server: %v", err)
		}
	}

	// Restart docker container
	fmt.Println("Starting Docker container . . .")
	err = runCommand("docker", "compose", "-f", filepath.Join(sourceDir, "docker-compose.yml"), "up", "-d")
	if err != nil {
		log.Fatalf("Error starting Docker container: %v", err)
	}
}