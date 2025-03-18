class ReviewEsWordBlock {
	constructor(scene,word,index) {
    this.reviewed = false
    this.activated = false
		this.scene = scene
		this.word = word
		this.index = index
		this.match_per = 0

		const style_spanish = {
			font: `${this.word.size_spanish}px Arial`, 
			fill: '#ffffff',
			stroke: '#000000',
			strokeThickness:  8
		}

		const style_spanish_masked = {
			font: `${this.word.size_spanish}px Arial`, 
			fill: '#00ff00',
			stroke: '#000000',
			strokeThickness:  8
		}

		const style_sp = {font: `${this.word.size_sp}px Arial`, fill: '#ffffff'}

		this.spanish = scene.add.text(0,0,'',style_spanish)
		this.spanish.visible = false

		this.spanish_masked = scene.add.text(0,0,'',style_spanish_masked)
		this.spanish_masked.visible = false

		this.sp = scene.add.text(0,0,'',style_sp)
		this.sp.visible = false

		// Create a rectangle to use as a mask
		this.maskShape = this.scene.make.graphics();
		this.mask = new Phaser.Display.Masks.GeometryMask(this.scene, this.maskShape);
  }

	set(x_offset,y_offset,txt_spanish,txt_sp){
    this.activated = true
		this.txt_sp = txt_sp

		const block_offset = (this.index * this.word.size_spanish)

		const x_spanish = block_offset + x_offset
		const y_spanish = 0 + y_offset + 16

		const x_sp = block_offset + x_offset + 4
		const y_sp = y_offset + 4

		this.spanish.setText(txt_spanish)
		this.spanish.x = x_spanish
		this.spanish.y = y_spanish
		this.spanish.visible = true

		this.spanish_masked.setText(txt_spanish)
		this.spanish_masked.x = x_spanish
		this.spanish_masked.y = y_spanish
		this.spanish_masked.visible = true

		this.sp.setText(txt_sp.join(''))
		this.sp.x = x_sp
		this.sp.y = y_sp
		this.sp.visible = true
	}

	highlight(user_input) {
    if (this.activated === false) {
      return
    }
		if (user_input === false || user_input === undefined) {
			this.match_per = 0
			return
		}

		const txt_sp_joined = this.txt_sp.join('')
		if (user_input.startsWith(txt_sp_joined)) {
			this.match_per = 1
			return user_input.slice(txt_sp_joined.length);
		} else {
			// partial match
			const matched = this.sequentialMatch(this.txt_sp,user_input)
			if (matched === 0) {
				// no match
				this.match_per = 0
				this.spanish.setColor('#ffffff'); // White
				return false;
			} else {
				// matched
				this.match_per = matched / this.txt_sp.length
				this.update_mask()
				// it still false because its not a complete match
				return false;
			} // if
		}	// if			
	} // highlight

	sequentialMatch(arr, str) {
		let matchCount = 0;
		let currentIndex = 0;
		
		for (let i = 0; i < arr.length; i++) {
			const char = arr[i];
			const found = str.indexOf(char, currentIndex);
			
			if (found !== -1 && found >= currentIndex) {
				matchCount++;
				currentIndex = found + 1;
			}
		}
		
		return matchCount;
	}

	update_mask() {
		this.maskShape.clear();
		
		const width = this.spanish.width;
		const height = this.spanish.height;
		const maskWidth = width * this.match_per;
		
		this.maskShape.fillStyle(0xffffff);
		this.maskShape.fillRect(this.spanish.x, this.spanish.y, maskWidth, height);
		
		this.spanish_masked.setMask(this.mask);
	}

	update(){
		this.update_mask()
	}

	remove(){
		this.activated = false
		this.spanish.visible = false
		this.spanish_masked.visible = false
		this.sp.visible = false
	}
}
