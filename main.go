package main

import (
	"fmt"
	"github.com/klauspost/cpuid"
	"os"
	"os/user"
	"runtime"
)

func get_hostname() string {
	hostname, err := os.Hostname()

	if err != nil {
		fmt.Println(err)
	}

	return string(hostname)
}

func get_username() string {
	username, err := user.Current()

	if err != nil {
		fmt.Println(err)
	}

	return string(username.Username)
}

func get_desktop_env() string {
	de := os.Getenv("XDG_CURRENT_DESKTOP")

	if de == "" {
		return "None"
	}

	return de
}

func get_cpu() string {
	cores_count := runtime.NumCPU()
	cpu_name := cpuid.CPU.BrandName

	return fmt.Sprint(cpu_name, " x", cores_count)
}

func main() {
	user := get_username()
	host := get_hostname()
	desktop_environment := get_desktop_env()
	cpu := get_cpu()

	fmt.Printf("host: %s@%s\n", user, host)
	fmt.Printf("de: %s\n", desktop_environment)
	fmt.Printf("cpu: %s\n", cpu)
}
