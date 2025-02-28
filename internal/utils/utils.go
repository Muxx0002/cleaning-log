package utils

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"

	"github.com/google/logger"
)

func SearchHistoryFile(browser, root string) ([]string, error) {
	var results []string
	maxDepth := 5
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if strings.Count(filepath.ToSlash(path), "/")-strings.Count(filepath.ToSlash(root), "/") > maxDepth {
			return filepath.SkipDir
		}
		if strings.Contains(strings.ToLower(path), browser) {
			if d.Name() == "History" || d.Name() == "places.sqlite" {
				results = append(results, path)
			}
			return nil
		}
		return nil
	})
	if err != nil {
		logger.Errorf("ошибка поиска файлов, браузер %s. %s", browser, err)
		return nil, err
	}
	return results, nil
}

func KillAllBrowsers() error {
	browsers := []string{
		"chrome.exe",
		"msedge.exe",
		"opera.exe",
		"brave.exe",
		"vivaldi.exe",
		"yandex.exe",
		"epicbrowser.exe",
		"colibri.exe",
		"firefox.exe",
	}
	for _, browser := range browsers {
		cmd := exec.Command("tasklist", "/FI", "IMAGENAME eq "+browser)
		output, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("ошибка проверки процесса: %v", err)
		}
		if !strings.Contains(string(output), browser) {
			fmt.Printf("процесс %s не запущен.\n", browser)
			continue
		}
		cmd = exec.Command("taskkill", "/IM", browser, "/F")
		output, err = cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("неудачное завершение процесса %s: %v\nданные: %s", browser, err, string(output))
		}
		fmt.Printf("процесс завершен: %s\n", browser)
	}
	return nil
}

func DeleteFiles(dir, pattern string) {
	files, err := filepath.Glob(filepath.Join(dir, pattern))
	if err != nil {
		fmt.Println("Ошибка при поиске файлов:", err)
		return
	}

	for _, file := range files {
		err := os.Remove(file)
		if err != nil {
			fmt.Printf("Не удалось удалить файл %s: %v\n", file, err)
		} else {
			fmt.Printf("Удален файл: %s\n", file)
		}
	}
}

func RunCommand(name string, args ...string) {
	cmd := exec.Command(name, args...)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Ошибка выполнения команды %s %v: %v\n", name, args, err)
	}
}

func IsAdmin() bool {
	_, err := os.Open("\\\\.\\PHYSICALDRIVE0")
	return err == nil
}

func ReadKeywordsFromFile(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("ошибка открытия файла: %v", err)
	}
	defer file.Close()
	var keywords []string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line != "" {
			keywords = append(keywords, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("ошибка чтения файла: %v", err)
	}
	if len(keywords) == 0 {
		return nil, fmt.Errorf("файл пустой")
	}
	return keywords, nil
}
