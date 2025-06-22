<script lang="ts">
	import { serversStore } from 'stores/servers.svelte';
	import { userStore } from 'stores/user.svelte';
	import ChatWindowInput from './chatWindowInput/ChatWindowInput.svelte';
	import ChatWindowSkeleton from './ChatWindowSkeleton.svelte';
	import ChatWindowMessages from './chatWindowMessage/ChatWindowMessages.svelte';
	import ChatWindowErrors from './ChatWindowErrors.svelte';
	import ChatWindowCall from './chatWindowCall/ChatWindowCall.svelte';

	interface Props {
		id: string;
		tab: 'chat' | 'call';
		channelId?: string;
		serverId?: string;
		friendId?: string;
	}

	let { id, tab, channelId = '', serverId = '', friendId = '' }: Props = $props();

	const server = $derived(serversStore.getServer(serverId));
	const channel = $derived(serversStore.getChannel(serverId, channelId));
	const friend = $derived(userStore.getFriend(friendId));
	const messages = $derived(serversStore.getMessages(serverId, channelId));
	let scrollContainer = $state<HTMLDivElement>();
</script>

<ChatWindowSkeleton {id} {tab} {channel} {server} {friend}>
	{#if tab === 'chat'}
		<div
			class="relative flex min-h-0 w-full flex-grow flex-col gap-y-2 overflow-y-auto pt-2 pb-4"
			bind:this={scrollContainer}
		>
			{#await messages then allMessages}
				<ChatWindowMessages {channel} {server} messages={allMessages} {scrollContainer} />
			{/await}
		</div>
		<ChatWindowInput {channel} {server} {friend} />
	{:else if tab === 'call'}
		<ChatWindowCall {server} {channel} />
	{/if}

	<ChatWindowErrors />
</ChatWindowSkeleton>
