// ========== RECHERCHE SIMPLIFIÉE ==========

let allCharacters = [];
let allStyles = [];

// Charger les personnages
async function loadCharacters() {
    try {
        const response = await fetch('https://www.demonslayer-api.com/api/v1/characters?limit=100');
        const data = await response.json();
        allCharacters = data.content || [];
        console.log('✅ Personnages chargés:', allCharacters.length);
    } catch (error) {
        console.error('❌ Erreur chargement personnages:', error);
        allCharacters = [];
    }
}

// Charger les styles de combat depuis l'API
async function loadStyles() {
    try {
        const response = await fetch('https://www.demonslayer-api.com/api/v1/combat-styles?limit=100');
        const data = await response.json();
        allStyles = data.content || [];
        console.log('✅ Styles de combat chargés:', allStyles.length);
    } catch (error) {
        console.error('❌ Erreur chargement styles:', error);
        allStyles = [];
    }
}

// Rechercher
function search(query) {
    if (!query || query.length < 1) return [];
    
    const q = query.toLowerCase();
    const results = [];
    
    // Est-ce un ID ?
    const isNumber = /^\d+$/.test(query);
    
    if (isNumber) {
        const id = parseInt(query);
        
        // Chercher personnage par ID
        const char = allCharacters.find(c => c.id === id);
        if (char) {
            results.push({
                type: 'char',
                name: char.name,
                img: char.img,
                url: `/characters/${char.id}`
            });
        }
        
        // Chercher style par ID
        const style = allStyles.find(s => s.id === id);
        if (style) {
            results.push({
                type: 'style',
                name: style.name,
                url: `/combat-styles/${encodeURIComponent(style.name)}`
            });
        }
    }
    
    // Chercher par nom dans personnages
    allCharacters.forEach(char => {
        if (char.name && char.name.toLowerCase().includes(q) && results.length < 5) {
            results.push({
                type: 'char',
                name: char.name,
                img: char.img,
                url: `/characters/${char.id}`
            });
        }
    });
    
    // Chercher par nom dans styles
    allStyles.forEach(style => {
        if (style.name && style.name.toLowerCase().includes(q) && results.length < 8) {
            results.push({
                type: 'style',
                name: style.name,
                url: `/combat-styles/${encodeURIComponent(style.name)}`
            });
        }
    });
    
    // Chercher dans citations
    allCharacters.forEach(char => {
        if (char.quote && char.quote.toLowerCase().includes(q) && results.length < 8) {
            results.push({
                type: 'quote',
                name: char.name,
                quote: char.quote,
                img: char.img,
                url: `/characters/${char.id}`
            });
        }
    });
    
    return results.slice(0, 8);
}

// Afficher résultats
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
                    <div style="width:50px;height:50px;border-radius:50%;background:rgba(58,78,68,0.8);display:flex;align-items:center;justify-content:center;font-size:24px;">⚔️</div>
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

// Cacher résultats
function hideResults() {
    const div = document.getElementById('searchResults');
    if (div) {
        setTimeout(() => div.classList.remove('active'), 200);
    }
}

// Init
document.addEventListener('DOMContentLoaded', async () => {
    const input = document.getElementById('searchInput');
    const resultsDiv = document.getElementById('searchResults');
    
    if (!input) {
        console.log('Pas de recherche sur cette page');
        return;
    }
    
    console.log('🔍 Recherche activée');
    
    // Charger données
    await loadCharacters();
    await loadStyles();
    
    // Input
    input.addEventListener('input', async (e) => {
        const query = e.target.value;
        if (query.length < 1) {
            hideResults();
            return;
        }
        const results = search(query);
        showResults(results);
    });
    
    // Focus
    input.addEventListener('focus', (e) => {
        if (e.target.value.length >= 1) {
            const results = search(e.target.value);
            showResults(results);
        }
    });
    
    // Blur
    input.addEventListener('blur', hideResults);
    
    // Enter
    input.addEventListener('keydown', (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            const results = search(input.value);
            if (results.length > 0) {
                window.location.href = results[0].url;
            }
        }
    });
    
    // Click outside
    document.addEventListener('click', (e) => {
        if (!input.contains(e.target) && resultsDiv && !resultsDiv.contains(e.target)) {
            hideResults();
        }
    });
});

// Style
const css = document.createElement('style');
css.textContent = `
    .search-result-race {
        font-size: 12px;
        color: #666;
        margin-top: 2px;
    }
`;
document.head.appendChild(css);

console.log('🔍 Module de recherche chargé');