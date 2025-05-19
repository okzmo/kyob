class Audio {
	notification = $state<HTMLAudioElement>();

	playSound(sound: 'notification') {
		switch (sound) {
			case 'notification':
				this.notification?.play();
				break;
		}
	}
}

export const sounds = new Audio();
