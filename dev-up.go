package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/cheggaaa/pb/v3"
	"gopkg.in/yaml.v2"
)

type EnvVariable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type Config struct {
	Profile string `yaml:"profile"`
	Execute []struct {
		Program string `yaml:"program,omitempty"`
		Command string `yaml:"command,omitempty"`
	} `yaml:"execute"`
	Lifecycle []struct {
		Install struct {
			App     string        `yaml:"app"`
			Version string        `yaml:"version"`
			Env     []EnvVariable `yaml:"env"`
		} `yaml:"install"`
	} `yaml:"lifecycle"`
}

func readConf(filename string) (*Config, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	c := &Config{}
	err = yaml.Unmarshal(buf, c)
	if err != nil {
		return nil, fmt.Errorf("in file %q: %v", filename, err)
	}

	return c, nil
}

func showProgress() func() {
	count := 5000
	// create and start new bar
	bar := pb.StartNew(count)

	for i := 0; i < count; i++ {
		bar.Increment()
		time.Sleep(time.Millisecond)
	}
	//return a function that can be called upon completion of installation
	return func() {
		bar.Finish()
	}
}

func verifyAndinstallHomeBrew() {
	fmt.Println("Verifying HomeBrew installation")
	cmd := exec.Command("brew", "version")
	output, _ := cmd.CombinedOutput()
	if strings.Contains(string(output), "Unknown command") {
		fmt.Println("HomeBrew is not installed , We shall install it...and NO , you have no say in this")
		fmt.Println("OH , and here is a progressbar that means absolutely nothing 😁")
		brewInstallComplete := showProgress()
		brewInstallCommand := exec.Command("/bin/bash", "-c", "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)")
		brewInstalloutput, _ := brewInstallCommand.CombinedOutput()
		log.Println(string(brewInstalloutput))
		brewInstallComplete()
	}
}

func verifyAndInstall(name string, version string) {
	fmt.Println("Verifying and installing " + name + " " + version)
	cmd := exec.Command(name, "--version")
	output, _ := cmd.CombinedOutput()
	fmt.Println(string(output))
	if strings.Contains(string(output), "Unknown command") {
		fmt.Println(name + "is not installed , We shall install it...and NO , you have no say in this")
		fmt.Println("OH , and here is a progressbar that means absolutely nothing 😁")
		brewInstallComplete := showProgress()
		brewInstallCommand := exec.Command("brew", "install", name+"@"+version)
		brewInstalloutput, _ := brewInstallCommand.CombinedOutput()
		log.Println(string(brewInstalloutput))
		brewInstallComplete()
	}
}

func writeToRcFile(env []EnvVariable) {
	var str strings.Builder

	for i := 0; i < len(env); i++ {
		str.WriteString(env[i].Name + "=" + env[i].Value + "\n")
		fmt.Println(i, env[i])
	}

	f, err := os.OpenFile(".dev-up_shell_profile",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println(err)
	}
	defer f.Close()
	if _, err := f.WriteString(str.String()); err != nil {
		log.Println(err)
	}
}

func main() {
	c, err := readConf("./example-conf.yaml")
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		panic(err)
	}

	verifyAndinstallHomeBrew()

	lifecycles := c.Lifecycle
	for i := 0; i < len(lifecycles); i++ {
		currentCycle := lifecycles[i].Install
		verifyAndInstall(currentCycle.App, currentCycle.Version)
		writeToRcFile(currentCycle.Env)
	}

}
