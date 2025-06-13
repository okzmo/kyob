<script lang="ts">
	import type { Snippet } from 'svelte';
	import Corners from '../Corners/Corners.svelte';
	import Tooltip from '../Tooltip/Tooltip.svelte';

	interface Props {
		variants: 'danger' | 'icon' | 'nostyle';
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
	<Tooltip text={tooltip}>
		<button
			type="button"
			class={[
				'group bg-main-900 text-main-50 relative flex h-full items-center justify-center transition duration-100 hover:cursor-pointer ',
				classes
			]}
			{onclick}
		>
			{#if corners}
				<Corners color={cornerColor} class={cornerClass} />
			{/if}
			{@render children()}
		</button>
	</Tooltip>
{/if}

{#if variants === 'nostyle'}
	<Tooltip text={tooltip}>
		<button type="button" class={classes} {onclick}>
			{#if corners}
				<Corners color={cornerColor} class={cornerClass} />
			{/if}
			{@render children()}
		</button>
	</Tooltip>
{/if}
