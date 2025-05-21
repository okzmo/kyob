<script lang="ts">
	import { onMount } from 'svelte';
	import { backend } from '../../stores/backend.svelte';
	import { goto } from '$app/navigation';
	import { userStore } from '../../stores/user.svelte';
	import { serversStore } from '../../stores/servers.svelte';
	import { page } from '$app/state';

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
			serversStore.servers = res.value.servers;
		}
	});
</script>

{@render children()}
