<script lang="ts">
	import type { Component } from 'svelte';
	import Corners from '../Corners/Corners.svelte';

	interface IconProps {
		height: number;
		width: number;
		class?: string;
	}

	interface Props {
		Icon: Component<IconProps>;
		href?: string;
		label: string;
	}

	let { Icon, href, label }: Props = $props();
</script>

{#if href}
	<a
		{href}
		aria-label={label}
		class="text-main-400 group hocus:text-accent-50 hocus:bg-accent-100/15 relative flex h-[2.25rem] w-[2.25rem] items-center justify-center transition hover:cursor-pointer"
	>
		<Corners color="border-accent-100" hide />
		<Icon height={22} width={22} />
	</a>
{:else}
	<button
		aria-label={label}
		class="text-main-400 group hocus:text-accent-50 hocus:bg-accent-100/15 relative flex h-[2.25rem] w-[2.25rem] items-center justify-center transition hover:cursor-pointer"
	>
		<Corners color="border-accent-100" hide />
		<Icon height={22} width={22} />
	</button>
{/if}
