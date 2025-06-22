<script lang="ts">
	import { sounds } from 'stores/audio.svelte';
</script>

<div id="voice-audio-container"></div>

<audio id="notif-sound" bind:this={sounds.notification} src="/sounds/notif-sound.mp3" volume={0.5}
></audio>
<audio id="mute-on-sound" bind:this={sounds.muteOn} src="/sounds/mute-on-sound.mp3" volume={0.5}
></audio>
<audio id="mute-off-sound" bind:this={sounds.muteOff} src="/sounds/mute-off-sound.mp3" volume={0.5}
></audio>
<audio id="call-on-sound" bind:this={sounds.callOn} src="/sounds/call-on-sound.mp3" volume={0.5}
></audio>
<audio id="call-off-sound" bind:this={sounds.callOff} src="/sounds/call-off-sound.mp3" volume={0.5}
></audio>
<audio
	id="ringtone-sound"
	bind:this={sounds.ringTone}
	src="/sounds/ringtone-sound.mp3"
	volume={0.25}
></audio>
