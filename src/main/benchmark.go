package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"config"

	"github.com/lioneagle/goutil/src/file"
)

func benchmark(config *config.RunConfig) error {
	setEnv(config)

	err := buildTestTempDir(config)
	if err != nil {
		return err
	}

	err = removeBenchmarkFiles(config)
	if err != nil {
		return err
	}

	err = doBenchmark(config)
	if err != nil {
		return err
	}

	return generateTestFile(config)
}

func doBenchmark(config *config.RunConfig) error {
	cpuProfileFileName := getCpuProfileFileName(config, config.Benchmark.Package)
	memProfileFileName := getMemProfileFileName(config, config.Benchmark.Package)

	cmd := exec.Command("go", "test", config.Benchmark.Package,
		"-bench", config.Benchmark.Regexp,
		"-benchtime", fmt.Sprintf("%ds", config.Benchmark.BenchTime),
		"-cpuprofile", cpuProfileFileName,
		"-memprofile", memProfileFileName)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func generateTestFile(config *config.RunConfig) error {
	testFileName := getTestFileName(config, config.Benchmark.Package)

	cmd := exec.Command("go", "test", config.Benchmark.Package,
		"-bench", ".",
		"-c",
		"-o", testFileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func removeBenchmarkFiles(config *config.RunConfig) error {
	filenames := []string{
		getTestFileName(config, config.Benchmark.Package),
		getCpuProfileFileName(config, config.Benchmark.Package),
		getMemProfileFileName(config, config.Benchmark.Package),
	}

	return file.RemoveExistFiles(filenames)
}

func getCpuProfileFileName(config *config.RunConfig, packageName string) string {
	return getTestPath(config) + filepath.Base(packageName) + "_cpu.prof"
}

func getMemProfileFileName(config *config.RunConfig, packageName string) string {
	return getTestPath(config) + filepath.Base(packageName) + "_mem.prof"
}
