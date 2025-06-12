<script lang="ts">
	import { page } from '$app/state';
	import { ContextMenu } from 'bits-ui';
	import { core } from 'stores/core.svelte';
	import { contextMenuTargets, type ContextMenuTarget } from 'types/types';
	import Topbar from './Topbar/Topbar.svelte';
	import Userbar from './Userbar/Userbar.svelte';
	import Searchbar from './Searchbar/Searchbar.svelte';
	import ChatWindow from './ChatWindow/ChatWindow.svelte';
	import ContextMenuSkeleton from './ui/ContextMenu/ContextMenuSkeleton.svelte';
	import { windows } from 'stores/windows.svelte';
	import { goback } from 'stores/goback.svelte';
	import Serverbar from './Serverbar/Serverbar.svelte';
	import { searchValidMessageParent } from 'utils/dom';
	import { Tooltip } from 'bits-ui';

	let { children } = $props();

	let contextMenuTarget: string | undefined = $state();
	let contextMenuTargetAuthor: string | undefined = $state();
	let onSettingsPage = $derived(page.url.pathname.includes('/settings'));

	function onContextMenu(e: MouseEvent) {
		const targetId = (e.target as HTMLElement).id;
		const targetAuthor = (e.target as HTMLElement).dataset?.authorId;
		const identifier = targetId.split('-')[0] as ContextMenuTarget;

		if (!contextMenuTargets.includes(identifier)) {
			const { id, author } = searchValidMessageParent(e.target as HTMLElement);
			if (id && author) {
				contextMenuTarget = id;
				contextMenuTargetAuthor = author;
			}
		} else {
			contextMenuTarget = targetId;
			contextMenuTargetAuthor = targetAuthor;

			switch (identifier) {
				case 'inServer':
					core.createChannelModal.x = e.clientX;
					core.createChannelModal.y = e.clientY;
					break;
				case 'mainMap':
					core.createServerModal.x = e.clientX;
					core.createServerModal.y = e.clientY;
					core.joinServerModal.x = e.clientX;
					core.joinServerModal.y = e.clientY;
					break;
			}
		}
	}
</script>

<Tooltip.Provider>
	<ContextMenu.Root>
		<ContextMenu.Trigger
			id={page.params.server_id ? `inServer-${page.params.server_id}` : 'mainMap'}
			class="fixed top-0 left-0 h-screen w-screen"
			oncontextmenu={onContextMenu}
		>
			{#if !onSettingsPage}
				<Topbar canGoBack={goback.active} />
				<Userbar />
				{#if page.params.server_id}
					<Serverbar />
				{/if}
				<Searchbar />
				{#each windows.openWindows as chatWindow (chatWindow.id)}
					<ChatWindow
						id={chatWindow.id}
						tab={chatWindow.tab}
						serverId={chatWindow.serverId}
						channelId={chatWindow.channelId}
						friendId={chatWindow.friendId}
					/>
				{/each}
			{/if}

			{@render children()}
		</ContextMenu.Trigger>
		<ContextMenuSkeleton
			bind:target={contextMenuTarget}
			bind:targetAuthor={contextMenuTargetAuthor}
		/>
	</ContextMenu.Root>
</Tooltip.Provider>
