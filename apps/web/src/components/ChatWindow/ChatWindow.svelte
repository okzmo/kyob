<script lang="ts">
	import { channelsStore } from '../../stores/channels.svelte';
	import { serversStore } from '../../stores/servers.svelte';
	import ChatWindowInput from './ChatWindowInput.svelte';
	import ChatWindowMessage from './ChatWindowMessage.svelte';
	import ChatWindowSkeleton from './ChatWindowSkeleton.svelte';

	interface Props {
		id: string;
		channelId: number;
		serverId: number;
	}

	let { id, channelId, serverId }: Props = $props();

	const channel = $state(channelsStore.getChannel(channelId));
	const server = $state(serversStore.getServer(serverId));
</script>

<ChatWindowSkeleton {id} {channel} {server}>
	<div class="flex h-[calc(100%-3.75rem)] w-full flex-col-reverse gap-y-2 overflow-auto py-3">
		{#each { length: 8 }}
			<ChatWindowMessage displayName="Okzmo" username="okzmo" time="Today at 2:30am" />
		{/each}
	</div>
	<ChatWindowInput {channel} {server} />
</ChatWindowSkeleton>
