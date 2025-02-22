class ReviewSpanishWord {
    constructor(scene, index) {
        this.index = index;
        this.scene = scene;
        this.size_en = 16;
        this.size_sp = 32;  // Spanish text size

        const style_english = {
            font: `${this.size_en}px Arial`,
            fill: '#ffffff'
        };
        this.txt_english = scene.add.text(0, 0, '', style_english);
        this.txt_english.visible = true;

        this.review_border = new ReviewBorder(scene, this);
        
        // Single block for Spanish word
        this.block = new ReviewSpanishWordBlock(scene, this);
    }

    set(word, d_index, pos) {
        this.d_index = d_index;
        this.spanish = word.spanish;
        this.english = word.english;

        // Calculate position based on grid layout
        const row = Math.floor(this.index / 7);
        this.x = 48 + (row * 256) + (row * 16);
        
        const y_pos = this.index % 7;
        const spacing = 40 + 16;
        this.y = 48 + (y_pos * (this.size_sp + spacing));

        // Set English text
        this.txt_english.setText(this.english);
        this.txt_english.x = this.x + 8;
        this.txt_english.y = this.y + 80 - 24;

        // Set Spanish word block
        this.block.set(this.x, this.y, this.spanish);
    }

    highlight(user_input) {
        this.block.highlight(user_input);
    }

    valid(user_input) {
        // Normalize both strings for comparison
        const normalizedInput = this.normalizeSpanish(user_input.toLowerCase());
        const normalizedSpanish = this.normalizeSpanish(this.spanish.toLowerCase());
        return normalizedSpanish === normalizedInput;
    }

    normalizeSpanish(text) {
        return text.normalize("NFD").replace(/[\u0300-\u036f]/g, "");
    }

    update() {
        this.review_border.update();
        this.block.update();
    }
}

// Make it available globally
window.ReviewSpanishWord = ReviewSpanishWord;
