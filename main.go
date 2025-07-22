package main

import (
    "fmt"
    "log"
    "net"
    "net/http"
    "os"
    "os/exec"
    "runtime"
    "time"
    "github.com/fatih/color"
    "waguri-auth/api"
    "strings"
    "waguri-auth/common"
)

// Clear terminal screen
func clearTerminal() {
    switch runtime.GOOS {
    case "windows":
        cmd := exec.Command("cmd", "/c", "cls")
        cmd.Stdout = os.Stdout
        cmd.Run()
    default:
        cmd := exec.Command("clear")
        cmd.Stdout = os.Stdout
        cmd.Run()
    }
}

func getLocalIP() string {
    interfaces, err := net.Interfaces()
    if err != nil {
        return "Unable to get interfaces"
    }

    for _, iface := range interfaces {
        // Skip loopback and down interfaces
        if iface.Flags&net.FlagUp == 0 || iface.Flags&net.FlagLoopback != 0 {
            continue
        }

        addrs, err := iface.Addrs()
        if err != nil {
            continue
        }

        for _, addr := range addrs {
            var ip net.IP

            switch v := addr.(type) {
            case *net.IPNet:
                ip = v.IP
            case *net.IPAddr:
                ip = v.IP
            }

            if ip == nil || ip.IsLoopback() {
                continue
            }

            ip = ip.To4()
            if ip == nil {
                continue // not an IPv4 address
            }

            // Skip link-local 169.254.x.x
            if strings.HasPrefix(ip.String(), "169.254") {
                continue
            }

            return ip.String()
        }
    }

    return "No valid IP found"
}

// Timestamp formatter
func logPrefix() string {
    return fmt.Sprintf("[%s] [?]  ~", time.Now().Format("15:04:05"))
}

func main() {
    common.LoadConfig()

    // Clear screen
    clearTerminal()

    cyan := color.New(color.FgCyan).SprintFunc()
    green := color.New(color.FgHiGreen).SprintFunc()
    yellow := color.New(color.FgHiYellow).SprintFunc()

    // Log styled info
    fmt.Printf("%s %s %s\n", cyan(logPrefix()), green("INFO"), "- LOADED CONFIG")
    fmt.Printf("%s %s %s\n", cyan(logPrefix()), green("IP"), "- "+yellow(getLocalIP()))
    fmt.Printf("%s %s %s\n", cyan(logPrefix()), green("PORT"), "- "+yellow(common.AppConfig.Port))

    // Register routes
    http.HandleFunc("/create-key", api.AuthMiddleware(api.CreateKeyHandler))
    http.HandleFunc("/delete-key", api.AuthMiddleware(api.DeleteKeyHandler))
    http.HandleFunc("/login", api.LoginHandler)

    log.Fatal(http.ListenAndServe(":"+common.AppConfig.Port, nil))
}
