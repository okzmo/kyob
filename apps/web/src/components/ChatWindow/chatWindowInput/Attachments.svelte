<script lang="ts">
	import Close from 'components/ui/icons/Close.svelte';
	import FileIcon from 'components/ui/icons/File.svelte';
	import { onDestroy } from 'svelte';

	interface Props {
		attachments: File[];
	}

	let { attachments = $bindable() }: Props = $props();
	let objectUrls = $state<Record<number, string | null>>([]);

	function removeAttachment(file: File, idx: number) {
		attachments = attachments.filter((f) => f.name !== file.name);
		delete objectUrls[idx];
	}

	function getFileType(file: File): 'image' | 'video' | 'unknown' {
		if (file.type.startsWith('image/')) return 'image';
		if (file.type.startsWith('video/')) return 'video';
		return 'unknown';
	}

	$effect(() => {
		attachments.forEach((file, idx) => {
			const fileType = getFileType(file);
			if (fileType === 'image' || fileType === 'video') {
				objectUrls[idx] = URL.createObjectURL(file);
			} else {
				objectUrls[idx] = null;
			}
		});
	});

	onDestroy(() => {
		Object.values(objectUrls).forEach((url) => url && URL.revokeObjectURL(url));
	});
</script>

<div class="flex gap-x-2 px-2 pb-2">
	{#each attachments as attachment, idx (idx)}
		{@const fileType = getFileType(attachment)}

		{#if fileType === 'image' && objectUrls[idx]}
			<figure class="attachment-input relative aspect-square h-20 w-20">
				<button
					class="hocus:border-red-400 hocus:bg-red-400/40 absolute top-1 right-1 border border-red-400/50 bg-red-400/20 transition-colors duration-100 hover:cursor-pointer"
					onclick={() => removeAttachment(attachment, idx)}
				>
					<Close height={16} width={16} class="text-red-400" />
				</button>
				<img src={objectUrls[idx]} alt={attachment.name} class="h-full w-full object-cover" />
			</figure>
		{:else if fileType === 'video'}
			<figure class="attachment-input relative aspect-square h-20 w-20">
				<button
					class="hocus:border-red-400 hocus:bg-red-400/40 absolute top-1 right-1 border border-red-400/50 bg-red-400/20 transition-colors duration-100 hover:cursor-pointer"
					onclick={() => removeAttachment(attachment, idx)}
				>
					<Close height={16} width={16} class="text-red-400" />
				</button>
				<video class="h-full w-full object-cover">
					<source src={objectUrls[idx]} />
					<track kind="captions" />
				</video>
			</figure>
		{:else}
			<figure
				class="attachment-input relative flex aspect-square h-20 w-20 flex-col items-center justify-center gap-y-1"
			>
				<button
					class="hocus:border-red-400 hocus:bg-red-400/40 absolute top-1 right-1 border border-red-400/50 bg-red-400/20 transition-colors duration-100 hover:cursor-pointer"
					onclick={() => removeAttachment(attachment, idx)}
				>
					<Close height={16} width={16} class="text-red-400" />
				</button>
				<FileIcon height={32} width={32} />
				<p class="w-[calc(100%-1.5rem)] truncate">{attachment.name}</p>
			</figure>
		{/if}
	{/each}
</div>

<style>
	.attachment-input::before {
		content: '';
		position: absolute;
		inset: 0;
		box-shadow: inset 0 0 0 1px #fafafa33;
	}
</style>
