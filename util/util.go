package util

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func ReadFile(path string) ([]string, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	file, ferr := os.Open(path)
	if ferr != nil {
		return nil, ferr
	}

	data := []string{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		data = append(data, scanner.Text())
	}

	file.Close()
	return data, nil
}

func ProcessShell(data []string, cwd string) ([]string, error) {
	pdata := make([]string, len(data)+1)
	if len(data) > 0 && strings.HasPrefix(data[0], "#!") {
		copy(pdata, data[:1])
		pdata[1] = "source ~/.bashrc"
		copy(pdata[2:], data[1:])
	}

	for idx, line := range pdata {
		if strings.HasPrefix(line, "$bin/canu") {
			items := []string{}
			items = append(items, fmt.Sprintf("%s/canu", cwd))
			newline := fmt.Sprintf("\"%s\"", strings.TrimSpace(strings.TrimPrefix(line, "$bin/canu")))
			items = append(items, newline)
			pdata[idx] = strings.Join(items, " ")
		}
	}
	return pdata, nil
}

func WriteToFile(path string, data []string) error {
	ofile, _ := os.Create(path)
	defer ofile.Close()
	writer := bufio.NewWriter(ofile)
	for _, line := range data {
		fmt.Fprintln(writer, line)
	}
	writer.Flush()
	return nil
}

func Command(cmdStr string, arg ...string) error {
	cmd := exec.Command(cmdStr, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		return err
	}
	return nil
}
