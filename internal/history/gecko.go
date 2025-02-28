package history

import (
	"browser/internal/utils"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/google/logger"
)

func DeleteGeckoHistory(path *string, keywords *[]string) (string, error) {
	db, err := sql.Open("sqlite3", *path)
	if err != nil {
		return "", fmt.Errorf("ошибка открытия базы данных: %w", err)
	}
	defer db.Close()

	if len(*keywords) == 0 {
		return "нет ключевых слов для удаления", nil
	}

	var conditions []string
	var params []interface{}

	for _, keyword := range *keywords {
		conditions = append(conditions, "(UPPER(url) LIKE ? OR UPPER(title) LIKE ?)")
		likePattern := "%" + strings.ToUpper(keyword) + "%"
		params = append(params, likePattern, likePattern)
	}

	if len(conditions) == 0 {
		return "нет условий для удаления", nil
	}

	query := "DELETE FROM moz_places WHERE " + strings.Join(conditions, " OR ")

	result, err := db.Exec(query, params...)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения SQL-запроса: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return "", fmt.Errorf("ошибка получения количества удаленных строк: %w", err)
	}

	return fmt.Sprintf("удалено %d записей", rowsAffected), nil
}

func Gecko(browsers *[]string, root *string, keywords *[]string) ([]string, error) {
	var result []string
	if len(*browsers) == 0 {
		return nil, errors.New("список браузеров пуст")
	}

	for _, path := range *browsers {
		allpath, err := utils.SearchHistoryFile(path, *root)
		if err != nil {
			logger.Errorf("ошибка поиска файлов браузера: %s", err)
			continue
		}

		for _, browser := range allpath {
			dd, err := DeleteGeckoHistory(&browser, keywords)
			if err != nil {
				logger.Errorf("ошибка удаления истории браузера %s. Ошибка: %s", browser, err)
				continue
			}
			result = append(result, fmt.Sprintf("%s для %s", dd, browser))
		}
	}

	if len(result) == 0 {
		return nil, errors.New("не удалось обработать ни один файл истории")
	}

	return result, nil
}
