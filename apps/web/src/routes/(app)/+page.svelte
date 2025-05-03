<script lang="ts">
	import ServerMap from '../../components/ServerMap/ServerMap.svelte';
	import { serversStore } from '../../stores/servers.svelte';
</script>

<ServerMap servers={serversStore.servers} />
