<script lang="ts">
	import { page } from '$app/state';
	import { goback } from '../../stores/goback.svelte';
	import Topbar from '../../components/Topbar/Topbar.svelte';
	import Userbar from '../../components/Userbar/Userbar.svelte';
	import Searchbar from '../../components/Searchbar/Searchbar.svelte';
	import { windows } from '../../stores/windows.svelte';
	import ChatWindow from '../../components/ChatWindow/ChatWindow.svelte';
	import { ContextMenu } from 'bits-ui';
	import ContextMenuSkeleton from '../../components/ui/ContextMenu/ContextMenuSkeleton.svelte';
	import { contextMenuTargets, type ContextMenuTarget } from '../../types/types';
	import CreateServerButton from '../../components/ui/CreateServerButton/CreateServerButton.svelte';
	import Serverbar from '../../components/Serverbar/Serverbar.svelte';

	let contextMenuTarget: string | undefined = $state();
	let { children } = $props();

	$effect(() => {
		if (page.url.pathname === '/') {
			goback.off();
		} else {
			goback.on();
		}
	});

	function onContextMenu(e: MouseEvent) {
		const targetId = (e.target as HTMLElement).id;
		const identifier = targetId.split('-')[0] as ContextMenuTarget;
		if (!contextMenuTargets.includes(identifier)) {
			e.preventDefault();
		} else {
			contextMenuTarget = targetId;
		}
	}
</script>

<ContextMenu.Root>
	<ContextMenu.Trigger class="fixed top-0 left-0 h-screen w-screen" oncontextmenu={onContextMenu}>
		<Topbar canGoBack={goback.active} />
		<Userbar />
		{#if !page.params.server_id}
			<CreateServerButton />
		{:else}
			<Serverbar />
		{/if}
		<Searchbar />

		{@render children()}

		{#each windows.openWindows as chatWindow (chatWindow.id)}
			<ChatWindow
				id={chatWindow.id}
				serverId={chatWindow.serverId}
				channelId={chatWindow.channelId}
			/>
		{/each}
	</ContextMenu.Trigger>
	<ContextMenuSkeleton bind:target={contextMenuTarget} />
</ContextMenu.Root>
