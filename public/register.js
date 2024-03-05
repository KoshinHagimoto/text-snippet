let csrfToken = "";

// CSRFトークンを取得する関数
function fetchCsrfToken() {
    fetch('/csrf-token')
    .then(response => response.json())
    .then(data => {
        csrfToken = data.csrfToken;
    })
    .catch(error => console.error('Error fetching CSRF token:', error));
}

// 新規登録を行う関数
function register() {
    const username = document.getElementById('register-username').value;
    const email = document.getElementById('register-email').value;
    const password = document.getElementById('register-password').value;

    fetch('/user/register', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'X-CSRF-Token': csrfToken
        },
        body: JSON.stringify({
            username,
            email,
            password
        })
    })
    .then(response => {
        if (response.ok) {
            console.log('Registration successful');
            // 登録成功後の処理（例: ログインページへのリダイレクト）をここに書く
            window.location.href = 'index.html'; // ホームページにリダイレクト
        } else {
            response.json().then(data => {
                console.error('Registration failed:', data.message);
                // エラーメッセージの表示など、失敗した場合の処理をここに書く
            });
        }
    })
    .catch(error => console.error('Error registering:', error));
}

// ページ読み込み時にCSRFトークンを取得する
document.addEventListener('DOMContentLoaded', () => {
    fetchCsrfToken();
});
