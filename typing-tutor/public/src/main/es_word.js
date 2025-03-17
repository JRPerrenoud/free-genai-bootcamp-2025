class EsWord {
	constructor(scene) {
		this.falling = false
		this.targeted = false
		this.scene = scene
		this.size_spanish = 32
		this.size_english = 12

		// Support words up to 15 characters long
		this.blocks = [
			new EsWordBlock(scene, this, 0),
			new EsWordBlock(scene, this, 1),
			new EsWordBlock(scene, this, 2),
			new EsWordBlock(scene, this, 3),
			new EsWordBlock(scene, this, 4),
			new EsWordBlock(scene, this, 5),
			new EsWordBlock(scene, this, 6),
			new EsWordBlock(scene, this, 7),
			new EsWordBlock(scene, this, 8),
			new EsWordBlock(scene, this, 9),
			new EsWordBlock(scene, this, 10),
			new EsWordBlock(scene, this, 11),
			new EsWordBlock(scene, this, 12),
			new EsWordBlock(scene, this, 13),
			new EsWordBlock(scene, this, 14)
		]
	}

	skyline_block_random_unhit_index() {
		// Create an array of indices where hit is false
		const unhitIndices = this.scene.skyline.reduce((indices, item, index) => {
			if (item.hit === false) {
				indices.push(index);
			}
			return indices;
		}, []);
		
		// If there are no unhit items, return -1
		if (unhitIndices.length === 0) {
			return -1;
		}
		
		// Generate a random index from the unhitIndices array
		const randomIndex = Math.floor(Math.random() * unhitIndices.length);
		
		// Return the randomly selected index from the original array
		return unhitIndices[randomIndex];
	}	

	set(word, d_index) {
		this.wordId = word.id
		this.d_index = d_index
		this.spanish = word.spanish
		this.english = word.english
		this.spanish_len = word.spanish.length
		
		const unhit_index = this.skyline_block_random_unhit_index()

		if (unhit_index === -1) {
			if (this.scene.skyline.every(item => item.hit === true)) {
				this.scene.scene.start('SceneReview')
			}
		}

		let offset
		// its 39 instead of 40 because we start at zero
		// we need to subtract 1 from spanish length because 0 will take up a single character
		if (unhit_index + (this.spanish_len - 1) > 39) {
			offset = (unhit_index + (this.spanish_len - 1)) - 39
			this.pos = unhit_index - offset
		} else {
			this.pos = unhit_index
		}

		this.x = this.pos * 32
		
		// Fixed speed for all words
		this.speed = 0.4

		this.y = -this.size_spanish
		this.falling = true

		// Process each character of the Spanish word
		for (let i = 0; i < this.spanish.length; i++) {
			this.blocks[i].set(
				this.x,
				this.y,
				this.spanish.charAt(i),
				i < this.spanish.length ? this.spanish.charAt(i) : '',
				this.size_spanish,
				this.size_english
			)
		}
	}

	highlight(user_input) {
		let remaining_text = user_input
		this.blocks.forEach((block, index) => {
			if (index < this.spanish_len) {
				remaining_text = block.highlight(remaining_text);
			}
		});			
	}

	update() {
		if (this.falling === false) {
			return
		}
		
		const hitline = this.scene.game.config.height - 80
		if (this.y >= hitline) {
			this.hit()
		} else {
			this.y += this.speed
			this.blocks.forEach((block, index) => {
				if (index < this.spanish_len) {
					block.update()
				}
			})
		}
	}

	hit() {
		this.scene.word_manager.last_word = {
			d_index: this.d_index,
			status: 'failed'
		}
		this.scene.userInput = ''
		this.scene.userText.setText(this.userInput)
		this.scene.timer.paused = true
		this.scene.sound.play(this.spanish.replace(/ /g, '_'))
		this.scene.hit(this.spanish_len, this.pos)
		this.scene.mode = 'correction'
		this.scene.correction_box.set(Data.words[this.d_index], this.d_index)
		_player_data.add_failed_word(this.d_index)
		_player_data.send_review(this.wordId, false)
		this.remove()
	}

	x_middle() {
		const len = this.spanish_len * this.size_spanish
		return this.x + (len/2)
	}

	valid_target(user_input) {
		if (this.falling === false) {
			return
		}
		
		// Normalize input and word for accent-insensitive comparison
		const normalizedInput = user_input.normalize("NFD").replace(/[\u0300-\u036f]/g, "").toLowerCase();
		const normalizedSpanish = this.spanish.normalize("NFD").replace(/[\u0300-\u036f]/g, "").toLowerCase();
		
		// First try exact match including accents
		if (this.spanish.toLowerCase() === user_input.toLowerCase()) {
			return true;
		}
		
		// Then try accent-insensitive match
		return normalizedSpanish === normalizedInput;
	}
	
	remove() {
		this.d_index = null
		this.falling = false
		this.spanish = null
		this.english = null
		this.spanish_len = null
		this.y = -this.size_spanish
		this.blocks.forEach((block, index) => {
			block.remove()
		});
	}
}
