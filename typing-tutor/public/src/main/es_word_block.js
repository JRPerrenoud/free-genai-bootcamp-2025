class EsWordBlock {
	constructor(scene, word, index) {
		this.falling = false
		this.scene = scene
		this.word = word
		this.index = index
		this.match_per = 0

		const style_spanish = {
			font: `${this.word.size_spanish}px Arial`, 
			fill: '#ffffff',
			stroke: '#000000',
			strokeThickness: 8
		}

		const style_spanish_masked = {
			font: `${this.word.size_spanish}px Arial`, 
			fill: '#00ff00',
			stroke: '#000000',
			strokeThickness: 8
		}

		const style_english = {
			font: `${this.word.size_english}px Arial`, 
			fill: '#ffffff'
		}

		this.spanish = scene.add.text(0, 0, '', style_spanish)
		this.spanish.visible = false

		this.spanish_masked = scene.add.text(0, 0, '', style_spanish_masked)
		this.spanish_masked.visible = false

		this.english = scene.add.text(0, 0, '', style_english)
		this.english.visible = false

		// Create a rectangle to use as a mask
		this.maskShape = this.scene.make.graphics();
		this.mask = new Phaser.Display.Masks.GeometryMask(this.scene, this.maskShape);
	}

	set(x_offset, y_offset, txt_spanish, txt_english) {
		this.falling = true
		
		// Calculate position based on index
		this.x = x_offset + (this.index * this.word.size_spanish)
		this.y = y_offset
		
		// Set text content
		this.spanish.setText(txt_spanish)
		this.spanish_masked.setText(txt_spanish)
		this.english.setText(txt_english)
		
		// Position text elements
		this.spanish.setPosition(this.x, this.y)
		this.spanish_masked.setPosition(this.x, this.y)
		
		// Position English text below Spanish
		this.english.setPosition(
			this.x + (this.word.size_spanish / 2) - (this.english.width / 2),
			this.y + this.word.size_spanish + 2
		)
		
		// Make text visible
		this.spanish.visible = true
		this.spanish_masked.visible = true
		this.english.visible = true
		
		// Reset mask
		this.maskShape.clear()
		this.maskShape.fillStyle(0xffffff)
		this.maskShape.fillRect(this.x, this.y, 0, this.word.size_spanish)
		this.spanish.setMask(this.mask)
	}

	highlight(user_input) {
		if (this.falling === false) {
			return user_input
		}
		
		const char = this.spanish.text.toLowerCase()
		
		if (user_input.length > 0 && char === user_input[0].toLowerCase()) {
			// Full match for this character
			this.match_per = 1
			this.maskShape.clear()
			this.maskShape.fillStyle(0xffffff)
			this.maskShape.fillRect(
				this.x, 
				this.y, 
				this.word.size_spanish, 
				this.word.size_spanish
			)
			
			// Return remaining input without the matched character
			return user_input.substring(1)
		} else if (user_input.length === 0) {
			// No input, reset mask
			this.match_per = 0
			this.maskShape.clear()
			this.maskShape.fillStyle(0xffffff)
			this.maskShape.fillRect(this.x, this.y, 0, this.word.size_spanish)
		}
		
		// Return unchanged input if no match
		return user_input
	}

	update() {
		if (this.falling === false) {
			return
		}
		
		// Update positions as the word falls
		this.y = this.word.y
		this.x = this.word.x + (this.index * this.word.size_spanish)
		
		this.spanish.setPosition(this.x, this.y)
		this.spanish_masked.setPosition(this.x, this.y)
		this.english.setPosition(
			this.x + (this.word.size_spanish / 2) - (this.english.width / 2),
			this.y + this.word.size_spanish + 2
		)
		
		// Update mask position
		this.maskShape.clear()
		this.maskShape.fillStyle(0xffffff)
		this.maskShape.fillRect(
			this.x, 
			this.y, 
			this.word.size_spanish * this.match_per, 
			this.word.size_spanish
		)
	}

	remove() {
		this.falling = false
		this.spanish.visible = false
		this.spanish_masked.visible = false
		this.english.visible = false
		this.match_per = 0
	}
}
