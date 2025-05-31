<script lang="ts">
	import { fly } from 'svelte/transition';
	import Corners from '../Corners/Corners.svelte';
	import Check from '../icons/Check.svelte';
	import LoadingIcon from '../icons/LoadingIcon.svelte';
	import type { Snippet } from 'svelte';

	interface Props {
		isEmpty?: boolean;
		isSubmitted?: boolean;
		isSubmitting?: boolean;
		buttonWidth?: number;
		onclick?: () => void;
		class?: string;
		type: 'button' | 'submit';
		children: Snippet;
	}

	let {
		isEmpty,
		isSubmitting,
		isSubmitted,
		buttonWidth,
		type,
		onclick,
		children,
		class: classes
	}: Props = $props();
</script>

<button
	{type}
	class={[
		'group inner-accent/15  bg-accent-100/15 text-accent-50 flex h-8 items-center justify-center overflow-hidden py-1 whitespace-nowrap transition-[box-shadow,color,width] duration-300',
		!isEmpty
			? 'hocus:bg-accent-100/25 hocus:inner-accent-no-shadow/25 hover:cursor-pointer'
			: 'hover:cursor-not-allowed',
		classes
	]}
	style="width: {buttonWidth}px;"
	disabled={isSubmitted || isSubmitting || isEmpty}
	onclick={type === 'button' ? onclick : null}
>
	<Corners color="border-accent-100" />
	{#if isSubmitting}
		<div class="absolute" transition:fly={{ duration: 200, delay: 100, y: 5 }}>
			<LoadingIcon height={20} width={20} />
		</div>
	{:else if isSubmitted}
		<div class="absolute" transition:fly={{ duration: 200, delay: 300, y: 5 }}>
			<Check height={20} width={20} />
		</div>
	{:else}
		<span transition:fly={{ duration: 100, y: 5 }}>{@render children()}</span>
	{/if}
</button>
