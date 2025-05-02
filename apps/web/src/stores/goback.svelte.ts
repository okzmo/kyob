class GoBack {
	active = $state(false);

	off() {
		this.active = false;
	}

	on() {
		this.active = true;
	}
}

export const goback = new GoBack();
