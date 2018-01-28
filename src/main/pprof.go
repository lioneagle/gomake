package main

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/lioneagle/gomake/src/config"

	"github.com/lioneagle/goutil/src/file"
)

func pprof(cfg *config.RunConfig) error {
	setEnv(cfg)

	testFileName := getTestFileName(cfg, cfg.Pprof.Package)
	cpuProfileFileName := getCpuProfileFileName(cfg, cfg.Pprof.Package)

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
		"-nodecount", fmt.Sprintf("%d", cfg.Pprof.NodeCount),
		testFileName, cpuProfileFileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	return cmd.Run()
}
