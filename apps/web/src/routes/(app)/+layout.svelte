<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { ContextMenu } from 'bits-ui';
	import ChatWindow from 'components/ChatWindow/ChatWindow.svelte';
	import Searchbar from 'components/Searchbar/Searchbar.svelte';
	import Topbar from 'components/Topbar/Topbar.svelte';
	import ContextMenuSkeleton from 'components/ui/ContextMenu/ContextMenuSkeleton.svelte';
	import Userbar from 'components/Userbar/Userbar.svelte';
	import { goback } from 'stores/goback.svelte';
	import { windows } from 'stores/windows.svelte';
	import { onMount } from 'svelte';
	import Audio from 'components/Audio.svelte';
	import GridDots from 'components/GridDots.svelte';
	import Serverbar from 'components/Serverbar/Serverbar.svelte';
	import AddFriendModal from 'components/Topbar/friends/AddFriendModal.svelte';
	import UserProfileNoTrigger from 'components/UserProfile/UserProfileNoTrigger.svelte';
	import { backend } from 'stores/backend.svelte';
	import { core } from 'stores/core.svelte';
	import { serversStore } from 'stores/servers.svelte';
	import { userStore } from 'stores/user.svelte';
	import { contextMenuTargets, type ContextMenuTarget } from 'types/types';

	let contextMenuTarget: string | undefined = $state();
	let contextMenuTargetAuthor: string | undefined = $state();
	let onSettingsPage = $derived(page.url.pathname.includes('/settings'));
	let { children } = $props();

	onMount(async () => {
		const inAuthPage = ['/signin', '/signup'].includes(page.url.pathname);
		const res = await backend.getSetup();

		if (res.isErr() && !inAuthPage) {
			if (res.error.code === 'ERR_UNAUTHORIZED') goto('/signin');
		}

		if (res.isOk() && inAuthPage) {
			goto('/');
		}

		if (res.isOk()) {
			userStore.user = res.value.user;
			userStore.friends = res.value.friends;
			serversStore.setupServers(res.value.servers);
			backend.setupWebsocket(res.value.user.id);

			document.documentElement.style.setProperty(
				'--user-color-85',
				`rgba(${res.value.user.main_color}, 0.85)`
			);
			document.documentElement.style.setProperty(
				'--user-color-95',
				`rgba(${res.value.user.main_color}, 0.95)`
			);
			document.documentElement.style.setProperty(
				'--user-color',
				`rgba(${res.value.user.main_color}, 1)`
			);
		}
	});

	$effect(() => {
		if (page.url.pathname === '/') {
			goback.off();
		} else {
			goback.on();
		}
	});

	function onContextMenu(e: MouseEvent) {
		const targetId = (e.target as HTMLElement).id;
		const targetAuthor = (e.target as HTMLElement).dataset?.authorId;
		const identifier = targetId.split('-')[0] as ContextMenuTarget;
		if (!contextMenuTargets.includes(identifier)) {
			e.preventDefault();
		} else {
			switch (identifier) {
				case 'inServer':
					core.openCreateChannelModal.x = e.clientX;
					core.openCreateChannelModal.y = e.clientY;
					break;
				case 'mainMap':
					core.openCreateServerModal.x = e.clientX;
					core.openCreateServerModal.y = e.clientY;
					core.openJoinServerModal.x = e.clientX;
					core.openJoinServerModal.y = e.clientY;
					break;
			}
			contextMenuTarget = targetId;
			contextMenuTargetAuthor = targetAuthor;
		}
	}
</script>

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

<Audio />
<UserProfileNoTrigger />
<AddFriendModal />

{#if !onSettingsPage}
	<GridDots />
{/if}
