package main

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Separator string `yaml:"Separator"`
	Distro    string `yaml:"DistroOverride"`
}

type Colors struct {
	main      string
	secondary string
	cpu       string
}

type Info struct {
	user        string
	host        string
	cpu         string
	uptime      string
	desktop_env string
}

func load_config() Config {
	var cfg Config

	cfg_dir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println(err)
	}

	file := cfg_dir + "/gofetch/gofetch.yml"
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

	return cfg
}

func get_colors(distro string) Colors {
	var clrs Colors

	switch strings.ToLower(distro) {
	case "debian":
		clrs.main = "\033[31m"     // red
		clrs.secondary = "\033[1m" // white
	case "arch":
		clrs.main = "\033[36m" // light blue
		clrs.secondary = "\033[1m"
	case "fedora":
		clrs.main = "\033[34m"      // blue
		clrs.secondary = "\033[35m" // purple
	}

	return clrs
}

func color_print(color string, text string) string {
	col_reset := "\033[0m"

	return fmt.Sprintf("%s%s%s", color, text, col_reset)
}

func fetch() (Info, string) {
	var data Info
	var cpu_color string

	data.user = get_username()
	data.host = get_hostname()
	data.desktop_env = get_desktop_env()
	data.cpu, cpu_color = get_cpu()
	data.uptime = get_uptime()

	return data, cpu_color
}

func main() {
	config := load_config()
	distro := get_distro()

	var data Info
	var colors Colors

	if config.Distro != "" {
		colors = get_colors(config.Distro)
	} else {
		colors = get_colors(distro)
	}

	data, cpu_color := fetch()
	colors.cpu = cpu_color

	info_list := [5]string{"distro", "host", "de", "cpu", "uptime"}

	for i := 0; i < 5; i++ {
		var template string

		switch info_list[i] {
		case "distro":
			template = color_print(colors.main, distro)
		case "host":
			template = fmt.Sprintf("%s@%s", color_print(colors.main, data.user), color_print(colors.main, data.host))
		case "de":
			template = data.desktop_env
		case "cpu":
			template = color_print(colors.cpu, data.cpu)
		case "uptime":
			template = data.uptime
		}

		fmt.Printf("%s %s %s\n", color_print(colors.secondary, info_list[i]),
			config.Separator, template)
	}
}
