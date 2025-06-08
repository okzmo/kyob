<script lang="ts">
	import { core } from 'stores/core.svelte';
	import { windows } from 'stores/windows.svelte';
	import { getFileType } from 'utils/attahcments';

	let { attachments } = $props();
</script>

<div
	class={[
		'pointer-events-auto',
		attachments.length > 1
			? 'grid w-fit max-w-[90%] grid-cols-2 gap-2 @3xl:max-w-[65%] @5xl:max-w-[40%]'
			: 'max-w-[75%] @3xl:max-w-[50%] @5xl:max-w-[35%]'
	]}
>
	{#each attachments as attachment, idx (idx)}
		{@const fileType = getFileType(attachment)}
		{#if fileType === 'image'}
			<button
				class={['attachment relative select-none', attachments.length > 1 && 'aspect-square']}
				onclick={() => {
					core.openAttachmentsModal.status = true;
					core.openAttachmentsModal.attachments = attachments;
					core.openAttachmentsModal.idx = idx;
					windows.setActiveWindow(null);
				}}
			>
				<img src={attachment} alt="Attachment" class="h-full w-full object-cover" />
			</button>
		{:else if fileType === 'video'}
			<video controls>
				<source src={attachment} />
				<track kind="captions" />
			</video>
		{/if}
	{/each}
</div>

<style>
	.attachment::before {
		content: '';
		position: absolute;
		inset: 0;
		box-shadow: inset 0 0 0 1px #fafafa33;
		transition: box-shadow ease-out 75ms;
	}

	.attachment:hover::before {
		box-shadow: inset 0 0 0 1px #fafafa;
		cursor: pointer;
	}
</style>
