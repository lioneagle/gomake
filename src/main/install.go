package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/lioneagle/gomake/src/config"

	"github.com/lioneagle/goutil/src/file"
)

func install(cfg *config.RunConfig) error {
	setEnv(cfg)

	if runtime.GOARCH == "amd64" {
		if cfg.Install.Arch == "64" || cfg.Install.Arch == "all" {
			err := installArch64(cfg)
			if err != nil {
				return err
			}
		}

		if cfg.Install.Arch == "32" || cfg.Install.Arch == "all" {
			return buildArch32(cfg)
		}
	} else if runtime.GOARCH == "386" {
		if cfg.Install.Arch == "32" || cfg.Install.Arch == "all" {
			return installArch32(cfg)
		}
	}

	return nil
}

func installArch64(cfg *config.RunConfig) error {
	fmt.Println("installing 64-bit ......")

	installName := getInstallOutputName(cfg)
	outputName := getArch64OutputName(cfg)

	cmd := exec.Command("go", "install", cfg.Install.Package)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	err = os.Rename(getBinPath(cfg)+installName, getBinPath(cfg)+outputName)
	if err != nil {
		return err
	}

	if cfg.Install.CopyToStdGoBin {
		_, err = file.CopyFile(filepath.FromSlash(cfg.OldGobin+"/"+outputName), getBinPath(cfg)+outputName)
		if err != nil {
			return err
		}
	}
	return nil
}

func installArch32(cfg *config.RunConfig) error {
	fmt.Println("installing 32-bit ......")

	installName := getInstallOutputName(cfg)
	outputName := getArch32OutputName(cfg)

	cmd := exec.Command("go", "install", cfg.Install.Package)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}

	err = os.Rename(getBinPath(cfg)+installName, getBinPath(cfg)+outputName)
	if err != nil {
		return err
	}

	if cfg.Install.CopyToStdGoBin {
		_, err = file.CopyFile(filepath.FromSlash(cfg.OldGobin+"/"+outputName), getBinPath(cfg)+outputName)
		if err != nil {
			return err
		}
	}
	return nil
}

func buildArch32(cfg *config.RunConfig) error {
	fmt.Println("building 32-bit ......")

	outputName := getArch32OutputName(cfg)

	err := os.Setenv("GOARCH", "386")
	if err != nil {
		return err
	}

	cmd := exec.Command("go", "build",
		"-o", getBinPath(cfg)+outputName, cfg.Install.Package)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Run()
	if err != nil {
		return err
	}

	if cfg.Install.CopyToStdGoBin {
		_, err = file.CopyFile(filepath.FromSlash(cfg.OldGobin+"/"+outputName), getBinPath(cfg)+outputName)
		if err != nil {
			return err
		}
	}
	return nil
}

func getOriginalOutputName(cfg *config.RunConfig) string {
	name := cfg.Install.Package
	if cfg.Install.OutputName != "" {
		name = cfg.Install.OutputName
	}
	return name
}

func getArch64OutputName(cfg *config.RunConfig) string {
	name := getOriginalOutputName(cfg)

	if cfg.Install.WithArchSuffix {
		name += "_64"
	}

	return addOsSuffix(name)
}

func getArch32OutputName(cfg *config.RunConfig) string {
	name := getOriginalOutputName(cfg)

	if cfg.Install.WithArchSuffix {
		name += "_32"
	}

	return addOsSuffix(name)
}

func getInstallOutputName(cfg *config.RunConfig) string {
	return addOsSuffix(cfg.Install.Package)
}
