<script lang="ts">
	import ChatWindowCallBar from './ChatWindowCallBar.svelte';
	import ChatWindowCallParticipants from './ChatWindowCallParticipants.svelte';

	let { server, channel } = $props();
</script>

<div class="flex h-full w-full overflow-auto">
	<ChatWindowCallParticipants {server} {channel} />

	<ChatWindowCallBar {server} {channel} />
</div>
