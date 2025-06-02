<script lang="ts">
	import { onMount } from 'svelte';
	import { backend } from 'stores/backend.svelte';
	import { goto } from '$app/navigation';

	let { children } = $props();

	onMount(async () => {
		const res = await backend.getSetup();

		if (res.isOk()) {
			goto('/');
		}
	});
</script>

{@render children()}
