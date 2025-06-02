<script lang="ts">
	import ChannelMap from 'components/ChannelMap/ChannelMap.svelte';
	import { serversStore } from 'stores/servers.svelte';
	import { page } from '$app/state';

	let server = $derived(serversStore.getServer(page.params.server_id));
	let channels = $derived(serversStore.getChannels(page.params.server_id));
</script>

<svelte:head>
	<title>Kyob | {server?.name}</title>
</svelte:head>

<ChannelMap {channels} />
