<script lang="ts">
	import { sounds } from 'stores/audio.svelte';
</script>

<audio id="notif-sound" bind:this={sounds.notification} src="/sounds/notif-sound.mp3"></audio>
