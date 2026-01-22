let allCharacters = [];
let allStyles = [];

async function loadCharacters() {
    try {
        const response = await fetch('/api/characters');
        const data = await response.json();
        allCharacters = data || [];
    } catch (error) {
        allCharacters = [];
    }
}

async function loadStyles() {
    try {
        const response = await fetch('/api/combat-styles');
        const data = await response.json();
        allStyles = data || [];
    } catch (error) {
        allStyles = [];
    }
}

function levenshteinDistance(str1, str2) {
    const len1 = str1.length;
    const len2 = str2.length;
    const matrix = [];

    for (let i = 0; i <= len1; i++) {
        matrix[i] = [i];
    }

    for (let j = 0; j <= len2; j++) {
        matrix[0][j] = j;
    }

    for (let i = 1; i <= len1; i++) {
        for (let j = 1; j <= len2; j++) {
            if (str1.charAt(i - 1) === str2.charAt(j - 1)) {
                matrix[i][j] = matrix[i - 1][j - 1];
            } else {
                matrix[i][j] = Math.min(
                    matrix[i - 1][j - 1] + 1,
                    matrix[i][j - 1] + 1,
                    matrix[i - 1][j] + 1
                );
            }
        }
    }

    return matrix[len1][len2];
}

function fuzzyMatch(text, query) {
    const textLower = text.toLowerCase();
    const queryLower = query.toLowerCase();

    if (textLower.includes(queryLower)) {
        return 0;
    }

    const distance = levenshteinDistance(textLower, queryLower);
    const maxLength = Math.max(textLower.length, queryLower.length);

    if (distance / maxLength < 0.3) {
        return distance;
    }

    return Infinity;
}

function search(query) {
    if (!query || query.length < 1) return [];

    const q = query.toLowerCase();
    const results = [];
    const isNumber = /^\d+$/.test(query);

    if (isNumber) {
        const id = parseInt(query);

        const char = allCharacters.find(c => c.id === id);
        if (char) {
            results.push({
                type: 'char',
                name: char.name,
                img: char.img,
                url: `/characters/${char.id}`,
                score: 0
            });
        }

        const style = allStyles.find(s => s.id === id);
        if (style) {
            results.push({
                type: 'style',
                name: style.name,
                url: `/combat-styles/${style.id}`,
                score: 0
            });
        }
    }

    allCharacters.forEach(char => {
        if (char.name) {
            const score = fuzzyMatch(char.name, q);
            if (score !== Infinity && !results.find(r => r.type === 'char' && r.name === char.name)) {
                results.push({
                    type: 'char',
                    name: char.name,
                    img: char.img,
                    url: `/characters/${char.id}`,
                    score: score
                });
            }
        }
    });

    allStyles.forEach(style => {
        if (style.name) {
            const score = fuzzyMatch(style.name, q);
            if (score !== Infinity && !results.find(r => r.type === 'style' && r.name === style.name)) {
                results.push({
                    type: 'style',
                    name: style.name,
                    url: `/combat-styles/${style.id}`,
                    score: score
                });
            }
        }
    });

    allCharacters.forEach(char => {
        if (char.quote && char.quote.toLowerCase().includes(q) && !results.find(r => r.name === char.name)) {
            results.push({
                type: 'quote',
                name: char.name,
                quote: char.quote,
                img: char.img,
                url: `/characters/${char.id}`,
                score: 100
            });
        }
    });

    results.sort((a, b) => a.score - b.score);
    return results.slice(0, 8);
}

function showResults(results) {
    const div = document.getElementById('searchResults');
    if (!div) return;

    if (results.length === 0) {
        div.innerHTML = '<div class="search-no-results">Aucun résultat</div>';
        div.classList.add('active');
        return;
    }

    let html = '';
    results.forEach(r => {
        if (r.type === 'char') {
            html += `
                <a href="${r.url}" class="search-result-item">
                    <img src="${r.img}" class="search-result-img" alt="${r.name}">
                    <div>
                        <div class="search-result-name">${r.name}</div>
                        <div class="search-result-race">Personnage</div>
                    </div>
                </a>
            `;
        } else if (r.type === 'style') {
            html += `
                <a href="${r.url}" class="search-result-item">
                    <img src="/static/images/combatstyle.webp" class="search-result-img" alt="${r.name}">
                    <div>
                        <div class="search-result-name">${r.name}</div>
                        <div class="search-result-race">Style de combat</div>
                    </div>
                </a>
            `;
        } else if (r.type === 'quote') {
            const short = r.quote.length > 50 ? r.quote.substring(0, 50) + '...' : r.quote;
            html += `
                <a href="${r.url}" class="search-result-item">
                    <img src="${r.img}" class="search-result-img" alt="${r.name}">
                    <div>
                        <div class="search-result-name">"${short}"</div>
                        <div class="search-result-race">Citation - ${r.name}</div>
                    </div>
                </a>
            `;
        }
    });

    div.innerHTML = html;
    div.classList.add('active');
}

function hideResults() {
    const div = document.getElementById('searchResults');
    if (div) {
        setTimeout(() => div.classList.remove('active'), 200);
    }
}

document.addEventListener('DOMContentLoaded', async () => {
    const input = document.getElementById('searchInput');
    const resultsDiv = document.getElementById('searchResults');

    if (!input) return;

    await loadCharacters();
    await loadStyles();

    input.addEventListener('input', async (e) => {
        const query = e.target.value;
        if (query.length < 1) {
            hideResults();
            return;
        }
        const results = search(query);
        showResults(results);
    });

    input.addEventListener('focus', (e) => {
        if (e.target.value.length >= 1) {
            const results = search(e.target.value);
            showResults(results);
        }
    });

    input.addEventListener('blur', hideResults);

    input.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            const results = search(input.value);
            if (results.length > 0) {
                window.location.href = results[0].url;
            }
        }
    });

    document.addEventListener('click', (e) => {
        if (!input.contains(e.target) && resultsDiv && !resultsDiv.contains(e.target)) {
            hideResults();
        }
    });
});

const css = document.createElement('style');
css.textContent = `
    .search-result-race {
        font-size: 12px;
        color: #666;
        margin-top: 2px;
    }
`;
document.head.appendChild(css);
