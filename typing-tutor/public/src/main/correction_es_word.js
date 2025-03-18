class CorrectionEsWord {
	constructor(scene) {
		this.scene = scene
		this.size_en = 16
		this.size_spanish = 48
		this.size_sp = 12
		
		const style_english = {
			font: `${this.size_en}px Arial`, 
			fill: '#ffffff'
		}
    const x_en = this.scene.game.config.width / 2
		this.txt_english = scene.add.text(x_en,0,'',style_english)
		this.txt_english.setOrigin(0.5,0);
		this.txt_english.visible = false

		// we can support words 8 characters long
		this.blocks = [
			new CorrectionEsWordBlock(scene,this,0),
			new CorrectionEsWordBlock(scene,this,1),
			new CorrectionEsWordBlock(scene,this,2),
			new CorrectionEsWordBlock(scene,this,3),
			new CorrectionEsWordBlock(scene,this,4),
			new CorrectionEsWordBlock(scene,this,5),
			new CorrectionEsWordBlock(scene,this,6),
			new CorrectionEsWordBlock(scene,this,7)
		]
	}

	set(word,d_index){
		this.d_index = d_index
		this.spanish_text = word.spanish
		this.spanish = word.spanish.toLowerCase()
		this.english = word.english
		this.spanish_len = word.spanish.length
		this.txt_english.visible = true

    const half_word_len = (this.spanish_len*this.size_spanish) / 2
    const half_box_len =  this.scene.correction_box.width / 2
    const box_x = this.scene.correction_box.x
    const box_y = this.scene.correction_box.y
    const box_h = this.scene.correction_box.height
		this.x = box_x + half_box_len - half_word_len
		this.y = box_y + 24

    this.txt_english.setText(this.english)
		this.txt_english.y = box_y + box_h - 64

		// Process each character of the Spanish word individually
		// instead of trying to access word.parts which doesn't exist
		for (let i = 0; i < this.spanish_text.length; i++) {
			this.blocks[i].set(
				this.x,
				this.y,
				this.spanish_text.charAt(i),
				this.spanish_text.charAt(i).toLowerCase().split(''),
			)
		}
	}

	check_input(user_input){
		// Use normalizeSpanish for accent-insensitive comparison
		const normalizedInput = normalizeSpanish(user_input);
		const normalizedSpanish = normalizeSpanish(this.spanish);
		
		// First try exact match including accents
		if (this.spanish === user_input.toLowerCase()) {
			return true;
		}
		
		// Then try accent-insensitive match
		return normalizedSpanish === normalizedInput;
	}

	highlight(user_input) {
		let remaining_text = user_input
		this.blocks.forEach((block, index) => {
			remaining_text = block.highlight(remaining_text);
		});			
	}

	valid(user_input){
		return this.spanish.toLowerCase() === user_input.toLowerCase()
	}

	update(){
		this.blocks.forEach((block, index) => {
			block.update()
		})
	}

	remove(){
		this.d_index = null
		this.spanish_text = null
		this.spanish = null
		this.spanish_len = null
		this.txt_english.visible = false
		this.y = -this.size_spanish
		this.blocks.forEach((block, index) => {
			block.remove()
		});
	}
}
