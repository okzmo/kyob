<script lang="ts">
	import { serversStore } from '../../stores/servers.svelte';
	import { userStore } from '../../stores/user.svelte';
	import ChatWindowInput from './chatWindowInput/ChatWindowInput.svelte';
	import ChatWindowMessage from './chatWindowMessage/ChatWindowMessage.svelte';
	import ChatWindowSkeleton from './ChatWindowSkeleton.svelte';

	interface Props {
		id: string;
		channelId?: string;
		serverId?: string;
		friendId?: string;
	}

	let scrollContent = $state<HTMLElement | null>();
	let { id, channelId = '', serverId = '', friendId = '' }: Props = $props();

	const server = $derived(serversStore.getServer(serverId));
	const channel = $derived(serversStore.getChannel(serverId, channelId));
	const friend = $derived(userStore.getFriend(friendId));
	const messages = $derived(serversStore.getMessages(serverId, channelId) ?? []);

	$effect(() => {
		if (messages) {
			messages.then(() => {
				if (scrollContent) scrollContent.scrollTo(0, scrollContent.scrollHeight);
			});
		}
	});
</script>

<ChatWindowSkeleton {id} {channel} {server} {friend}>
	<div
		bind:this={scrollContent}
		class="flex h-[calc(100%-3.5rem)] min-h-0 w-full flex-col gap-y-2 overflow-y-auto py-3"
	>
		{#await messages then allMessages}
			{#if allMessages}
				{#each allMessages as message (message.id)}
					<ChatWindowMessage
						id={message.id}
						author={message.author}
						content={message.content}
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
	<ChatWindowInput {channel} {server} {friend} />
</ChatWindowSkeleton>
