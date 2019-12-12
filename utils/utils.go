package utils

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func PathExists(path string) bool {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		// doesn't exist
		return false
	}

	return true
}

func Which(cmd string) string {
	p, _ := exec.LookPath(cmd)
	return p
}

func SHA256(src string) string {
	h := sha256.New()
	h.Write([]byte(src))
	return fmt.Sprintf("%x", h.Sum(nil))
}

func GetDir(filePath string) string {
	dir, err := filepath.Abs(filepath.Dir(filePath))
	if err != nil {
		log.Fatal(err)
	}
	return dir
}

func Log(message string) {
	log.Printf("[sealed_secrets_provider] ================= %s\n", message)
}

func ExecuteCmd(cmds ...string) error {
	cmd := strings.Join(cmds[:], " ")
	Log("executing a command: " + cmd)
	out, err := exec.Command("sh", "-c", cmd).Output()
	if err != nil {
		Log("Failed to execute the command: " + cmd)
		return err
	}

	Log(string(out))
	return nil
}

func GetFileName(path string) string {
	f, _ := os.Stat(path)
	return f.Name()
}

func GetFileContent(path string) string {
	c, _ := ioutil.ReadFile(path)
	return string(c)
}
