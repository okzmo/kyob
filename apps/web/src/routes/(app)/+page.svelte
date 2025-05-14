<script lang="ts">
	import ServerMap from '../../components/ServerMap/ServerMap.svelte';
	import { serversStore } from '../../stores/servers.svelte';

	let servers = $derived(serversStore.getServers());
</script>

<svelte:head>
	<title>Kyob | Servers</title>
</svelte:head>

<ServerMap {servers} />
