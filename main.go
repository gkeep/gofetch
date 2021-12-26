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
	separator string `yaml:"separator"`
	distro    string `yaml:"distro"`
}

type Colors struct {
	main  string
	reset string
}

func load_config() Config {
	f, err := os.Open("cfg.yml")
	if err != nil {
		fmt.Println(err)
	}
	defer f.Close()

	var cfg Config
	decoder := yaml.NewDecoder(f)
	err = decoder.Decode(&cfg)
	if err != nil {
		fmt.Println(err)
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
	return strings.Replace(string(distro), "\n", "", -1)
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
	default:
		fmt.Println("err")
	}

	clrs.reset = "\033[0m"

	return clrs
}

func color_print(col Colors, text string) string {
	return fmt.Sprintf("%s%s%s", col.main, text, col.reset)
}

func main() {
	config := load_config()
	distro := get_distro()

	var colors Colors

	if config.distro != "" {
		colors = get_colors(config.distro)
	} else {
		colors = get_colors(distro)
	}

	user := get_username()
	host := get_hostname()
	desktop_environment := get_desktop_env()
	cpu := get_cpu()

	fmt.Printf("distro %s %s\n", config.separator, color_print(colors, distro))
	fmt.Printf("host %s %s@%s\n", config.separator, color_print(colors, user), color_print(colors, host))
	fmt.Printf("de %s %s\n", config.separator, desktop_environment)
	fmt.Printf("cpu %s %s\n", config.separator, cpu)
}
