class SpanishWordBlock {
    constructor(scene, word, index) {
        this.falling = false;
        this.scene = scene;
        this.word = word;
        this.index = index;
        this.match_per = 0;

        // Single style for Spanish text since we don't need kanji/romaji distinction
        const style_text = {
            font: `${this.word.size}px Arial`,
            fill: '#ffffff',
            stroke: '#000000',
            strokeThickness: 4
        };

        const style_text_masked = {
            font: `${this.word.size}px Arial`,
            fill: '#00ff00',
            stroke: '#000000',
            strokeThickness: 4
        };

        this.text = scene.add.text(0, 0, '', style_text);
        this.text.visible = false;

        this.text_masked = scene.add.text(0, 0, '', style_text_masked);
        this.text_masked.visible = false;

        // Create a rectangle to use as a mask
        this.maskShape = this.scene.make.graphics();
        this.mask = new Phaser.Display.Masks.GeometryMask(this.scene, this.maskShape);
    }

    set(x_offset, y_offset, txt_spanish) {
        this.falling = true;
        this.txt_spanish = txt_spanish;

        const block_offset = (this.index * this.word.size);
        const x_pos = block_offset + x_offset;
        const y_pos = y_offset;

        this.text.setText(txt_spanish);
        this.text.x = x_pos;
        this.text.y = y_pos;
        this.text.visible = true;

        this.text_masked.setText(txt_spanish);
        this.text_masked.x = x_pos;
        this.text_masked.y = y_pos;
        this.text_masked.visible = true;
    }

    highlight(user_input) {
        if (this.falling === false) {
            return;
        } else if (this.word.targeted === true) {
            this.match_per = 1;
            return;
        } else if (!user_input) {
            this.match_per = 0;
            return;
        }

        // Normalize both strings for comparison
        const normalizedInput = this.word.normalizeSpanish(user_input.toLowerCase());
        const normalizedSpanish = this.word.normalizeSpanish(this.txt_spanish.toLowerCase());

        if (normalizedSpanish === normalizedInput) {
            this.match_per = 1;
            return '';  // Complete match
        } else if (normalizedSpanish.startsWith(normalizedInput)) {
            // Partial match from the start
            this.match_per = normalizedInput.length / normalizedSpanish.length;
            this.update_mask();
            return false;
        } else {
            // No match
            this.match_per = 0;
            this.text.setColor('#ffffff');
            return false;
        }
    }

    update_mask() {
        if (this.match_per === 0) {
            this.text_masked.visible = false;
            return;
        }

        this.maskShape.clear();
        this.maskShape.fillStyle(0xffffff);
        
        const bounds = this.text.getBounds();
        const maskWidth = bounds.width * this.match_per;
        
        this.maskShape.fillRect(bounds.x, bounds.y, maskWidth, bounds.height);
        this.text_masked.setMask(this.mask);
        this.text_masked.visible = true;
    }

    update() {
        if (this.falling === false) {
            return;
        }

        this.text.y = this.word.y;
        this.text_masked.y = this.word.y;
        this.update_mask();
    }

    destroy() {
        this.text.destroy();
        this.text_masked.destroy();
        this.maskShape.destroy();
    }
}

// Make it available globally
window.SpanishWordBlock = SpanishWordBlock;
