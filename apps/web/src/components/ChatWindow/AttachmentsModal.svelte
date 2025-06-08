<script lang="ts">
	import { Dialog } from 'bits-ui';
	import Corners from 'components/ui/Corners/Corners.svelte';
	import ArrowLeft from 'components/ui/icons/ArrowLeft.svelte';
	import Close from 'components/ui/icons/Close.svelte';
	import { core } from 'stores/core.svelte';
	import { windows } from 'stores/windows.svelte';

	let attachments = $derived(core.openAttachmentsModal.attachments);
	let maxLimit = $derived(core.openAttachmentsModal.idx >= attachments.length - 1);
	let minLimit = $derived(core.openAttachmentsModal.idx <= 0);

	function previous(e: MouseEvent) {
		e.stopImmediatePropagation();
		core.openAttachmentsModal.idx--;
	}

	function next(e: MouseEvent) {
		e.stopImmediatePropagation();
		core.openAttachmentsModal.idx++;
	}
</script>

<Dialog.Root
	open={core.openAttachmentsModal.status}
	onOpenChange={(s) => {
		if (!s) {
			core.openAttachmentsModal.status = false;
			core.openAttachmentsModal.idx = 0;
			windows.reuseLastWindow();
		}
	}}
>
	<Dialog.Overlay class="bg-main-950/60 fixed inset-0 z-[999] backdrop-blur-[4px]" />
	<Dialog.Content
		class="fixed inset-0 z-[999]"
		onclick={(e) => {
			e.stopImmediatePropagation();
			core.openAttachmentsModal.status = false;
			core.openAttachmentsModal.idx = 0;
			windows.reuseLastWindow();
		}}
	>
		<Dialog.Close
			type="button"
			class="group text-main-50/50 hocus:text-main-50 absolute top-10 right-10 p-2 transition-colors duration-100 hover:cursor-pointer"
		>
			<Corners color="border-main-50/50" class="group-hocus:border-main-50 duration-100" />
			<Close width={20} height={20} />
		</Dialog.Close>

		<div
			class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2"
			onclick={(e) => e.stopImmediatePropagation()}
		>
			<img
				src={attachments[core.openAttachmentsModal.idx]}
				alt=""
				class="max-h-[70vh] max-w-[90vw] object-contain"
				style="display: block;"
			/>
		</div>

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
