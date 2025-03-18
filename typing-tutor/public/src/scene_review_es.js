class SceneReviewEs extends Phaser.Scene {
	constructor() {
		super('SceneReviewEs');
	}	

	create() {
		this.userInput = '';

		this.marker_pos = 0

		this.score_text = this.add.text(32, 24, '', {
			font: '16px Arial',
			fill: '#ffffff'
		}).setOrigin(0.5);
		this.score_text.setText("Game Over");

		this.userText = this.add.text(this.game.config.width / 2, this.game.config.height / 2, '', {
			font: '40px Arial',
			fill: '#434343'
		}).setOrigin(0.5);

		this.failed_words = []
		_player_data.failed_words.forEach((d_index,index) => {
			const data = Data.words[d_index]
			const word = new ReviewEsWord(this,index)
			word.set(data,d_index,index)
			this.failed_words.push(word)
		})

		// Enable keyboard input
		this.input.keyboard.on('keydown', (event) => {
			if  (event.key === 'Dead') {
				// Ignore the 'Dead' key event
				return;
			} else if (event.key === 'q') {
				this.scene.start('SceneMainEs')
			} else if (event.key === 'Shift') {
				this.userInput = ''
			} else if (event.keyCode === 8 && this.userInput.length > 0) {
				// Handle backspace
				this.userInput = this.userInput.slice(0, -1);
			} else if (event.keyCode === 13) {
				this.reviewed()
				return;
			} else if ((event.keyCode >= 65 && event.keyCode <= 90) || event.keyCode === 32) {
				// Add character for letters and space
				this.userInput += event.key.toLowerCase();
			}
			this.userText.setText(this.userInput);
		});
	}

	reviewed(){
		// Add boundary check to prevent accessing beyond array length
		if (this.marker_pos >= this.failed_words.length) {
			// All words reviewed, could show a message or return to main game
			this.userText.setText("All words reviewed! Press 'q' to continue");
			return;
		}
		
		if (this.failed_words[this.marker_pos].valid(this.userInput)) {
			this.failed_words[this.marker_pos].reviewed = true
			this.marker_pos += 1
			this.userInput = ''
			this.userText.setText(this.userInput);
			
			// Add another check after incrementing marker_pos
			if (this.marker_pos >= this.failed_words.length) {
				this.userText.setText("All words reviewed! Press 'q' to continue");
			}
		}
	}

	update() {
		this.failed_words.forEach((word,index) => {
			word.highlight(this.userInput)
			word.update()
		})
	}
}
