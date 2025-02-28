package history

import (
	"browser/internal/utils"
	"database/sql"

	"github.com/google/logger"
)

func DeleteGeckoHistory(path *string, keywords *[]string) (string, error) {
	db, err := sql.Open("sqlite3", *path)
	if err != nil {
		return "ошибка", err
	}
	defer db.Close()

	query := "DELETE FROM moz_places WHERE "
	for i, keyword := range *keywords {
		if i > 0 {
			query += " OR "
		}
		// Применяем поиск к url и title для каждого ключевого слова
		query += "(url LIKE '%" + keyword + "%' OR title LIKE '%" + keyword + "%')"
	}

	_, err = db.Exec(query)
	if err != nil {
		return "ошибка при запросе", err
	}

	return "удалено", nil
}

func Gecko(browsers *[]string, root *string, keywords *[]string) ([]string, error) {
	var result []string

	for _, path := range *browsers {
		allpath, err := utils.SearchHistoryFile(path, *root)
		if err != nil {
			logger.Errorf("ошибка поиска файлов браузера: %s", err)
			return nil, err
		}

		for _, browser := range allpath {
			dd, err := DeleteGeckoHistory(&browser, keywords)
			if err != nil {
				logger.Errorf("ошибка удаления истории браузера %s. Ошибка: %s", browser, err)
				return nil, err
			}
			result = append(result, dd)
		}
	}

	return result, nil
}
