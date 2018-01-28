package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/lioneagle/gomake/src/config"

	"github.com/lioneagle/goutil/src/file"
)

func benchmark(cfg *config.RunConfig) error {
	setEnv(cfg)

	err := buildTestTempDir(cfg)
	if err != nil {
		return err
	}

	err = removeBenchmarkFiles(cfg)
	if err != nil {
		return err
	}

	err = doBenchmark(cfg)
	if err != nil {
		return err
	}

	err = generateTestFile(cfg)
	if err != nil {
		return err
	}

	if ok, _ := file.PathOrFileIsExist(getCpuProfileFileName(cfg, cfg.Benchmark.Package)); !ok {
		return nil
	}

	if cfg.Benchmark.GoTorch {
		err = generateTorchFile(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func doBenchmark(cfg *config.RunConfig) error {
	cpuProfileFileName := getCpuProfileFileName(cfg, cfg.Benchmark.Package)
	memProfileFileName := getMemProfileFileName(cfg, cfg.Benchmark.Package)

	cmd := exec.Command("go", "test", cfg.Benchmark.Package,
		"-bench", cfg.Benchmark.Regexp,
		"-benchtime", fmt.Sprintf("%ds", cfg.Benchmark.BenchTime),
		"-cpuprofile", cpuProfileFileName,
		"-memprofile", memProfileFileName)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func generateTorchFile(cfg *config.RunConfig) error {
	torchFileName := getTorchFileName(cfg, cfg.Benchmark.Package)
	cpuProfileFileName := getCpuProfileFileName(cfg, cfg.Benchmark.Package)

	cmd := exec.Command("go-torch", cpuProfileFileName,
		"-f", torchFileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func generateTestFile(cfg *config.RunConfig) error {
	testFileName := getTestFileName(cfg, cfg.Benchmark.Package)

	cmd := exec.Command("go", "test", cfg.Benchmark.Package,
		"-bench", ".",
		"-c",
		"-o", testFileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func removeBenchmarkFiles(cfg *config.RunConfig) error {
	filenames := []string{
		getTestFileName(cfg, cfg.Benchmark.Package),
		getCpuProfileFileName(cfg, cfg.Benchmark.Package),
		getMemProfileFileName(cfg, cfg.Benchmark.Package),
	}

	return file.RemoveExistFiles(filenames)
}

func getCpuProfileFileName(cfg *config.RunConfig, packageName string) string {
	return getTestPath(cfg) + filepath.Base(packageName) + "_cpu.prof"
}

func getMemProfileFileName(cfg *config.RunConfig, packageName string) string {
	return getTestPath(cfg) + filepath.Base(packageName) + "_mem.prof"
}

func getTorchFileName(cfg *config.RunConfig, packageName string) string {
	if cfg.Benchmark.Regexp == "." {
		return fmt.Sprintf("%s%s_benchtime%d_cpu.svg", getTestPath(cfg), filepath.Base(packageName), cfg.Benchmark.BenchTime)
	}
	regexp := strings.Replace(cfg.Benchmark.Regexp, "*", "", -1)
	return fmt.Sprintf("%s%s_%s_benchtime%d_cpu.svg", getTestPath(cfg), filepath.Base(packageName), filepath.Base(regexp), cfg.Benchmark.BenchTime)
}
