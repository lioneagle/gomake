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

	"github.com/lioneagle/gomake/src/config"

	"github.com/lioneagle/goutil/src/chars"
	"github.com/lioneagle/goutil/src/file"
)

func coverage(cfg *config.RunConfig) error {
	setEnv(cfg)

	coverageFileName := getCoverageFileName(cfg)
	tempCoverageFileName := getTempCoverageFileName(cfg)

	file.RemoveExistFile(coverageFileName)

	packages := parsePackages(cfg)
	if len(packages) <= 0 {
		return nil
	}

	err := coverageOnePackage(cfg, packages[0], coverageFileName)
	if err != nil {
		return err
	}

	for i := 1; i < len(packages); i++ {
		err = coverageOnePackage(cfg, packages[i], tempCoverageFileName)
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

	if ok, _ := file.PathOrFileIsExist(coverageFileName); !ok {
		return nil
	}

	if err = showTotalStat(cfg); err != nil {
		return err
	}

	if err = generateHtml(cfg); err != nil {
		return err
	}

	if cfg.Coverage.ShowHtml {
		showHtml(cfg)
	}

	return nil
}

func showTotalStat(cfg *config.RunConfig) error {
	cmd := exec.Command("go", "tool", "cover",
		"-func", getCoverageFileName(cfg))
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

func parsePackages(cfg *config.RunConfig) []string {
	var packages []string

	if cfg.Coverage.Packages == "." {
		packages = getAllPackages(cfg)
	} else {
		packages = strings.Split(cfg.Coverage.Packages, ",")
	}

	//filters := []string{"main", "github", "vendor"}
	filters := []string{"main", "vendor"}

	return chars.FilterReverse(packages, filters)
}

func getAllPackages(cfg *config.RunConfig) []string {
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

		/*if cfg.Coverage.Verbose {
			fmt.Printf(line)
		}*/

	}

	cmd.Wait()

	return ret
}

func coverageOnePackage(cfg *config.RunConfig, packageName, coverageFileName string) error {
	var cmd *exec.Cmd

	if cfg.Coverage.Verbose {
		cmd = exec.Command("go", "test", "-cover",
			"-coverprofile", coverageFileName,
			"-run", cfg.Coverage.Regexp,
			"-v",
			packageName)
	} else {
		cmd = exec.Command("go", "test", "-cover",
			"-coverprofile", coverageFileName,
			"-run", cfg.Coverage.Regexp,
			packageName)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func generateHtml(cfg *config.RunConfig) error {
	cmd := exec.Command("go", "tool", "cover",
		"-html", getCoverageFileName(cfg),
		"-o", getCoverageHtmlFileName(cfg))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func showHtml(cfg *config.RunConfig) error {
	cmd := exec.Command("cmd", "/C", filepath.FromSlash(getCoverageHtmlFileName(cfg)))
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

func getCoverageFileName(cfg *config.RunConfig) string {
	return getTestPath(cfg) + "coverage.out"
}
func getTempCoverageFileName(cfg *config.RunConfig) string {
	return getTestPath(cfg) + "coverage1.out"
}

func getCoverageHtmlFileName(cfg *config.RunConfig) string {
	return getTestPath(cfg) + "coverage.html"
}
