<script lang="ts">
	import { goto } from '$app/navigation';
	import { page } from '$app/state';
	import { onMount } from 'svelte';

	onMount(() => {
		if (page.url.pathname === '/settings') {
			goto('/settings/account', { replaceState: true });
		}
	});
</script>
