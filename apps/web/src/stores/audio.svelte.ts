class Audio {
	notification = $state<HTMLAudioElement>();
	muteOn = $state<HTMLAudioElement>();
	muteOff = $state<HTMLAudioElement>();
	callOn = $state<HTMLAudioElement>();
	callOff = $state<HTMLAudioElement>();

	playSound(sound: 'notification' | 'mute-on' | 'mute-off' | 'call-on' | 'call-off') {
		switch (sound) {
			case 'notification':
				this.notification?.play();
				break;
			case 'mute-on':
				this.muteOn?.play();
				break;
			case 'mute-off':
				this.muteOff?.play();
				break;
			case 'call-on':
				this.callOn?.play();
				break;
			case 'call-off':
				this.callOff?.play();
				break;
		}
	}
}

export const sounds = new Audio();
