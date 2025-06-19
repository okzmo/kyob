import { print } from "utils/basics";

class Audio {
  notification = $state<HTMLAudioElement>();
  muteOn = $state<HTMLAudioElement>();
  muteOff = $state<HTMLAudioElement>();
  callOn = $state<HTMLAudioElement>();
  callOff = $state<HTMLAudioElement>();
  ringTone = $state<HTMLAudioElement>();

  playSound(sound: 'notification' | 'mute-on' | 'mute-off' | 'call-on' | 'call-off' | 'ring-tone') {
    print("playing sound")
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
      case 'ring-tone':
        this.ringTone?.play();
        break;
    }
  }

  stopSound(sound: 'notification' | 'mute-on' | 'mute-off' | 'call-on' | 'call-off' | 'ring-tone') {
    switch (sound) {
      case 'notification':
        if (!this.notification) return
        this.notification.pause();
        this.notification.currentTime = 0;
        break;
      case 'mute-on':
        if (!this.muteOn) return
        this.muteOn.pause();
        this.muteOn.currentTime = 0;
        break;
      case 'mute-off':
        if (!this.muteOff) return
        this.muteOff.pause();
        this.muteOff.currentTime = 0;
        break;
      case 'call-on':
        if (!this.callOn) return
        this.callOn.pause()
        this.callOn.currentTime = 0;
        break;
      case 'call-off':
        if (!this.callOff) return
        this.callOff.pause()
        this.callOff.currentTime = 0;
        break;
      case 'ring-tone':
        if (!this.ringTone) return
        this.ringTone.pause()
        this.ringTone.currentTime = 0;
        break;
    }
  }
}

export const sounds = new Audio();
