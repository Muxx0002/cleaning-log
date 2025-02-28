package logcleaner

import (
	"browser/internal/utils"
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func ClearRegistryBasic() {
	fmt.Println("Очистка основных логов реестра...")
	utils.RunCommand("reg", "delete", `HKEY_CURRENT_USER\Software\Microsoft\Windows\CurrentVersion\Explorer\RunMRU`, "/va", "/f")
	utils.RunCommand("reg", "delete", `HKEY_CURRENT_USER\Software\Microsoft\Windows\Shell\BagMRU`, "/f")
	utils.RunCommand("reg", "delete", `HKEY_CURRENT_USER\Software\Microsoft\Windows\Shell\Bags`, "/f")
	fmt.Println("Основные логи реестра очищены.")
}

func ClearPrefetch() {
	fmt.Println("Очистка файлов Prefetch...")
	utils.DeleteFiles(filepath.Join(os.Getenv("SystemRoot"), "Prefetch"), "*.pf")
	utils.DeleteFiles(filepath.Join(os.Getenv("SystemRoot"), "Prefetch"), "*.db")
	fmt.Println("Файлы Prefetch очищены.")
}

func ClearMinidump() {
	fmt.Println("Очистка файлов Minidump...")
	utils.DeleteFiles(filepath.Join(os.Getenv("SystemRoot"), "Minidump"), "*.*")
	fmt.Println("Файлы Minidump очищены.")
}

func ClearWindowsLogs() {
	fmt.Println("Очистка журналов Windows...")

	out, err := exec.Command("wevtutil", "el").Output()
	if err != nil {
		fmt.Println("Ошибка при получении списка журналов:", err)
		return
	}
	logs := string(out)
	scanner := bufio.NewScanner(strings.NewReader(logs))

	for scanner.Scan() {
		log := scanner.Text()
		fmt.Printf("Очистка журнала %s...\n", log)
		utils.RunCommand("wevtutil", "cl", log)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Ошибка при обработке списка журналов:", err)
	}

	fmt.Println("Журналы Windows очищены.")
}
