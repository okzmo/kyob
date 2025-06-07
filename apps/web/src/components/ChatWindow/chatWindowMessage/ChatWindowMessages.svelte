<script lang="ts">
	import { userStore } from 'stores/user.svelte';
	import ChatWindowMessage from './ChatWindowMessage.svelte';

	let { messages, server, channel } = $props();
	let scrollContent = $state<HTMLElement | null>();

	$effect(() => {
		if (messages.length && scrollContent) {
			scrollContent.scrollTo(0, scrollContent.scrollHeight);
		}
	});
</script>

<div
	bind:this={scrollContent}
	class="@container flex h-[calc(100%-3.5rem)] min-h-0 w-full flex-col gap-y-2 overflow-y-auto pt-2 pb-4"
>
	{#if messages.length > 0}
		{#each messages as message (message.id)}
			<ChatWindowMessage
				id={message.id}
				author={message.author}
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
