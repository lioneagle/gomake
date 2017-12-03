package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"config"

	"github.com/lioneagle/goutil/src/chars"
	"github.com/lioneagle/goutil/src/file"
)

func coverage(config *config.RunConfig) error {
	setEnv(config)

	coverageFileName := getCoverageFileName(config)
	tempCoverageFileName := getTempCoverageFileName(config)

	file.RemoveExistFile(coverageFileName)

	packages := parsePackages(config)
	if len(packages) <= 0 {
		return nil
	}

	err := coverageOnePackage(config, packages[0], coverageFileName)
	if err != nil {
		return err
	}

	for i := 1; i < len(packages); i++ {
		err = coverageOnePackage(config, packages[i], tempCoverageFileName)
		if err != nil {
			return err
		}

		if ok, _ := file.PathOrFileIsExist(tempCoverageFileName); ok {
			err = mergeCoverageOutput(coverageFileName, tempCoverageFileName)
			if err != nil {
				return err
			}

			file.RemoveExistFile(tempCoverageFileName)
		}
	}

	if ok, _ := file.PathOrFileIsExist(tempCoverageFileName); !ok {
		return nil
	}

	if err = showTotalStat(config); err != nil {
		return err
	}

	if err = generateHtml(config); err != nil {
		return err
	}

	if config.Coverage.ShowHtml {
		showHtml(config)
	}

	return nil
}

func showTotalStat(config *config.RunConfig) error {
	cmd := exec.Command("go", "tool", "cover",
		"-func", getCoverageFileName(config))
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)

	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		if strings.HasPrefix(line, "total") {
			fmt.Println(chars.StringPackSpace(line))
			break
		}
	}

	cmd.Wait()

	return nil
}

func mergeCoverageOutput(destFileName, srcFileName string) error {
	srcData, err := ioutil.ReadFile(srcFileName)
	if err != nil {
		fmt.Printf("ERROR: cannot open file %s\r\n", srcFileName)
		return err
	}

	pos := bytes.Index(srcData, []byte("\n"))
	if pos == -1 {
		return nil
	}

	return file.AppendFile(destFileName, srcData[pos+1:], 0x777)
}

func parsePackages(config *config.RunConfig) []string {
	var packages []string

	if config.Coverage.Packages == "." {
		packages = getAllPackages(config)
	} else {
		packages = strings.Split(config.Coverage.Packages, ",")
	}

	filters := []string{"main", "github", "vendor"}

	return chars.FilterReverse(packages, filters)
}

func getAllPackages(config *config.RunConfig) []string {
	var ret []string
	cmd := exec.Command("go", "list", "./...")
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Println(err)
		return ret
	}

	cmd.Start()

	reader := bufio.NewReader(stdout)

	for {
		line, err := reader.ReadString('\n')
		if err != nil || io.EOF == err {
			break
		}
		ret = append(ret, strings.TrimSpace(line))
	}

	cmd.Wait()

	return ret
}

func coverageOnePackage(config *config.RunConfig, packageName, coverageFileName string) error {
	var cmd *exec.Cmd

	if config.Coverage.Verbose {
		cmd = exec.Command("go", "test", "-cover",
			"-coverprofile", coverageFileName,
			"-run", config.Coverage.Regexp,
			"-v",
			packageName)
	} else {
		cmd = exec.Command("go", "test", "-cover",
			"-coverprofile", coverageFileName,
			"-run", config.Coverage.Regexp,
			packageName)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func generateHtml(config *config.RunConfig) error {
	cmd := exec.Command("go", "tool", "cover",
		"-html", getCoverageFileName(config),
		"-o", getCoverageHtmlFileName(config))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func showHtml(config *config.RunConfig) error {
	cmd := exec.Command("cmd", "/C", filepath.FromSlash(getCoverageHtmlFileName(config)))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func getCoverageFileName(config *config.RunConfig) string {
	return getTestPath(config) + "coverage.out"
}
func getTempCoverageFileName(config *config.RunConfig) string {
	return getTestPath(config) + "coverage1.out"
}

func getCoverageHtmlFileName(config *config.RunConfig) string {
	return getTestPath(config) + "coverage.html"
}
