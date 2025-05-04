<script lang="ts">
	import { onMount } from 'svelte';
	import ChannelMap from '../../../components/ChannelMap/ChannelMap.svelte';
	import type { Channel } from '../../../types/types';
	import { serversStore } from '../../../stores/servers.svelte';
	import { page } from '$app/state';

	let channels = $state<Channel[] | undefined>();

	onMount(() => {
		channels = serversStore.getChannels(Number(page.params.server_id));
	});
</script>

<ChannelMap {channels} />
