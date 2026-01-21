
// Ajouter aux favoris
async function addFavorite(type, value) {
    try {
        const response = await fetch(`/favorites/add/${type}/${value}`, {
            method: 'POST'
        });

        const data = await response.json();

        if (response.ok) {
            console.log('✅ Ajouté aux favoris');
            // Mettre à jour l'icône
            updateFavoriteIcon(type, value, true);
            showNotification('Ajouté aux favoris ❤️');
        } else {
            console.error('❌ Erreur:', data.error);
            showNotification('Erreur lors de l\'ajout');
        }
    } catch (error) {
        console.error('❌ Erreur:', error);
        showNotification('Erreur de connexion');
    }
}

// Retirer des favoris
async function removeFavorite(type, value) {
    try {
        const response = await fetch(`/favorites/remove/${type}/${value}`, {
            method: 'POST'
        });

        const data = await response.json();

        if (response.ok) {
            console.log('✅ Retiré des favoris');
            // Recharger la page pour mettre à jour l'affichage
            window.location.reload();
        } else {
            console.error('❌ Erreur:', data.error);
            showNotification('Erreur lors de la suppression');
        }
    } catch (error) {
        console.error('❌ Erreur:', error);
        showNotification('Erreur de connexion');
    }
}

// Mettre à jour l'icône de favori
function updateFavoriteIcon(type, value, isFavorite) {
    const icon = document.querySelector(`[data-favorite-type="${type}"][data-favorite-value="${value}"]`);
    if (icon) {
        if (isFavorite) {
            icon.textContent = '❤️';
            icon.classList.add('is-favorite');
        } else {
            icon.textContent = '🤍';
            icon.classList.remove('is-favorite');
        }
    }
}

// Afficher une notification
function showNotification(message) {
    // Créer l'élément de notification
    const notification = document.createElement('div');
    notification.className = 'favorite-notification';
    notification.textContent = message;
    
    // Ajouter au body
    document.body.appendChild(notification);
    
    // Afficher avec animation
    setTimeout(() => {
        notification.classList.add('show');
    }, 10);
    
    // Retirer après 2 secondes
    setTimeout(() => {
        notification.classList.remove('show');
        setTimeout(() => {
            document.body.removeChild(notification);
        }, 300);
    }, 2000);
}

// Ajouter le style pour les notifications
const style = document.createElement('style');
style.textContent = `
    .favorite-notification {
        position: fixed;
        top: 100px;
        right: 20px;
        background: rgba(58, 78, 68, 0.95);
        color: white;
        padding: 15px 30px;
        border-radius: 10px;
        font-size: 16px;
        z-index: 10000;
        opacity: 0;
        transform: translateX(400px);
        transition: all 0.3s ease;
    }
    
    .favorite-notification.show {
        opacity: 1;
        transform: translateX(0);
    }
    
    .favorite-btn {
        background: none;
        border: none;
        font-size: 24px;
        cursor: pointer;
        transition: transform 0.3s;
        padding: 5px;
    }
    
    .favorite-btn:hover {
        transform: scale(1.2);
    }
    
    .favorite-btn.is-favorite {
        animation: heartBeat 0.3s;
    }
    
    @keyframes heartBeat {
        0%, 100% { transform: scale(1); }
        50% { transform: scale(1.3); }
    }
`;
document.head.appendChild(style);

console.log('❤️ Module favoris chargé');