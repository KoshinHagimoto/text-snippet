package dao

import (
	"database/sql"
	"fmt"
	"os"
	"text-snippet/app/object"

	"github.com/joho/godotenv"
)

type SnippetDAO struct {
	db *sql.DB
}

// 環境変数を読み込む
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

// SnippetDAOを初期化する
func InitSnippetDAO() (*SnippetDAO, error) {
	loadEnv()

	DATABASE_NAME := os.Getenv("DATABASE_NAME")
	DATABASE_USER := os.Getenv("DATABASE_USER")
	DATABASE_PASSWORD := os.Getenv("DATABASE_PASSWORD")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s", DATABASE_USER, DATABASE_PASSWORD, DATABASE_NAME))
	if err != nil {
		return nil, err
	}

	// データベースに接続できるか確認
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// テーブルが存在しない場合は作成する
	createTable := `
	CREATE TABLE IF NOT EXISTS snippets (
		id INT AUTO_INCREMENT PRIMARY KEY,
		content TEXT NOT NULL,
		language VARCHAR(255), 
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP NULL, 
		UNIQUE(unique_string)
	);
	`

	if _, err := db.Exec(createTable); err != nil {
		return nil, err
	}

	return &SnippetDAO{db}, nil
}

// DBをクローズする
func (d SnippetDAO) Finalize() {
	d.db.Close()
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
