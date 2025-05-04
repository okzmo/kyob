<script lang="ts">
	import '@fontsource/nunito/400.css';
	import '@fontsource/nunito/500.css';
	import '@fontsource/nunito/600.css';
	import '@fontsource/nunito/700.css';
	import '@fontsource/vollkorn/400-italic.css';
	import '@fontsource/outfit/700.css';
	import '@fontsource/outfit/800.css';
	import '@fontsource/outfit/900.css';
	import '../app.css';
	import GridDots from '../components/GridDots.svelte';
	import { onMount } from 'svelte';
	import { backend } from '../stores/backend.svelte';
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { userStore } from '../stores/user.svelte';
	import { serversStore } from '../stores/servers.svelte';

	let { children } = $props();

	onMount(async () => {
		const inAuthPage = ['/signin', '/signup'].includes(page.url.pathname);
		const res = await backend.getSetup();

		if (res.isErr() && !inAuthPage) {
			if (res.error.code === 'ERR_SETUP_UNAUTHORIZED') goto('/signin');
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

<GridDots />
