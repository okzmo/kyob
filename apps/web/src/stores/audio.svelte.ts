class Audio {
  notification = $state<HTMLAudioElement>();
  muteOn = $state<HTMLAudioElement>();
  muteOff = $state<HTMLAudioElement>();
  callOn = $state<HTMLAudioElement>();
  callOff = $state<HTMLAudioElement>();

  playSound(sound: 'notification' | 'mute-on' | 'mute-off' | 'call-on' | 'call-off') {
    switch (sound) {
      case 'notification':
        if (!this.notification) return
        this.notification.currentTime = 0;
        this.notification.play();
        break;
      case 'mute-on':
        if (!this.muteOn) return
        this.muteOn.currentTime = 0;
        this.muteOn.play();
        break;
      case 'mute-off':
        if (!this.muteOff) return
        this.muteOff.currentTime = 0;
        this.muteOff.play();
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
