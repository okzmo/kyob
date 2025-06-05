<script lang="ts">
	import ChannelMap from 'components/ChannelMap/ChannelMap.svelte';
	import { serversStore } from 'stores/servers.svelte';
	import { page } from '$app/state';
	import { onMount } from 'svelte';

	let server = $derived(serversStore.getServer(page.params.server_id));
	let channels = $derived(serversStore.getChannels(page.params.server_id));

	onMount(() => {
		document.documentElement.style.setProperty('--server-color-85', `rgba(24,24,27, 0.85)`);
		document.documentElement.style.setProperty('--server-color-95', `rgba(24,24,27, 0.95)`);
		document.documentElement.style.setProperty('--server-color', `rgba(24,24,27, 1)`);
	});
</script>

<svelte:head>
	<title>Kyob | {server?.name}</title>
</svelte:head>

<ChannelMap {channels} />
