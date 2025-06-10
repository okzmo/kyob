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
	const messages = $derived(serversStore.getMessages(serverId, channelId) ?? []);
</script>

<ChatWindowSkeleton {id} {tab} {channel} {server} {friend}>
	{#if tab === 'chat'}
		{#await messages then allMessages}
			<ChatWindowMessages {channel} {server} messages={allMessages} />
		{/await}
		<ChatWindowInput {channel} {server} {friend} />
	{:else if tab === 'call'}
		<ChatWindowCall />
	{/if}

	<ChatWindowErrors />
</ChatWindowSkeleton>
