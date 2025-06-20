<script lang="ts">
	import Sidebar from 'components/ui/Sidebar/Sidebar.svelte';

	let { children } = $props();
</script>

<div class="mx-auto flex max-w-7xl gap-x-10">
	<Sidebar type="server" />
	<main class="h-screen w-full overflow-auto pt-20">
		{@render children()}
	</main>
</div>
