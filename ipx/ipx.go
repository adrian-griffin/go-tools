package main

import (
    //"flag"
    "fmt"
    "io"
    "net/http"
    "os/exec"
    "strings"
)

func main () {
    gateway := getDefaultGW()
    lanInt, lanIp := getLan()
    dns := getDnsServers()
    natIp := getPublicIP()


    fmt.Println("--- Default route & Network Information")
    fmt.Println("Default Gateway:", gateway)
    fmt.Println("LAN IP:", lanIp)
    fmt.Println("LAN Interface:", lanInt)
    fmt.Println("DNS Servers:", dns)
    fmt.Println("NAT/Public IP:", natIp)
    fmt.Println("---------------------------------------")
}

func getDefaultGW() string {
    // Collects default gateway information from `ip route` & grabs IP segments
    output, _ := exec.Command("sh", "-c", "ip route | grep default").Output()
    segments := strings.Fields(string(output))
    return segments[2]
}

func getLan() (string, string) {
    // Checks what interface traffic takes to reach the internet (8.8.8.8); collects interface name & lan IP address
    output, _ := exec.Command("sh", "-c", "ip route get 8.8.8.8").Output()
    segments := strings.Fields(string(output))
    return segments[4], segments[6]
}

func getDnsServers() string {
    // Checks resolv.conf for nameserver line
    output, _ := exec.Command("sh", "-c", "cat /etc/resolv.conf | grep nameserver").Output()
    segments := strings.Fields(string(output))
    return segments[1]
}

func getPublicIP() string {
    // Performs http get to collect ifconfig.me ip, returns ip
    ifConfig, _ := http.Get("https://ifconfig.me")
    defer ifConfig.Body.Close()
    ip, _ := io.ReadAll(ifConfig.Body)
    return string(ip)
}