<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		title: string;
		description?: string;
		children: Snippet;
	}

	let { title, description, children }: Props = $props();
</script>

<section class="">
	<h4 class="text-lg font-semibold">{title}</h4>
	{#if description}
		<p class="text-main-400 text-sm">{description}</p>
	{/if}
	{@render children()}
</section>
