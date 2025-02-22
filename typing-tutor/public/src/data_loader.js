// Function to load words from the backend API
async function loadWordsFromBackend(groupId) {
    try {
        const response = await fetch(`http://localhost:5000/api/groups/${groupId}/words`);
        if (!response.ok) {
            throw new Error('Failed to fetch words');
        }
        const data = await response.json();
        
        // Format the data for the game
        window.Data = {
            words: data.items.map(word => ({
                id: word.id,
                english: word.english,
                spanish: word.spanish,
                correct_count: word.correct_count || 0,
                wrong_count: word.wrong_count || 0
            }))
        };
        
        return true;
    } catch (error) {
        console.error('Error loading words:', error);
        return false;
    }
}

// Initialize the game with words from a specific group
async function initGameWithGroup(groupId) {
    const loaded = await loadWordsFromBackend(groupId);
    if (loaded) {
        // Get the game container
        const gameContainer = document.getElementById('game-container');
        
        // Create the game configuration
        const config = {
            type: Phaser.AUTO,
            width: 1280,
            height: 720,
            backgroundColor: '#000000',
            parent: 'game-container',
            scene: [SceneMain, SceneReview]
        };
        
        // Create the game instance
        window.game = new Phaser.Game(config);
    } else {
        // Show error message in the game container
        const gameContainer = document.getElementById('game-container');
        gameContainer.innerHTML = '<p style="color: white;">Failed to load words. Please try again.</p>';
    }
}

// Make functions available globally
window.loadWordsFromBackend = loadWordsFromBackend;
window.initGameWithGroup = initGameWithGroup;
