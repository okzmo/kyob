<script lang="ts">
	import type { Snippet } from 'svelte';
	import Corners from '../Corners/Corners.svelte';
	import Tooltip from '../Tooltip/Tooltip.svelte';

	interface Props {
		variants: 'danger' | 'green' | 'icon' | 'nostyle';
		onclick: () => void;
		class?: string | string[];
		cornerColor?: string;
		cornerClass?: string;
		corners?: boolean;
		tooltip?: string;
		children: Snippet;
	}

	let {
		variants,
		onclick,
		class: classes,
		corners = false,
		cornerColor = 'border-main-700',
		cornerClass,
		tooltip,
		children
	}: Props = $props();
</script>

{#if variants === 'green'}
	<button
		type="button"
		class={[
			'group inner-green-400/20 hocus:inner-green-400/40 hocus:bg-green-400/25 relative bg-green-400/15 px-2 py-1 text-green-400 transition duration-100 hover:cursor-pointer',
			classes
		]}
		{onclick}
	>
		{#if corners}
			<Corners color="bg-green-400" />
		{/if}
		{@render children()}
	</button>
{/if}

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

{#if variants === 'icon'}
	<button
		type="button"
		class={[
			'group bg-main-900 text-main-50 relative flex aspect-square items-center justify-center transition duration-100 hover:cursor-pointer',
			classes
		]}
		{onclick}
	>
		{#if tooltip}
			<Tooltip text={tooltip}></Tooltip>
		{/if}
		{#if corners}
			<Corners color={cornerColor} class={cornerClass} />
		{/if}
		{@render children()}
	</button>
{/if}

{#if variants === 'nostyle'}
	<button type="button" class={classes} {onclick}>
		{#if tooltip}
			<Tooltip text={tooltip}></Tooltip>
		{/if}
		{#if corners}
			<Corners color={cornerColor} class={cornerClass} />
		{/if}
		{@render children()}
	</button>
{/if}
