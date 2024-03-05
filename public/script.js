let editor;

require.config({ paths: { 'vs': 'https://cdnjs.cloudflare.com/ajax/libs/monaco-editor/0.20.0/min/vs' }});

require(['vs/editor/editor.main'], function() {
    editor = monaco.editor.create(document.getElementById('editor'), {
        language: 'javascript',
        theme: 'vs-dark'
    });
});

function changeLanguage() {
    const languageSelector = document.getElementById('languageSelector');
    const newLanguage = languageSelector.value;

    monaco.editor.setModelLanguage(editor.getModel(), newLanguage);
}

function submitSnippet() {
    const content = editor.getValue();
    const language = document.getElementById('languageSelector').value;

    const expiresAt = new Date();
    expiresAt.setDate(expiresAt.getDate() + 30);
    const expiresAtString = expiresAt.toISOString().split('T')[0];

    fetch('/snippet', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            content,
            language,
            expires_at: expiresAtString
        })
    })
    .then(response => response.json())
    .then(data => {
        console.log('Snippet submitted:', data);
        fetchSnippets();
    })
    .catch(error => console.error('Error:', error));
}

function fetchSnippets() {
    fetch('/snippets')
    .then(response => response.json())
    .then(data => {
        const snippetsDiv = document.getElementById('snippets');
        snippetsDiv.innerHTML = '';
        data.forEach(snippet => {
            const snippetDiv = document.createElement('div');
            snippetDiv.classList.add('snippet');
            snippetDiv.textContent = snippet.content;
            snippetsDiv.appendChild(snippetDiv);
        });
    })
    .catch(error => console.error('Error:', error));
}

document.addEventListener('DOMContentLoaded', () => {
    fetchSnippets();
});

let csrfToken = "";

function fetchCsrfToken() {
    fetch('/csrf-token')
    .then(response => response.json())
    .then(data => {
        csrfToken = data.csrfToken;
    })
    .catch(error => console.error('Error fetching CSRF token:', error));
}

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
    .then(response => response.json())
    .then(data => {
        console.log('Registration successful:', data);
    })
    .catch(error => console.error('Error registering:', error));
}

function login() {
    const username = document.getElementById('login-username').value;
    const password = document.getElementById('login-password').value;

    fetch('/user/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
            'X-CSRF-Token': csrfToken
        },
        body: JSON.stringify({
            username,
            password
        })
    })
    .then(response => response.json())
    .then(data => {
        console.log('Login successful:', data);
        fetchSnippets();
    })
    .catch(error => console.error('Error logging in:', error));
}

document.addEventListener('DOMContentLoaded', () => {
    fetchCsrfToken();
    fetchSnippets();
});
