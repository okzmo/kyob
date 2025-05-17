<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Popover, type WithoutChildrenOrChild } from 'bits-ui';
	import { flyBlur } from '../../../utils/transition';

	let {
		x,
		y,
		ref = $bindable(null),
		children,
		...restProps
	}: WithoutChildrenOrChild<Popover.ContentProps> & {
		children?: Snippet;
		x?: number;
		y?: number;
	} = $props();
</script>

<Popover.Content bind:ref {...restProps} forceMount={true}>
	{#snippet child({ wrapperProps, props, open })}
		{#if open}
			<div {...wrapperProps}>
				<div {...props} transition:flyBlur={{ x, y }}>
					{@render children?.()}
				</div>
			</div>
		{/if}
	{/snippet}
</Popover.Content>
