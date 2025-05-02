<script lang="ts">
	import type { Component } from 'svelte';

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
		class="text-main-400 hocus:text-accent-50 hocus:bg-accent-100/15 flex h-[2.25rem] w-[2.25rem] items-center justify-center rounded-lg transition-colors hover:cursor-pointer"
	>
		<Icon height={22} width={22} />
	</a>
{:else}
	<button
		aria-label={label}
		class="text-main-400 hocus:text-accent-50 hocus:bg-accent-100/15 flex h-[2.25rem] w-[2.25rem] items-center justify-center rounded-lg transition-colors hover:cursor-pointer"
	>
		<Icon height={22} width={22} />
	</button>
{/if}
