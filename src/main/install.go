package main

import (
	"fmt"
	"os"
	"os/exec"

	"config"
)

func install(config *config.RunConfig) error {
	setEnv(config)

	goarch := os.Getenv("GOARCH")

	if goarch == "amd64" {
		if config.Install.Arch == "64" || config.Install.Arch == "all" {
			err := installArch64(config)
			if err != nil {
				return err
			}
		}

		if config.Install.Arch == "32" || config.Install.Arch == "all" {
			return buildArch32(config)
		}
	} else if goarch == "386" {
		if config.Install.Arch == "32" || config.Install.Arch == "all" {
			return installArch32(config)
		}
	}

	return nil
}

func installArch64(config *config.RunConfig) error {
	fmt.Println("installing 64-bit ......")

	installName := getInstallOutputName(config)
	outputName := getArch64OutputName(config)

	cmd := exec.Command("go", "install", config.Install.Package)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return os.Rename(getBinPath(config)+installName, getBinPath(config)+outputName)
}

func installArch32(config *config.RunConfig) error {
	fmt.Println("installing 32-bit ......")

	installName := getInstallOutputName(config)
	outputName := getArch32OutputName(config)

	cmd := exec.Command("go", "install", config.Install.Package)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	return os.Rename(getBinPath(config)+installName, getBinPath(config)+outputName)
}

func buildArch32(config *config.RunConfig) error {
	fmt.Println("building 32-bit ......")

	outputName := getArch32OutputName(config)

	err := os.Setenv("GOARCH", "386")
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "build",
		"-o", getBinPath(config)+outputName, config.Install.Package)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func getOriginalOutputName(config *config.RunConfig) string {
	name := config.Install.Package
	if config.Install.OutputName != "" {
		name = config.Install.OutputName
	}
	return name
}

func getArch64OutputName(config *config.RunConfig) string {
	name := getOriginalOutputName(config)

	if config.Install.WithArchSuffix {
		name += "_64"
	}

	return addOsSuffix(name)
}

func getArch32OutputName(config *config.RunConfig) string {
	name := getOriginalOutputName(config)

	if config.Install.WithArchSuffix {
		name += "_32"
	}

	return addOsSuffix(name)
}

func getInstallOutputName(config *config.RunConfig) string {
	return addOsSuffix(config.Install.Package)
}
