<script lang="ts">
	import { Tooltip } from 'bits-ui';
	import { fly } from 'svelte/transition';

	interface Props {
		text: string;
		y?: number;
		x?: number;
	}

	let { text, y = 5, x }: Props = $props();
</script>

<Tooltip.Root delayDuration={200}>
	<Tooltip.Trigger>
		{#snippet child({ props })}
			<div {...props} class="absolute inset-0"></div>
		{/snippet}
	</Tooltip.Trigger>
	<Tooltip.Content sideOffset={8} forceMount>
		{#snippet child({ wrapperProps, props, open })}
			{#if open}
				<div {...wrapperProps}>
					<div {...props} class="z-[999]" transition:fly={{ duration: 100, y, x }}>
						<div class="bg-main-900 inner-main-800 !text-main-50 px-2 py-1 text-sm select-none">
							{text}
						</div>
					</div>
				</div>
			{/if}
		{/snippet}
	</Tooltip.Content>
</Tooltip.Root>
