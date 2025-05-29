<script lang="ts">
	import ServerMap from '../../components/ServerMap/ServerMap.svelte';
	import { serversStore } from '../../stores/servers.svelte';

	let servers = $derived(serversStore.getServers().filter((s) => s.id !== 'global'));
</script>

<svelte:head>
	<title>Kyob | Servers</title>
</svelte:head>

<ServerMap {servers} />
