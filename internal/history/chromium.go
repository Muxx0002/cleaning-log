package history

import (
	"browser/internal/utils"
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/google/logger"
	_ "github.com/mattn/go-sqlite3"
)

func DeleteChromiumHistory(path *string, keywords *[]string) (string, error) {
	db, err := sql.Open("sqlite3", *path)
	if err != nil {
		return "", fmt.Errorf("ошибка открытия БД: %w", err)
	}
	defer db.Close()

	if len(*keywords) == 0 {
		return "список ключевых слов пуст", nil
	}

	var conditions []string
	var params []interface{}
	for _, keyword := range *keywords {
		conditions = append(conditions, "LOWER(title) LIKE ?")
		params = append(params, "%"+strings.ToLower(keyword)+"%")
	}
	if len(conditions) == 0 {
		return "нет условий для удаления", nil
	}

	query := "DELETE FROM urls WHERE " + strings.Join(conditions, " OR ")

	log.Println("SQL Query:", query)
	log.Println("Parameters:", params)

	result, err := db.Exec(query, params...)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения SQL-запроса: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("ошибка получения количества удаленных строк: %w", err)
	}

	return fmt.Sprintf("удалено %d строк", rowsAffected), nil
}

func Chromium(browsers *[]string, root *string, keywords *[]string) ([]string, error) {
	var result []string
	for _, path := range *browsers {
		allpath, err := utils.SearchHistoryFile(path, *root)
		if err != nil {
			logger.Errorf("ошибка поиска файлов браузера: %s", err)
			return nil, err
		}
		for _, browser := range allpath {
			dd, err := DeleteChromiumHistory(&browser, keywords)
			if err != nil {
				logger.Errorf("ошибка удаления истории браузера %s. Ошибка: %s", browser, err)
				return nil, err
			}
			result = append(result, dd+browser)
		}
	}
	return result, nil
}
