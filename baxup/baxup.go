package main

// Baxup v0.85.2

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"flag"
	"log"
	"strings"
	//"archive/tar"
	//"compress/gzip"
)
const (
	//// IMPORTANT: In order to change the default docker backup & docker compose container root paths, these variables must be changed manually and the executable rebuilt!
	//// Don't forget the trailing `/`

	// Defines the root directory for all docker compose container directories, as well as the directory to store compressed backups 
	// all docker containers must be located at `/opt/docker/container1`,`/opt/docker/container2`, etc., change this root path below:
	DefaultTargetRoot = "/opt/docker/" 
	// Local location to store backed-up tarballs
	DefaultBackupRoot = "/opt/docker-backups/"
)
// declare config struct
type Config struct {
	TargetRootPath  string
	BackupRootPath  string
	TargetName      string
	RemoteUser      string
	RemoteHost      string
	RemoteSend      bool
	DockerEnabled   bool
}

// runCommand function, requires command name (docker, tar, etc); accepts multiple arguments
func runCommand(commandName string, args ...string) error {
	cmd := exec.Command(commandName, args...)
	
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// collects docker image information and digests, stores alongside `docker-compose.yml` file in newly compressed tarball
func getDockerImages(composeFile string, outputFile string) error {
	cmd := exec.Command("docker", "compose", "-f", composeFile, "images", "--quiet")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("Failed to get docker images: %v", err)
	}

	// loop over image ids to gather docker image digests
	imageLines := string(output)
	imageList := strings.Split(imageLines, "\n")
	var imageInfo string

	for _, imageID := range imageList {
		if imageID == "" {
			continue
		}
		
		// actually get image digest
		cmdInspect := exec.Command("docker", "inspect", "--format", "{{index .RepoDigests 0}}", imageID)
		digestOutput, err := cmdInspect.Output()
		if err != nil {
			return fmt.Errorf("Failed to inspect docker images: %v", err)
		}
		imageInfo += fmt.Sprintf("Image: %s Digest: %s\n", imageID, digestOutput)
	}

	err = os.WriteFile(outputFile, []byte(imageInfo), 0644)
	if err != nil {
		return fmt.Errorf("Failed to write docker image version info to file: %v", err)
	}

	fmt.Println("Docker image information and digests saved to", outputFile)
	return nil
}

func checkDockerRunState(composeFile string) (bool, error) {
	cmd := exec.Command("docker", "compose", "-f", composeFile, "ps", "--services", "--filter","--status=running")
	output, err := cmd.Output()
	if err != nil {
		return false, fmt.Errorf("Failed to check Docker container status: %v", err)
	}
	runningServices := strings.TrimSpace(string(output))
	if runningServices == "" {
		return false, nil
	}
	return true, nil
}

func main () {

	var err error

	// Defines config values
	config := Config{
		TargetRootPath: DefaultTargetRoot,
		BackupRootPath: DefaultBackupRoot,
	}

	// Flag definitions
	targetName := flag.String("target", "", "Name of target directory, not including root path")
	remoteSend := flag.Bool("remote-send", false, "Enable sending backup file to remote machine. Additional flags needed.")
	remoteUser := flag.String("remote-user", "", "Remote machine username. SSH key required.")
	remoteHost := flag.String("remote-host", "", "Remote machine IP address. SSH key required.")
	remoteFile := flag.String("remote-file", "", "Remote filepath. Defaults to /home/$USER/$TARGETNAME.bak.tar.gz")
	dockerBool := flag.Bool("docker", true, "Docker target? Default: true")

	flag.Parse()

	sourceDir := filepath.Join(config.TargetRootPath, *targetName)
	backupFile := filepath.Join(config.BackupRootPath, *targetName+".bak.tar.gz")
	imageVersionFile := filepath.Join(sourceDir, "docker-image-versions.txt")

	if *targetName == "" {
		fmt.Println("Target must be specified!")
		fmt.Println("Exiting..")
		os.Exit(1)
	}

	if *dockerBool {
		composeFilePath := filepath.Join(sourceDir, "docker-compose.yml")

		// Check that docker-compose.yml file exists
		if _, err := os.Stat(composeFilePath); os.IsNotExist(err) {
			log.Fatalf("docker-compose.yml not found at %s", composeFilePath)
		}
		// Ensure Docker container is running
		running, err := checkDockerRunState(composeFilePath)
		if err != nil {
			log.Fatalf("Error checking Docker container status: %v", err)
		}
		if !running {
			log.Fatalf("FATAL ERROR: Docker container is not running or not locateable, exiting!")
		}

		// Get Docker image information & store it in the working dir
		fmt.Println("-------------------------------------------------------------------------")
		fmt.Println("Getting Docker image versions . . .")
		err = getDockerImages(filepath.Join(sourceDir, "docker-compose.yml"), imageVersionFile)
		if err != nil {
			log.Fatalf("Error retrieving Docker image versions: %v", err)
		}

		// Stop docker container
		fmt.Println("-------------------------------------------------------------------------")
		fmt.Println("Stopping Docker container . . .")
		fmt.Println("Issuing docker compose down on ", filepath.Join(sourceDir, "docker-compose.yml"))
		fmt.Println("-------------------------------------------------------------------------")
		err = runCommand("docker", "compose", "-f", filepath.Join(sourceDir, "docker-compose.yml"), "down")
		if err != nil {
			log.Fatalf("Error stopping Docker container: %v", err)
		}
	}

	// Compress target directory
	fmt.Println("-------------------------------------------------------------------------")
	fmt.Println("Compressing container directory . . .")
	fmt.Println("-------------------------------------------------------------------------")
	err := runCommand("tar", "-cvzf", backupFile, "-C", config.TargetRootPath, *targetName)
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

		// Set default remote file path to remote user's homedir if none is specified
		remoteFilePath := *remoteFile
		if remoteFilePath == "" {
			remoteFilePath = fmt.Sprintf("/home/%s/%s.bak.tar.gz", *remoteUser, *targetName)
		}

		fmt.Println("Copying to remote machine . . .")
		// Checksum forced
		rsyncArgs := []string{
			"-avz", "--checksum", "-e", "ssh", backupFile, fmt.Sprintf("%s@%s:%s", *remoteUser, *remoteHost, *remoteFilePath),
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
