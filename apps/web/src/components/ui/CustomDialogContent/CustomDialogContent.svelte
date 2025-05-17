<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Dialog, type WithoutChildrenOrChild } from 'bits-ui';
	import { scaleBlur } from '../../../utils/transition';

	let {
		ref = $bindable(null),
		children,
		...restProps
	}: WithoutChildrenOrChild<Dialog.ContentProps> & {
		children?: Snippet;
	} = $props();
</script>

<Dialog.Content bind:ref {...restProps} forceMount={true}>
	{#snippet child({ props, open })}
		{#if open}
			<div {...props} transition:scaleBlur={{}}>
				{@render children?.()}
			</div>
		{/if}
	{/snippet}
</Dialog.Content>
