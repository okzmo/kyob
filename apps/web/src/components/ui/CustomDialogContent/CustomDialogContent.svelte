<script lang="ts">
	import type { Snippet } from 'svelte';
	import { Dialog, type WithoutChildrenOrChild } from 'bits-ui';
	import { scaleBlur } from '../../../utils/transition';
	import Corners from '../Corners/Corners.svelte';
	import Close from '../icons/Close.svelte';

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
			<div
				class="bg-main-900 inner-main-800 fixed top-1/2 left-1/2 w-[550px] -translate-1/2"
				{...props}
				transition:scaleBlur={{}}
			>
				<Corners color="border-main-700" />
				<div class="border-b-main-800 relative mb-8 w-full border-b py-7">
					<Dialog.Close
						type="button"
						class="text-main-400 hocus:text-main-50 absolute top-1/2 right-5 -translate-y-1/2 transition-colors hover:cursor-pointer"
					>
						<Close width={18} height={18} />
					</Dialog.Close>
				</div>
				{@render children?.()}
			</div>
		{/if}
	{/snippet}
</Dialog.Content>
