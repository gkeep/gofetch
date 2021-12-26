package main

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"

	"github.com/klauspost/cpuid"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Separator string `yaml:"Separator"`
	Distro    string `yaml:"DistroOverride"`
}

type Colors struct {
	main string
	cpu  string
}

func load_config() Config {
	var cfg Config
	file := "cfg.yml"
	f, err := os.Open(file)

	if os.IsNotExist(err) {
		cfg.Separator = ":"
		return cfg
	}

	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Println(err)
	}

	if cfg.Separator == "" {
		cfg.Separator = ":"
	}

	return cfg
}

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

func get_cpu(clrs Colors) (string, Colors) {
	cpu_name := cpuid.CPU.BrandName
	cpu_name = strings.Replace(cpu_name, "(R)", "", -1)
	cpu_name = strings.Replace(cpu_name, "(TM)", "", -1)

	switch cpuid.CPU.VendorString {
	case "Intel":
		clrs.cpu = "\033[34m"
	case "AuthenticAMD":
		clrs.cpu = "\033[31m"
	default:
		clrs.cpu = "\033[0m"
	}

	return fmt.Sprint(cpu_name), clrs
}

func get_distro() string {
	command := exec.Command("/usr/bin/lsb_release", "-i")

	distro, err := command.Output()
	if err != nil {
		fmt.Println(err)
	}

	distro = distro[16:] // strip "Distributor ID:"
	return strings.Replace(string(distro), "\n", "", -1)
}

func get_uptime() string {
	command := exec.Command("/usr/bin/uptime", "-p")

	uptime, err := command.Output()
	if err != nil {
		fmt.Println(err)
	}

	return strings.Replace(string(uptime), "\n", "", -1)
}

func get_colors(distro string) Colors {
	var clrs Colors

	switch strings.ToLower(distro) {
	case "debian":
		clrs.main = "\033[31m" // red
	case "arch":
		clrs.main = "\033[1;34m" // light blue
	case "fedora":
		clrs.main = "\033[34m" // blue
	}

	return clrs
}

func color_print(color string, text string) string {
	col_reset := "\033[0m"

	return fmt.Sprintf("%s%s%s", color, text, col_reset)
}

func main() {
	config := load_config()
	distro := get_distro()

	var colors Colors

	if config.Distro != "" {
		colors = get_colors(config.Distro)
	} else {
		colors = get_colors(distro)
	}

	user := get_username()
	host := get_hostname()
	desktop_environment := get_desktop_env()
	cpu, colors := get_cpu(colors)
	uptime := get_uptime()

	fmt.Printf("distro %s %s\n", config.Separator, color_print(colors.main, distro))
	fmt.Printf("host %s %s@%s\n", config.Separator, color_print(colors.main, user), color_print(colors.main, host))
	fmt.Printf("de %s %s\n", config.Separator, desktop_environment)
	fmt.Printf("cpu %s %s\n", config.Separator, color_print(colors.cpu, cpu))
	fmt.Printf("uptime %s %s\n", config.Separator, uptime)
}
