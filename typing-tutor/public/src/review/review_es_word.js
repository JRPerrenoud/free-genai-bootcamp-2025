class ReviewEsWord {
	constructor(scene,index) {
		this.index = index
		this.scene = scene
		this.size_en = 16
		this.size_spanish = 32
		this.size_sp = 12

		const style_english = {
			font: `${this.size_en}px Arial`, 
			fill: '#ffffff'
		}
		this.txt_english = scene.add.text(0,0,'',style_english)
		this.txt_english.visible = true

		this.review_border = new ReviewBorder(scene,this)

		// we can support words 8 characters long
		this.blocks = [
			new ReviewEsWordBlock(scene,this,0),
			new ReviewEsWordBlock(scene,this,1),
			new ReviewEsWordBlock(scene,this,2),
			new ReviewEsWordBlock(scene,this,3),
			new ReviewEsWordBlock(scene,this,4),
			new ReviewEsWordBlock(scene,this,5),
			new ReviewEsWordBlock(scene,this,6),
			new ReviewEsWordBlock(scene,this,7)
		]
	}

	set(word,d_index, pos){
		this.d_index = d_index
		this.spanish_text = word.spanish
		this.spanish = word.spanish.toLowerCase()
		this.english = word.english
		this.spanish_len = word.spanish.length

		const row = Math.floor(this.index / 7)
		this.x = 48 + (row * 256) + (row * 16)
		
		const y_pos = this.index % 7
		const spacing = 40 + 16
		this.y = 48 + (y_pos*(this.size_spanish+spacing))

    this.txt_english.setText(this.english)
		this.txt_english.x = this.x + 8
		this.txt_english.y = this.y + 80 - 24

		word.parts.forEach((part, index) => {
			this.blocks[index].set(
				this.x,
				this.y,
				part.spanish,
				part.spanish.toLowerCase().split(''),
			)
		});
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
		this.review_border.update()
		this.blocks.forEach((block, index) => {
			block.update()
		})
	}
}
