package util

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"
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

func RandomString(n int) string {
	rand.Seed(time.Now().UnixNano())
	var letter = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letter[rand.Intn(len(letter))]
	}
	return string(b)
}

func CreateShellScript() string {
	rand_str := RandomString(10)
	rand_script := fmt.Sprintf("%s.sh", rand_str)
	return rand_script
}
