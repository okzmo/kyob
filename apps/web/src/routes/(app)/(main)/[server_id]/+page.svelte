<script lang="ts">
	import ChannelMap from 'components/ChannelMap/ChannelMap.svelte';
	import { serversStore } from 'stores/servers.svelte';
	import { page } from '$app/state';

	let server = $derived(serversStore.getServer(page.params.server_id));
	let channels = $derived(serversStore.getChannels(page.params.server_id));

	$effect(() => {
		document.documentElement.style.setProperty(
			'--server-color-85',
			`rgba(${server?.main_color}, 0.85)`
		);
		document.documentElement.style.setProperty(
			'--server-color-95',
			`rgba(${server?.main_color}, 0.95)`
		);
		document.documentElement.style.setProperty('--server-color', `rgba(${server?.main_color}, 1)`);
	});
</script>

<svelte:head>
	<title>Kyob | {server?.name}</title>
</svelte:head>

<ChannelMap {channels} />
