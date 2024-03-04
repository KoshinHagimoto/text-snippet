package dao

import (
	"database/sql"
	"text-snippet/app/object"

	_ "github.com/go-sql-driver/mysql"
)

type SnippetDAO struct {
	db *sql.DB
}

func NewSnippetDAO(db *sql.DB) *SnippetDAO {
	return &SnippetDAO{db: db}
}

// SnippetをDBに保存する
func (d SnippetDAO) Save(snippet object.Snippet) (int, error) {
	result, err := d.db.Exec("INSERT INTO snippets (content, language, expires_at) VALUES (?, ?, ?)", snippet.Content, snippet.Language, snippet.ExpiresAt)
	if err != nil {
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

// すべてのSnippetを取得する
func (d SnippetDAO) GetAll() ([]object.Snippet, error) {
	snippets := []object.Snippet{}

	rows, err := d.db.Query("SELECT id, content, language, created_at, expires_at FROM snippets")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s object.Snippet
		err := rows.Scan(&s.ID, &s.Content, &s.Language, &s.CreatedAt, &s.ExpiresAt)
		if err != nil {
			return nil, err
		}
		snippets = append(snippets, s)
	}
	return snippets, nil
}
