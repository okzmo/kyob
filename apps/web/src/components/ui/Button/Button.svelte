<script lang="ts">
	import type { Snippet } from 'svelte';
	import Corners from '../Corners/Corners.svelte';

	interface Props {
		variants: 'danger' | 'accent' | 'main';
		onclick: () => void;
		class?: string;
		corners?: boolean;
		children: Snippet;
	}

	let { variants, onclick, class: classes, corners = false, children }: Props = $props();
</script>

{#if variants === 'danger'}
	<button
		type="button"
		class={[
			'group inner-red-400/20 hocus:inner-red-400/40 hocus:bg-red-400/25 relative bg-red-400/15 px-2 py-1 text-red-400 transition duration-100 hover:cursor-pointer',
			classes
		]}
		{onclick}
	>
		{#if corners}
			<Corners color="bg-red-400" />
		{/if}
		{@render children()}
	</button>
{/if}
