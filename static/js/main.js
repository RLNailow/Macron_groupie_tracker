// ========== VARIABLES GLOBALES ==========
let isLoginMode = true;

// ========== GESTION DU BOUTON RETOUR ==========
// Afficher le bouton retour uniquement si on n'est pas sur la page d'accueil
window.addEventListener('DOMContentLoaded', () => {
    const backBtn = document.getElementById('backBtn');
    const currentPath = window.location.pathname;
    
    // Afficher le bouton retour sur toutes les pages sauf "/"
    if (currentPath !== '/' && backBtn) {
        backBtn.style.display = 'block';
    }
});

// Fonction pour retourner à l'accueil
function goHome() {
    window.location.href = '/';
}

// ========== GESTION DE LA SIDEBAR ==========
function toggleSidebar() {
    const sidebar = document.getElementById('sidebar');
    const overlay = document.getElementById('sidebarOverlay');

    if (sidebar) {
        sidebar.classList.toggle('active');
    }

    // Empêcher le scroll quand la sidebar est ouverte
    if (sidebar && sidebar.classList.contains('active')) {
        document.body.style.overflow = 'hidden';
    } else {
        document.body.style.overflow = 'auto';
    }
}

// Fermer la sidebar en cliquant sur l'overlay
document.addEventListener('DOMContentLoaded', () => {
    const overlay = document.getElementById('sidebarOverlay');
    if (overlay) {
        overlay.addEventListener('click', toggleSidebar);
    }
});

// ========== GESTION DU MODAL LOGIN/REGISTER ==========
function openAuthModal() {
    const modal = document.getElementById('authModal');
    if (modal) {
        modal.classList.add('active');
        document.body.style.overflow = 'hidden';
    }
}

function closeAuthModal() {
    const modal = document.getElementById('authModal');
    if (modal) {
        modal.classList.remove('active');
        document.body.style.overflow = 'auto';
    }
}

// Fermer le modal en cliquant en dehors
document.addEventListener('DOMContentLoaded', () => {
    const modal = document.getElementById('authModal');
    if (modal) {
        modal.addEventListener('click', (e) => {
            if (e.target.id === 'authModal') {
                closeAuthModal();
            }
        });
    }
});

// Basculer entre Login et Register
function toggleAuthMode() {
    isLoginMode = !isLoginMode;
    
    const title = document.getElementById('authTitle');
    const toggleText = document.getElementById('authToggleText');
    
    if (title && toggleText) {
        if (isLoginMode) {
            title.textContent = 'Connexion';
            toggleText.textContent = "Pas encore de compte ? S'inscrire";
        } else {
            title.textContent = 'Inscription';
            toggleText.textContent = "Déjà un compte ? Se connecter";
        }
    }
}

// ========== GESTION DE L'AUTHENTIFICATION ==========
async function handleAuth(event) {
    event.preventDefault();
    
    const email = document.getElementById('authEmail').value;
    const password = document.getElementById('authPassword').value;
    
    // Validation basique
    if (!email || !password) {
        alert('Veuillez remplir tous les champs');
        return;
    }
    
    // Endpoint selon le mode (login ou register)
    const endpoint = isLoginMode ? '/login' : '/register';
    
    try {
        const response = await fetch(endpoint, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password })
        });
        
        const data = await response.json();
        
        if (response.ok) {
            // Fermer le modal
            closeAuthModal();
            
            // Recharger la page pour mettre à jour l'interface
            window.location.reload();
        } else {
            // Erreur
            alert(data.error || 'Erreur lors de l\'authentification');
        }
    } catch (error) {
        alert('Erreur de connexion au serveur');
    }
}

// ========== DÉCONNEXION ==========
function logout() {
    // Rediriger vers /logout qui supprime le cookie
    window.location.href = '/logout';
}

// ========== UTILITAIRES ==========
// Fonction pour vider la barre de recherche (utilisée sur la page d'accueil)
function clearSearch() {
    const searchInput = document.getElementById('searchInput');
    const searchResults = document.getElementById('searchResults');
    
    if (searchInput) {
        searchInput.value = '';
        searchInput.focus();
    }
    
    if (searchResults) {
        searchResults.innerHTML = '';
        searchResults.classList.remove('active');
    }
}

// ========== ANIMATIONS AU SCROLL ==========
// Ajouter une classe quand on scroll (pour des effets futurs)
window.addEventListener('scroll', () => {
    const header = document.querySelector('.header');
    if (header) {
        if (window.scrollY > 50) {
            header.style.boxShadow = '0 4px 20px rgba(0, 0, 0, 0.2)';
        } else {
            header.style.boxShadow = 'none';
        }
    }
});

// ========== TOUCHES CLAVIER ==========
// Fermer les modals avec la touche Échap
document.addEventListener('keydown', (e) => {
    if (e.key === 'Escape') {
        // Fermer la sidebar
        const sidebar = document.getElementById('sidebar');
        if (sidebar && sidebar.classList.contains('active')) {
            toggleSidebar();
        }
        
        // Fermer le modal auth
        const modal = document.getElementById('authModal');
        if (modal && modal.classList.contains('active')) {
            closeAuthModal();
        }
    }
});

// ========== FONCTION "CHOOSE FOR ME !" ==========
async function goToRandomCharacter() {
    try {
        const response = await fetch('/api/random-character');

        if (!response.ok) {
            throw new Error('Erreur lors de la récupération d\'un personnage aléatoire');
        }

        const data = await response.json();
        window.location.href = `/characters/${data.id}`;
    } catch (error) {
        alert('Erreur lors de la sélection d\'un personnage aléatoire: ' + error.message);
    }
}