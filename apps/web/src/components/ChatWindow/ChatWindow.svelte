<script lang="ts">
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

	const channel = $state(serversStore.getChannel(serverId, channelId));
	const server = $state(serversStore.getServer(serverId));
</script>

<ChatWindowSkeleton {id} {channel} {server}>
	<div
		class={[
			'flex w-full flex-col-reverse gap-y-2 overflow-auto py-3',
			server?.is_member ? 'h-[calc(100%-3.75rem)]' : 'h-full'
		]}
	>
		{#each { length: 8 }, message}
			<ChatWindowMessage id={message} displayName="Okzmo" username="okzmo" time="Today at 2:30am" />
		{/each}
	</div>
	{#if server?.is_member}
		<ChatWindowInput {channel} {server} />
	{/if}
</ChatWindowSkeleton>
