<script lang="ts">
	import { userStore } from 'stores/user.svelte';
	import ChatWindowMessage from './ChatWindowMessage.svelte';
	import { serversStore } from 'stores/servers.svelte';
	import type { Channel, Message, Server, User } from 'types/types';
	import { onDestroy, onMount } from 'svelte';

	interface Props {
		messages: Message[];
		server: Server;
		channel: Channel;
	}

	const DEFAULT_AUTHOR: Partial<User> = {
		id: 'unknown',
		avatar: 'https://i.pinimg.com/736x/c3/7b/e8/c37be8f2419d84e7d38addf481eba9e6.jpg',
		display_name: 'Unknown user',
		username: 'Unknown user'
	};

	let { messages, server, channel }: Props = $props();

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
	class="relative flex h-[calc(100%-3.5rem)] min-h-0 w-full flex-col-reverse gap-y-2 overflow-y-auto pt-2 pb-4"
>
	{#if messages.length > 0}
		{#each messages as message (message.id)}
			{@const author = serversStore.getMemberById(server.id, message.author_id)!}
			{@const friend = userStore.getFriend(message.author_id)}

			<ChatWindowMessage
				id={message.id}
				author={server.id === 'global' ? friend || userStore.user! : author || DEFAULT_AUTHOR}
				content={message.content}
				time={message.created_at}
				isUserMentioned={message.mentions_users?.includes(userStore.user?.id || '') ||
					message.everyone}
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
