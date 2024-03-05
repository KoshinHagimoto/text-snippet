package config

import (
	"crypto/hmac"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"net/smtp"
	"os"
	"time"

	"github.com/joho/godotenv"
)

// 環境変数を読み込む
func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
}

// SnippetDAOを初期化する
func InitDAO() (*sql.DB, error) {
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
	createSnippetTable := `
	CREATE TABLE IF NOT EXISTS snippets (
		id INT AUTO_INCREMENT PRIMARY KEY,
		content TEXT NOT NULL,
		language VARCHAR(255), 
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
		expires_at TIMESTAMP NULL
	);
	`

	if _, err := db.Exec(createSnippetTable); err != nil {
		return nil, err
	}

	// Users テーブルが存在しない場合は作成する
	createUsersTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INT AUTO_INCREMENT PRIMARY KEY,
        username VARCHAR(255) NOT NULL UNIQUE,
        email VARCHAR(255) NOT NULL UNIQUE,
        password_hash VARCHAR(255) NOT NULL,
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
        email_verified BOOLEAN DEFAULT FALSE
    );
    `
	if _, err := db.Exec(createUsersTable); err != nil {
		return nil, err
	}

	return db, nil
}

// 署名付きURLを生成する
func GenerateSignedURL(userID int, userEmail string) string {
	secretKey := os.Getenv("SIGNED_URL_SECRET")
	expiration := time.Now().Add(30 * time.Minute).Unix()
	userHash := hmac.New(sha256.New, []byte(secretKey))
	userHash.Write([]byte(userEmail))
	userHashStr := hex.EncodeToString(userHash.Sum(nil))

	unsignedURL := fmt.Sprintf("/verify/email?id=%d&user=%s&expiration=%d", userID, userHashStr, expiration)
	signature := hmac.New(sha256.New, []byte(secretKey))
	signature.Write([]byte(unsignedURL))
	signatureStr := hex.EncodeToString(signature.Sum(nil))

	return fmt.Sprintf("%s&signature=%s", unsignedURL, signatureStr)
}

// メールを送信する
func SendVerificationEmail(userEmail, verificationURL string) {
	from := os.Getenv("SMTP_FROM")
	password := os.Getenv("SMTP_PASSWORD")
	to := []string{userEmail}
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := os.Getenv("SMTP_PORT")

	message := []byte("To: " + userEmail + "\r\n" +
		"Subject: Email Verification\r\n\r\n" +
		"Click the link below to verify your email address:\r\n" +
		verificationURL + "\r\n")

	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	if err != nil {
		fmt.Println("error sending email: ", err)
		return
	}

	fmt.Println("email sent successfully")
}
