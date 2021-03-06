package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"

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

func get_cpu() (string, string) {
	cpu_name := cpuid.CPU.BrandName
	cpu_name = strings.Replace(cpu_name, "(R)", "", -1)
	cpu_name = strings.Replace(cpu_name, "(TM)", "", -1)

	var color string = "\033[0m"

	switch cpuid.CPU.VendorString {
	case "GenuineIntel":
		color = "\033[34m"
	case "AuthenticAMD":
		color = "\033[31m"
	}

	return fmt.Sprint(cpu_name), color
}

func get_distro() string {
	file, err := os.Open("/etc/os-release")
	if err != nil {
		fmt.Println(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)
	var distro string

	for scanner.Scan() {
		text := scanner.Text()

		if strings.HasPrefix(text, "ID=") {
			distro = text[3:]
			break
		}
	}

	return strings.ToLower(distro)
}

func get_pkg_count(distro string) string {
	var cmd []string

	switch distro {
	case "debian":
		cmd = strings.Split("/usr/bin/dpkg --list", " ")
	case "fedora":
		cmd = strings.Split("/usr/bin/dnf list installed", " ")
	case "arch":
		cmd = strings.Split("/usr/bin/pacman -Q", " ")
	default:
		return "unknown"
	}

	command := exec.Command(cmd[0])
	for i := 1; i < len(cmd); i++ {
		command.Args = append(command.Args, cmd[i])
	}

	output, err := command.Output()
	if err != nil {
		fmt.Println(err)
	}

	pkg_count := len(strings.Split(string(output), "\n"))
	return fmt.Sprint(pkg_count)
}

func get_uptime() string {
	command := exec.Command("/usr/bin/uptime", "-p")

	uptime, err := command.Output()
	if err != nil {
		fmt.Println(err)
	}

	return strings.Replace(string(uptime), "\n", "", -1)
}
