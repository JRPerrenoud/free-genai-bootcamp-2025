class SpanishWord {
    constructor(scene) {
        this.falling = false;
        this.targeted = false;
        this.scene = scene;
        this.size = 32; // Single size since we don't need separate sizes for different scripts
        
        // Initialize with empty values
        this.wordId = null;
        this.d_index = null;
        this.spanish = '';
        this.english = '';
        this.x = 0;
        this.y = 0;
        this.speed = 0.4;
        
        // Create text objects for Spanish and English
        this.spanishText = scene.add.text(0, 0, '', {
            fontSize: this.size + 'px',
            fontFamily: 'Arial',
            fill: '#ffffff'
        });
        this.englishText = scene.add.text(0, 0, '', {
            fontSize: (this.size / 2) + 'px',
            fontFamily: 'Arial',
            fill: '#888888'
        });
        
        // Hide texts initially
        this.spanishText.setVisible(false);
        this.englishText.setVisible(false);
    }

    skyline_block_random_unhit_index() {
        // Create an array of indices where hit is false
        const unhitIndices = this.scene.skyline.reduce((indices, item, index) => {
            if (item.hit === false) {
                indices.push(index);
            }
            return indices;
        }, []);
        
        if (unhitIndices.length === 0) {
            return -1;
        }
        
        const randomIndex = Math.floor(Math.random() * unhitIndices.length);
        return unhitIndices[randomIndex];
    }   

    set(word, d_index) {
        this.wordId = word.id;
        this.d_index = d_index;
        this.spanish = word.spanish;
        this.english = word.english;
        
        const unhit_index = this.skyline_block_random_unhit_index();

        if (unhit_index === -1) {
            if (this.scene.skyline.every(item => item.hit === true)) {
                this.scene.scene.start('SceneReview');
            }
            return;
        }

        // Calculate position based on word length
        const wordWidth = this.spanish.length * (this.size / 2); // Approximate width
        const maxPos = 40 - Math.ceil(wordWidth / this.size); // 40 is the game width in blocks
        
        // Ensure the word fits within the game bounds
        this.pos = Math.min(unhit_index, maxPos);
        this.x = this.pos * this.size;
        this.y = -this.size; // Start above the screen
        
        // Update text objects
        this.spanishText.setText(this.spanish);
        this.englishText.setText(this.english);
        this.spanishText.setPosition(this.x, this.y);
        this.englishText.setPosition(this.x, this.y + this.size);
        
        // Show texts and start falling
        this.spanishText.setVisible(true);
        this.englishText.setVisible(true);
        this.falling = true;
    }

    highlight(userInput) {
        // If no input, reset highlighting
        if (!userInput) {
            this.spanishText.setColor('#ffffff');
            return '';
        }

        const normalizedInput = this.normalizeSpanish(userInput.toLowerCase());
        const normalizedWord = this.normalizeSpanish(this.spanish.toLowerCase());

        // Check if input matches the start of the word
        if (normalizedWord.startsWith(normalizedInput)) {
            // Highlight the matched portion in green
            this.spanishText.setColor('#00ff00');
            return userInput;
        } else {
            // Show error in red
            this.spanishText.setColor('#ff0000');
            return userInput;
        }
    }

    // Normalize Spanish text for comparison (remove accents for easier typing)
    normalizeSpanish(text) {
        return text.normalize("NFD").replace(/[\u0300-\u036f]/g, "");
    }

    update() {
        if (!this.falling) {
            return;
        }

        this.y += this.speed;
        this.spanishText.setPosition(this.x, this.y);
        this.englishText.setPosition(this.x, this.y + this.size);

        // Check if word has hit the bottom
        if (this.y >= this.scene.game.config.height - this.size) {
            this.falling = false;
            this.destroy();
            this.scene.word_missed(this);
        }
    }

    destroy() {
        this.spanishText.destroy();
        this.englishText.destroy();
        this.falling = false;
    }

    match(input) {
        const normalizedInput = this.normalizeSpanish(input.toLowerCase());
        const normalizedWord = this.normalizeSpanish(this.spanish.toLowerCase());
        return normalizedInput === normalizedWord;
    }
}

// Make it available globally
window.SpanishWord = SpanishWord;
