<script lang="ts">
	import { serversStore } from '../../stores/servers.svelte';
	import { userStore } from '../../stores/user.svelte';
	import ChatWindowInput from './chatWindowInput/ChatWindowInput.svelte';
	import ChatWindowMessage from './ChatWindowMessage.svelte';
	import ChatWindowSkeleton from './ChatWindowSkeleton.svelte';

	interface Props {
		id: string;
		channelId: string;
		serverId: string;
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
						userId={message.author.id || ''}
						avatar={message.author.avatar || ''}
						content={message.content}
						displayName={message.author.display_name || 'Name'}
						username={message.author.username || 'username'}
						time={message.created_at}
						isUserMentioned={message.mentions_users?.includes(userStore.user?.id || '')}
						isEdited={message.created_at !== message.updated_at}
						{server}
						{channel}
					/>
				{/each}
			{:else}
				no messages
			{/if}
		{/await}
	</div>
	<ChatWindowInput {channel} {server} />
</ChatWindowSkeleton>
