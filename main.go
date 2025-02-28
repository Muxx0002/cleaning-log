package main

import (
	"browser/internal/history"
	"browser/internal/logcleaner"
	"browser/internal/utils"
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/logger"
)

var ChromiumBrowsers = []string{
	"chrome",
	"edge",
	"opera",
	"brave",
	"vivaldi",
	"yandex",
	"epic Privacy",
	"colibri",
}

var GeckoBrowsers = []string{
	"firefox",
	"Basilisk",
	"icecat",
	"k-meleon",
	"pale moon",
	"tor",
	"Waterfox",
}

var wordlists []string

var (
	LOCALAPPDATA   = os.Getenv("LOCALAPPDATA")
	ROAMINGAPPDATA = os.Getenv("APPDATA")
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Необходимо запустить программу от имени администратора.")
		fmt.Println("1 - Очистка истории браузеров")
		fmt.Println("2 - Очистка основных логов в реестре")
		fmt.Println("3 - Очистка файлов Prefetch и Minidump")
		fmt.Println("4 - Очистка журналов Windows")
		fmt.Println("5 - Очистка всех логов в реестре, файлов Prefetch и Minidump")
		fmt.Println("6 - Очистка всех логов, файлов Prefetch, журналов Windows и истории браузеров")
		fmt.Println("Нажмите ENTER для выхода")

		fmt.Print("Ваш выбор: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)
		switch choice {
		case "1":
			if err := DeleteHistoryBrowser(); err != nil {
				fmt.Printf("Ошибка при очистке истории браузеров: %v\n", err)
			}
		case "2":
			logcleaner.ClearRegistryBasic()
		case "3":
			logcleaner.ClearPrefetch()
			logcleaner.ClearMinidump()
		case "4":
			logcleaner.ClearWindowsLogs()
		case "5":
			logcleaner.ClearRegistryBasic()
			logcleaner.ClearPrefetch()
			logcleaner.ClearMinidump()
		case "6":
			if err := DeleteHistoryBrowser(); err != nil {
				fmt.Printf("Ошибка при очистке истории браузеров: %v\n", err)
			}
			logcleaner.ClearRegistryBasic()
			logcleaner.ClearPrefetch()
			logcleaner.ClearMinidump()
			logcleaner.ClearWindowsLogs()
		case "":
			fmt.Println("Выход из программы.")
			return
		default:
			fmt.Println("Некорректный выбор. Попробуйте снова.")
		}
	}
}

func DeleteHistoryBrowser() error {
	var err error
	wordlists, err = utils.ReadKeywordsFromFile("words.txt")
	if err != nil {
		return fmt.Errorf("ошибка чтения файла ключевых слов: %v", err)
	}

	if len(wordlists) == 0 {
		fmt.Println("Список ключевых слов пуст")
	} else {
		fmt.Println("Используемые ключевые слова:", wordlists)
	}

	err = utils.KillAllBrowsers()
	if err != nil {
		logger.Errorf("ошибка закрытия браузеров: %v", err)
		return err
	}

	chromium, err := history.Chromium(&ChromiumBrowsers, &LOCALAPPDATA, &wordlists)
	if err != nil {
		logger.Errorf("ошибка очистки браузеров Chromium: %v", err)
	} else {
		fmt.Println("Результат очистки Chromium-браузеров:")
		for _, res := range chromium {
			fmt.Println(res)
		}
	}
	gecko, err := history.Gecko(&GeckoBrowsers, &ROAMINGAPPDATA, &wordlists)
	if err != nil {
		logger.Errorf("ошибка очистки браузеров Gecko: %v", err)
		if len(chromium) == 0 {
			return err
		}
	} else {
		fmt.Println("Результат очистки Gecko-браузеров:")
		for _, res := range gecko {
			fmt.Println(res)
		}
	}
	if len(chromium) == 0 && len(gecko) == 0 {
		return errors.New("не удалось очистить ни один браузер")
	}
	return nil
}
