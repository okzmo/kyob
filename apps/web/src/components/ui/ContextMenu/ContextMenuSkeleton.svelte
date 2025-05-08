<script lang="ts">
	import { ContextMenu } from 'bits-ui';
	import ServerContextMenu from './ServerContextMenu.svelte';
	import ChannelContextMenu from './ChannelContextMenu.svelte';
	import MessageContextMenu from './MessageContextMenu.svelte';
	import ChannelMapContextMenu from './ChannelMapContextMenu/ChannelMapContextMenu.svelte';

	interface Props {
		target?: string | undefined;
	}

	let { target = $bindable() }: Props = $props();
	let targetId = $derived(Number(target?.split('-')[1]));
</script>

<ContextMenu.Portal>
	{#if target?.includes('serverButton')}
		<ServerContextMenu {targetId} />
	{:else if target?.includes('channelButton')}
		<ChannelContextMenu {targetId} />
	{:else if target?.includes('message')}
		<MessageContextMenu />
	{:else if target?.includes('inServer')}
		<ChannelMapContextMenu />
	{/if}
</ContextMenu.Portal>
