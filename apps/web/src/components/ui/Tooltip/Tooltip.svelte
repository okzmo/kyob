<script lang="ts">
	import { Tooltip } from 'bits-ui';
	import type { Snippet } from 'svelte';
	import { fly } from 'svelte/transition';

	interface Props {
		children: Snippet;
		text: string;
		y?: number;
		x?: number;
	}

	let { children, text, y = 5, x }: Props = $props();
</script>

<Tooltip.Provider>
	<Tooltip.Root delayDuration={200}>
		<Tooltip.Trigger>
			{@render children()}
		</Tooltip.Trigger>
		<Tooltip.Content sideOffset={8} forceMount>
			{#snippet child({ wrapperProps, props, open })}
				{#if open}
					<div {...wrapperProps}>
						<div {...props} class="z-[999]" transition:fly={{ duration: 100, y, x }}>
							<div class="bg-main-900 inner-main-800 px-2 py-1 text-sm">
								{text}
							</div>
						</div>
					</div>
				{/if}
			{/snippet}
		</Tooltip.Content>
	</Tooltip.Root>
</Tooltip.Provider>
