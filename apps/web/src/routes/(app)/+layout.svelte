<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import Audio from 'components/Audio.svelte';
	import Desktop from 'components/Desktop.svelte';
	import GridDots from 'components/GridDots.svelte';
	import AddFriendModal from 'components/modals/AddFriendModal.svelte';
	import AttachmentsModal from 'components/modals/AttachmentsModal.svelte';
	import UserProfileNoTrigger from 'components/UserProfile/UserProfileNoTrigger.svelte';
	import { backend } from 'stores/backend.svelte';
	import { goback } from 'stores/goback.svelte';
	import { serversStore } from 'stores/servers.svelte';
	import { userStore } from 'stores/user.svelte';
	import { onMount } from 'svelte';

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
		}
	});

	$effect(() => {
		if (page.url.pathname === '/') {
			goback.off();
		} else {
			goback.on();
		}
	});
</script>

<Desktop>{@render children()}</Desktop>

<Audio />
<UserProfileNoTrigger />
<AddFriendModal />
<AttachmentsModal />

{#if !onSettingsPage}
	<GridDots />
{/if}
