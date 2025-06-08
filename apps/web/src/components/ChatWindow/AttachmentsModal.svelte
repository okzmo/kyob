<script lang="ts">
	import { Dialog } from 'bits-ui';
	import Corners from 'components/ui/Corners/Corners.svelte';
	import ArrowLeft from 'components/ui/icons/ArrowLeft.svelte';
	import { core } from 'stores/core.svelte';

	let attachments = $derived(core.openAttachmentsModal.attachments);
	let index = $state(0);
	let maxLimit = $derived(index >= attachments.length - 1);
	let minLimit = $derived(index <= 0);

	function previous(e: MouseEvent) {
		e.stopImmediatePropagation();
		index--;
	}

	function next(e: MouseEvent) {
		e.stopImmediatePropagation();
		index++;
	}
</script>

<Dialog.Root open={core.openAttachmentsModal.status}>
	<Dialog.Overlay class="fixed inset-0 z-[999] bg-black/20" />
	<Dialog.Content
		class="fixed inset-0 z-[999]"
		onclick={() => {
			core.openAttachmentsModal.status = false;
			index = 0;
		}}
	>
		{#if attachments.length > 1}
			<button
				aria-label="Next"
				class={[
					'group text-main-50/50 absolute top-1/2 left-[10rem] -translate-y-1/2 p-2 transition-colors duration-100',
					minLimit ? 'hover:cursor-not-allowed' : 'hocus:text-main-50 hover:cursor-pointer'
				]}
				onclick={minLimit ? null : previous}
			>
				<Corners
					color="border-main-50/50"
					class={[minLimit ? '' : 'group-hocus:border-main-50 duration-100']}
				/>
				<ArrowLeft height={24} width={24} />
			</button>
		{/if}

		<div
			class="absolute top-1/2 left-1/2
	       -translate-x-1/2 -translate-y-1/2"
			onclick={(e) => e.stopImmediatePropagation()}
		>
			<img
				src={attachments[index]}
				alt=""
				class="max-h-[70vh] max-w-[90vw] object-contain"
				style="display: block;"
			/>
		</div>

		{#if attachments.length > 1}
			<button
				aria-label="Next"
				class={[
					'group text-main-50/50 absolute top-1/2 right-[10rem] -translate-y-1/2 p-2 transition-colors duration-100 ',
					maxLimit ? 'hover:cursor-not-allowed' : 'hocus:text-main-50 hover:cursor-pointer'
				]}
				onclick={maxLimit ? null : next}
			>
				<Corners
					color="border-main-50/50"
					class={[maxLimit ? '' : 'group-hocus:border-main-50 duration-100']}
				/>
				<ArrowLeft height={24} width={24} class="-scale-100" />
			</button>
		{/if}
	</Dialog.Content>
</Dialog.Root>
