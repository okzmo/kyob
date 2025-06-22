<script lang="ts">
	import Sidebar from 'components/ui/Sidebar/Sidebar.svelte';

	let { children } = $props();
</script>

<div class="mx-auto flex max-w-5xl gap-x-10">
	<Sidebar type="general" />
	<main class="h-screen w-full pt-20">
		{@render children()}
	</main>
</div>
