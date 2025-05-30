<script lang="ts">
	import { ContextMenu } from 'bits-ui';
	import ServerContextMenu from './ServerContextMenu.svelte';
	import ChannelContextMenu from './ChannelContextMenu.svelte';
	import MessageContextMenu from './MessageContextMenu.svelte';
	import ChannelMapContextMenu from './ChannelMapContextMenu/ChannelMapContextMenu.svelte';
	import ServerMapContextMenu from './ServerMapContextMenu/ServerMapContextMenu.svelte';
	import Corners from '../Corners/Corners.svelte';

	interface Props {
		target?: string | undefined;
		targetAuthor?: string | undefined;
	}

	let { target = $bindable(), targetAuthor = $bindable() }: Props = $props();
	let targetId = $derived(target?.split('-')[1] || '');

	function contextMenuMouseDown(e: MouseEvent) {
		e.stopImmediatePropagation();
	}
</script>

<ContextMenu.Portal>
	<ContextMenu.Content
		class="bg-main-900 inner-shadow-main-800 relative flex w-[225px] flex-col gap-y-1 p-2"
		onmousedown={contextMenuMouseDown}
	>
		<Corners color="border-main-700" />
		{#if target?.includes('serverButton')}
			<ServerContextMenu {targetId} />
		{:else if target?.includes('channelButton')}
			<ChannelContextMenu {targetId} />
		{:else if target?.includes('message')}
			<MessageContextMenu authorId={targetAuthor} {targetId} />
		{:else if target?.includes('inServer')}
			<ChannelMapContextMenu />
		{:else if target?.includes('mainMap')}
			<ServerMapContextMenu />
		{/if}
	</ContextMenu.Content>
</ContextMenu.Portal>
