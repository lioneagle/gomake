package main

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/lioneagle/gomake/src/config"

	"github.com/lioneagle/goutil/src/file"
)

func setEnv(cfg *config.RunConfig) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	oldPath := os.Getenv("GOPATH")
	newPath := oldPath
	if newPath == "" {
		newPath = dir
	} else {
		if runtime.GOOS == "windows" {
			newPath = fmt.Sprintf("%s;%s", oldPath, dir)
		} else {
			newPath = fmt.Sprintf("%s:%s", oldPath, dir)
		}
	}

	err = os.Setenv("GOPATH", newPath)
	if err != nil {
		return err
	}

	cfg.OldGobin = os.Getenv("GOBIN")

	return os.Setenv("GOBIN", filepath.FromSlash(dir+"/bin"))
}

func addOsSuffix(name string) string {
	if runtime.GOOS == "windows" {
		return name + ".exe"
	}
	return name
}

func getTestFileName(cfg *config.RunConfig, packageName string) string {
	name := getTestPath(cfg) + filepath.Base(packageName) + ".test"
	return addOsSuffix(name)
}

func getTestPath(cfg *config.RunConfig) string {
	return "./test_temp/"
}

func getBinPath(cfg *config.RunConfig) string {
	return "./bin/"
}

func buildTestTempDir(cfg *config.RunConfig) error {
	testTempDir := getTestPath(cfg)
	ok, _ := file.PathOrFileIsExist(testTempDir)
	if !ok {
		return os.Mkdir(testTempDir, os.ModeDir)
	}
	return nil
}
