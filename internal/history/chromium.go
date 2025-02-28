package history

import (
	"browser/internal/utils"
	"database/sql"

	"github.com/google/logger"
	_ "github.com/mattn/go-sqlite3"
)

func DeleteChromiumHistory(path *string, keywords *[]string) (string, error) {
	db, err := sql.Open("sqlite3", *path)
	if err != nil {
		return "", err
	}
	defer db.Close()

	query := "DELETE FROM urls WHERE "
	for i, keyword := range *keywords {
		if i > 0 {
			query += " OR "
		}
		query += "url LIKE '%" + keyword + "%'"
	}

	_, err = db.Exec(query)
	if err != nil {
		return "", err
	}

	return "удалено", nil
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
			dd, err := DeleteChromiumHistory(&browser, keywords) // Используем переданный параметр keywords вместо wordlists
			if err != nil {
				logger.Errorf("ошибка удаления истории браузера %s. Ошибка: %s", browser, err)
				return nil, err
			}
			result = append(result, dd+browser)
		}
	}

	return result, nil
}
