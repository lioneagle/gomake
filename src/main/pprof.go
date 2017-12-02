package main

import (
	"fmt"
	"os"
	"os/exec"

	"config"

	"github.com/lioneagle/goutil/src/file"
)

func pprof(config *config.RunConfig) error {
	setEnv(config)

	testFileName := getTestFileName(config, config.Pprof.Package)
	cpuProfileFileName := getCpuProfileFileName(config, config.Pprof.Package)

	ok, _ := file.PathOrFileIsExist(testFileName)
	if !ok {
		return nil
	}

	ok, _ = file.PathOrFileIsExist(cpuProfileFileName)
	if !ok {
		return nil
	}

	fmt.Println("testFileName =", testFileName)
	fmt.Println("cpuProfileFileName =", cpuProfileFileName)

	cmd := exec.Command("go", "tool", "pprof",
		"-nodecount", fmt.Sprintf("%d", config.Pprof.NodeCount),
		testFileName, cpuProfileFileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
