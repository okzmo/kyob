<script lang="ts">
	import type { Snippet } from 'svelte';

	interface Props {
		title: string;
		description?: string;
		children: Snippet;
		class?: string;
	}

	let { title, description, children, class: classes }: Props = $props();
</script>

<section class={classes}>
	<h4 class="text-lg font-semibold">{title}</h4>
	{#if description}
		<p class="text-main-400 text-sm">{description}</p>
	{/if}
	{@render children()}
</section>
