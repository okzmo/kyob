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
		scrollContainer: HTMLDivElement;
	}

	const DEFAULT_AUTHOR: Partial<User> = {
		id: 'unknown',
		avatar: 'https://i.pinimg.com/736x/c3/7b/e8/c37be8f2419d84e7d38addf481eba9e6.jpg',
		display_name: 'Unknown user',
		username: 'Unknown user'
	};

	let { messages, server, channel, scrollContainer }: Props = $props();

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

	$effect(() => {
		if (messages.length && scrollContainer) {
			scrollContainer.scrollTop = scrollContainer.scrollHeight;
		}
	});
</script>

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
	<div class="flex h-full w-full items-center justify-center">
		<p class="text-main-400 max-w-[20rem] text-center font-semibold">
			It seems nobody sent a message here yet, be the first one!
		</p>
	</div>
{/if}
