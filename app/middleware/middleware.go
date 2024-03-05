package middleware

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/http"
	"os"
	"strconv"
	"time"
)

// CORSMiddleware wraps handler with CORS header
func CORSMiddleware(handler http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		if r.Method == "OPTIONS" {
			w.Header().Set("Access-Control-Allow-Methods", "*")
			w.Header().Set("Access-Control-Allow-Headers", "*")
			return
		}
		handler(w, r)
	}
}

// 署名付きURLが有効か確認する
func VerifySignature(r *http.Request) bool {
	query := r.URL.Query()
	signature := query.Get("signature")
	if signature == "" {
		return false
	}

	// クエリーから署名を削除し、署名されたコンテンツを再作成する。
	query.Del("signature")
	signedContent := query.Encode()

	secretKey := os.Getenv("SIGNED_URL_SECRET")

	hmacHash := hmac.New(sha256.New, []byte(secretKey))
	hmacHash.Write([]byte(signedContent))
	expectedSignature := hex.EncodeToString(hmacHash.Sum(nil))

	// 署名が一致するか確認
	if !hmac.Equal([]byte(signature), []byte(expectedSignature)) {
		return false
	}

	// 期限を確認
	expiration := query.Get("expiration")
	if expiration == "" {
		return false
	}

	expirationTime, err := strconv.ParseInt(expiration, 10, 64)
	if err != nil {
		return false
	}

	// 現在の時間が期限を過ぎているか確認
	if time.Now().Unix() > expirationTime {
		return false
	}

	return true
}
