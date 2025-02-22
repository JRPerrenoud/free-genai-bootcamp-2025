class CorrectionSpanishWord {
    constructor(scene) {
        this.scene = scene;
        this.size_en = 16;
        this.size_sp = 48;  // Spanish text size

        const style_english = {
            font: `${this.size_en}px Arial`,
            fill: '#ffffff'
        };

        const x_en = this.scene.game.config.width / 2;
        this.txt_english = scene.add.text(x_en, 0, '', style_english);
        this.txt_english.setOrigin(0.5, 0);
        this.txt_english.visible = false;

        // Single block for Spanish word
        this.block = new CorrectionSpanishWordBlock(scene, this);
    }

    set(word, d_index) {
        this.d_index = d_index;
        this.spanish = word.spanish;
        this.english = word.english;
        this.txt_english.visible = true;

        // Calculate position based on word length and correction box dimensions
        const word_width = this.spanish.length * (this.size_sp / 2);  // Approximate width
        const half_box_len = this.scene.correction_box.width / 2;
        const box_x = this.scene.correction_box.x;
        const box_y = this.scene.correction_box.y;
        const box_h = this.scene.correction_box.height;

        this.x = box_x + half_box_len - (word_width / 2);
        this.y = box_y + 24;

        this.txt_english.setText(this.english);
        this.txt_english.y = box_y + box_h - 64;

        // Set the Spanish word block
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
        this.block.update();
    }

    remove() {
        this.d_index = null;
        this.spanish = null;
        this.english = null;
        this.txt_english.visible = false;
        this.y = -this.size_sp;
        this.block.remove();
    }
}

// Make it available globally
window.CorrectionSpanishWord = CorrectionSpanishWord;
