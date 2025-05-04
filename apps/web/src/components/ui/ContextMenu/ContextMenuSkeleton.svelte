<script lang="ts">
	import { ContextMenu } from 'bits-ui';
	import ServerContextMenu from './ServerContextMenu.svelte';
	import ChannelContextMenu from './ChannelContextMenu.svelte';
	import MessageContextMenu from './MessageContextMenu.svelte';

	interface Props {
		target?: string | undefined;
	}

	let { target = $bindable() }: Props = $props();
</script>

<ContextMenu.Portal>
	{#if target?.includes('serverButton')}
		<ServerContextMenu targetId={Number(target?.split('-')[1])} />
	{:else if target?.includes('channelButton')}
		<ChannelContextMenu />
	{:else if target?.includes('message')}
		<MessageContextMenu />
	{/if}
</ContextMenu.Portal>
