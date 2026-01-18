// ========== RECHERCHE DE PERSONNAGES ==========

// Cache pour stocker les personnages (évite de refaire l'API call)
let cachedCharacters = null;

// Charger tous les personnages au chargement de la page
async function loadAllCharacters() {
    if (cachedCharacters) {
        return cachedCharacters; // Utiliser le cache si disponible
    }
    
    try {
        const response = await fetch('https://www.demonslayer-api.com/api/v1/characters?limit=100');
        const data = await response.json();
        cachedCharacters = data.content || [];
        console.log('✅ Personnages chargés:', cachedCharacters.length);
        return cachedCharacters;
    } catch (error) {
        console.error('❌ Erreur chargement personnages:', error);
        return [];
    }
}

// Fonction de recherche
async function searchCharacters(query) {
    if (!query || query.length < 2) {
        return []; // Minimum 2 caractères
    }
    
    const characters = await loadAllCharacters();
    const lowerQuery = query.toLowerCase();
    
    // Filtrer les personnages dont le nom contient la recherche
    const results = characters.filter(char => 
        char.name && char.name.toLowerCase().includes(lowerQuery)
    );
    
    return results.slice(0, 5); // Maximum 5 résultats
}

// Afficher les résultats de recherche
function displaySearchResults(results) {
    const searchResults = document.getElementById('searchResults');
    
    if (!searchResults) return;
    
    // Si aucun résultat
    if (results.length === 0) {
        searchResults.innerHTML = `
            <div class="search-no-results">
                Aucun personnage trouvé
            </div>
        `;
        searchResults.classList.add('active');
        return;
    }
    
    // Afficher les résultats
    searchResults.innerHTML = results.map(char => `
        <a href="/characters/${char.id}" class="search-result-item">
            <img src="${char.img}" alt="${char.name}" class="search-result-img">
            <div>
                <div class="search-result-name">${char.name}</div>
                <div class="search-result-race">${char.race || 'Inconnu'}</div>
            </div>
        </a>
    `).join('');
    
    searchResults.classList.add('active');
}

// Cacher les résultats
function hideSearchResults() {
    const searchResults = document.getElementById('searchResults');
    if (searchResults) {
        setTimeout(() => {
            searchResults.classList.remove('active');
        }, 200); // Petit délai pour permettre le clic sur un résultat
    }
}

// ========== EVENT LISTENERS ==========
document.addEventListener('DOMContentLoaded', () => {
    const searchInput = document.getElementById('searchInput');
    const searchResults = document.getElementById('searchResults');
    
    if (!searchInput) return; // Pas de barre de recherche sur cette page
    
    // Précharger les personnages au chargement de la page
    loadAllCharacters();
    
    // Event: Saisie dans la barre de recherche
    searchInput.addEventListener('input', async (e) => {
        const query = e.target.value;
        
        if (query.length < 2) {
            hideSearchResults();
            return;
        }
        
        // Rechercher et afficher les résultats
        const results = await searchCharacters(query);
        displaySearchResults(results);
    });
    
    // Event: Perte de focus (cacher les résultats)
    searchInput.addEventListener('blur', () => {
        hideSearchResults();
    });
    
    // Event: Focus (réafficher les résultats si recherche en cours)
    searchInput.addEventListener('focus', async (e) => {
        const query = e.target.value;
        if (query.length >= 2) {
            const results = await searchCharacters(query);
            displaySearchResults(results);
        }
    });
    
    // Event: Touche Entrée (aller au premier résultat)
    searchInput.addEventListener('keydown', async (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            const query = searchInput.value;
            
            if (query.length >= 2) {
                const results = await searchCharacters(query);
                if (results.length > 0) {
                    // Rediriger vers le premier résultat
                    window.location.href = `/characters/${results[0].id}`;
                }
            }
        }
    });
    
    // Event: Clic en dehors de la recherche (fermer les résultats)
    document.addEventListener('click', (e) => {
        if (!searchInput.contains(e.target) && !searchResults.contains(e.target)) {
            hideSearchResults();
        }
    });
});

// ========== STYLES POUR LES RÉSULTATS (ajout dynamique) ==========
// Ajouter un style CSS pour l'élément .search-result-race
const style = document.createElement('style');
style.textContent = `
    .search-result-race {
        font-size: 12px;
        color: #666;
        margin-top: 2px;
    }
`;
document.head.appendChild(style);

console.log('🔍 Module de recherche chargé');