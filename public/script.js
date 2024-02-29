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
    const language = document.getElementById('languageSelector').value; // 言語の選択値を取得

    // 現在の日付から30日を加算してexpires_atに設定
    const expiresAt = new Date();
    expiresAt.setDate(expiresAt.getDate() + 30); // 30日後の日付を設定
    const expiresAtString = expiresAt.toISOString().split('T')[0]; // YYYY-MM-DD 形式に変換

    fetch('/snippet', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
            content: content,
            language: language, // 言語をリクエストに含める
            expires_at: expiresAtString // 30日後の日付をexpires_atとして設定
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
        snippetsDiv.innerHTML = ''; // Clear previous snippets
        data.forEach(snippet => {
            const snippetDiv = document.createElement('div');
            snippetDiv.classList.add('snippet');
            snippetDiv.textContent = snippet.content; // For simplicity, directly show content
            snippetsDiv.appendChild(snippetDiv);
        });
    })
    .catch(error => console.error('Error:', error));
}

document.addEventListener('DOMContentLoaded', fetchSnippets);
