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

	let scrollContent = $state<HTMLElement | null>();
	let { id, channelId, serverId }: Props = $props();

	const server = $derived(serversStore.getServer(serverId));
	const channel = $derived(serversStore.getChannel(serverId, channelId));
	const messages = $derived(serversStore.getMessages(serverId, channelId));

	$effect(() => {
		messages.then(() => {
			if (scrollContent) scrollContent.scrollTo(0, scrollContent.scrollHeight);
		});
	});
</script>

<ChatWindowSkeleton {id} {channel} {server}>
	<div
		bind:this={scrollContent}
		class="flex h-[calc(100%-3.5rem)] min-h-0 w-full flex-col gap-y-2 overflow-y-auto py-3"
	>
		{#await messages then allMessages}
			{#if allMessages}
				{#each allMessages as message (message.id)}
					<ChatWindowMessage
						id={message.id}
						avatar={message.author.avatar || ''}
						content={message.content}
						displayName={message.author.display_name || 'Name'}
						username={message.author.username || 'username'}
						time={message.created_at}
					/>
				{/each}
			{:else}
				no messages
			{/if}
		{/await}
	</div>
	<ChatWindowInput {channel} {server} />
</ChatWindowSkeleton>
