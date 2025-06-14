<script lang="ts">
	import { userStore } from 'stores/user.svelte';
	import ChatWindowMessage from './ChatWindowMessage.svelte';
	import { serversStore } from 'stores/servers.svelte';
	import type { Channel, Message, Server } from 'types/types';
	import { onDestroy, onMount } from 'svelte';

	interface Props {
		messages: Message[];
		server: Server;
		channel: Channel;
	}

	let { messages, server, channel }: Props = $props();
	let scrollContent = $state<HTMLElement | null>();

	onMount(() => {
		if (!channel.last_message_sent || !channel.last_message_read) {
			serversStore.markChannelAsRead(server.id, channel.id);
		}

		if (
			channel.last_message_sent &&
			channel.last_message_read &&
			channel.last_message_sent > channel.last_message_read
		) {
			serversStore.markChannelAsRead(server.id, channel.id);
		}
	});

	onDestroy(() => {
		serversStore.markChannelAsRead(server.id, channel.id);
	});
</script>

<div
	bind:this={scrollContent}
	class="relative flex h-[calc(100%-3.5rem)] min-h-0 w-full flex-col-reverse gap-y-2 overflow-y-auto pt-2 pb-4"
>
	{#if messages.length > 0}
		{#each messages as message (message.id)}
			{@const author = serversStore.getMemberById(server.id, message.author_id)!}
			<ChatWindowMessage
				id={message.id}
				{author}
				content={message.content}
				time={message.created_at}
				isUserMentioned={message.mentions_users?.includes(userStore.user?.id || '')}
				isEdited={message.created_at !== message.updated_at}
				{server}
				{channel}
				attachments={message.attachments || []}
			/>
		{/each}
	{:else}
		no messages
	{/if}
</div>
