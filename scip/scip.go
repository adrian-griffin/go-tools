package main

import (
    "fmt"
    "golang.org/x/crypto/ssh"
    "golang.org/x/crypto/ssh/agent"
    "github.com/schollz/progressbar/v3"
    "io"
    "net"
    "os"
    "time"
)

func main() {
    sourceFile := "path/to/source/file"
    destFile := "path/to/dest/file"
    user := "username"
    host := "destination_host"
    port := "22"
    jumpHost := "jumphost" // Set to "" if not using a jump host

    sshConfig := &ssh.ClientConfig{
        User:            user,
        Auth:            []ssh.AuthMethod{sshAgent()},
        HostKeyCallback: ssh.InsecureIgnoreHostKey(), // Change this to proper host key checking in production
        Timeout:         5 * time.Second,
    }

    var conn *ssh.Client
    var err error

    if jumpHost != "" {
        // Establish connection to jump host
        jumpConn, err := ssh.Dial("tcp", jumpHost+":22", sshConfig)
        if err != nil {
            fmt.Printf("Failed to connect to jump host: %v\n", err)
            return
        }
        defer jumpConn.Close()

        // Dial the final destination through the jump host
        finalConn, err := jumpConn.Dial("tcp", host+":"+port)
        if err != nil {
            fmt.Printf("Failed to dial final destination through jump host: %v\n", err)
            return
        }

        // Establish an SSH connection over the tunneled connection
        connSSH, chans, reqs, err := ssh.NewClientConn(finalConn, host+":"+port, sshConfig)
        if err != nil {
            fmt.Printf("Failed to create SSH connection to final destination: %v\n", err)
            return
        }

        conn = ssh.NewClient(connSSH, chans, reqs)
    } else {
        conn, err = ssh.Dial("tcp", host+":"+port, sshConfig)
        if err != nil {
            fmt.Printf("Failed to connect to destination host: %v\n", err)
            return
        }
    }

    defer conn.Close()

    session, err := conn.NewSession()
    if err != nil {
        fmt.Printf("Failed to create SSH session: %v\n", err)
        return
    }
    defer session.Close()

    // Progress bar setup
    file, err := os.Open(sourceFile)
    if err != nil {
        fmt.Printf("Failed to open source file: %v\n", err)
        return
    }
    defer file.Close()

    fileInfo, err := file.Stat()
    if err != nil {
        fmt.Printf("Failed to get file info: %v\n", err)
        return
    }

    bar := progressbar.DefaultBytes(
        fileInfo.Size(),
        "copying file",
    )

    pipeReader, pipeWriter := io.Pipe()
    go func() {
        io.Copy(pipeWriter, file)
        pipeWriter.Close()
    }()

    session.Stdin = io.TeeReader(pipeReader, bar)
    session.Stdout = os.Stdout
    session.Stderr = os.Stderr

    err = session.Run(fmt.Sprintf("scp -t %s", destFile))
    if err != nil {
        fmt.Printf("Failed to run SCP: %v\n", err)
    }
}

func sshAgent() ssh.AuthMethod {
    socket := os.Getenv("SSH_AUTH_SOCK")
    if socket == "" {
        return nil
    }

    conn, err := net.Dial("unix", socket)
    if err != nil {
        return nil
    }

    return ssh.PublicKeysCallback(agent.NewClient(conn).Signers)
}
