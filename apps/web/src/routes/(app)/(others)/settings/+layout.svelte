<script lang="ts">
	import Sidebar from 'components/ui/Sidebar/Sidebar.svelte';

	let { children } = $props();
</script>

<div class="mx-auto flex max-w-7xl gap-x-10">
	<Sidebar type="general" />
	<main class="mt-20 h-screen w-full overflow-auto">
		{@render children()}
	</main>
</div>
