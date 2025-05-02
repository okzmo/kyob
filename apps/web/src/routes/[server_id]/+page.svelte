<script lang="ts">
	import ChannelMap from '../../components/ChannelMap/ChannelMap.svelte';
	import { channelsStore } from '../../stores/channels.svelte';
</script>

<ChannelMap channels={channelsStore.channels} />
