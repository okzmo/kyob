<script lang="ts">
	import { onMount } from 'svelte';
	import Sidebar from '../../../components/ui/Sidebar/Sidebar.svelte';
	import { backend } from '../../../stores/backend.svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { userStore } from '../../../stores/user.svelte';
	import { serversStore } from '../../../stores/servers.svelte';

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
</script>

<div class="mx-auto flex max-w-7xl gap-x-10">
	<Sidebar type="general" />
	<main class="mt-20 h-screen w-full overflow-auto">
		{@render children()}
	</main>
</div>
