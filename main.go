package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/bclicn/color"
	"github.com/klauspost/cpuid"
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
	cpu_name := cpuid.CPU.BrandName
	cpu_name = strings.Replace(cpu_name, "(R)", "", -1)
	cpu_name = strings.Replace(cpu_name, "(TM)", "", -1)

	return fmt.Sprint(cpu_name)
}

func get_distro() string {
	command := exec.Command("/usr/bin/lsb_release", "-i")

	distro, err := command.Output()
	if err != nil {
		fmt.Println(err)
	}

	distro = distro[16:] // strip "Distributor ID:"
	return string(distro)
}

func main() {
	user := get_username()
	host := get_hostname()
	desktop_environment := get_desktop_env()
	cpu := get_cpu()
	distro := get_distro()

	fmt.Printf("distro: %s", color.Red(distro))
	fmt.Printf("host: %s@%s\n", color.Red(user), color.Red(host))
	fmt.Printf("de: %s\n", desktop_environment)
	fmt.Printf("cpu: %s\n", cpu)
}
