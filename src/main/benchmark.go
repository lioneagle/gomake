package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

	err = generateTestFile(config)
	if err != nil {
		return err
	}

	if ok, _ := file.PathOrFileIsExist(getCpuProfileFileName(config, config.Benchmark.Package)); !ok {
		return nil
	}

	if config.Benchmark.GoTorch {
		err = generateTorchFile(config)
		if err != nil {
			return err
		}
	}
	return nil
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

func generateTorchFile(config *config.RunConfig) error {
	torchFileName := getTorchFileName(config, config.Benchmark.Package)
	cpuProfileFileName := getCpuProfileFileName(config, config.Benchmark.Package)

	cmd := exec.Command("go-torch", cpuProfileFileName,
		"-f", torchFileName)
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

func getTorchFileName(config *config.RunConfig, packageName string) string {
	if config.Benchmark.Regexp == "." {
		return fmt.Sprintf("%s%s_benchtime%d_cpu.svg", getTestPath(config), filepath.Base(packageName), config.Benchmark.BenchTime)
	}
	regexp := strings.Replace(config.Benchmark.Regexp, "*", "", -1)
	return fmt.Sprintf("%s%s_%s_benchtime%d_cpu.svg", getTestPath(config), filepath.Base(packageName), filepath.Base(regexp), config.Benchmark.BenchTime)
}
