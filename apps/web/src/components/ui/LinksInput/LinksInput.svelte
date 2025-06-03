<script lang="ts">
	import type { Link } from 'types/types';
	import PlusSimple from '../icons/PlusSimple.svelte';
	import { Dialog } from 'bits-ui';
	import CustomDialogContent from '../CustomDialogContent/CustomDialogContent.svelte';
	import CloseDialogButton from '../CustomDialogContent/CloseDialogButton.svelte';
	import FooterDialog from '../CustomDialogContent/FooterDialog.svelte';
	import SubmitButton from '../SubmitButton/SubmitButton.svelte';
	import { generateRandomId } from 'utils/randomId';
	import Button from '../Button/Button.svelte';

	interface Props {
		links: Link[];
	}

	let { links = $bindable() }: Props = $props();

	let openDialog = $state(false);
	let dialogMode = $state<'create' | 'edit'>('create');
	let editLinkId = $state('');
	let label = $state('');
	let url = $state('');

	function openEditMode(linkId: string, linkLabel: string, linkUrl: string) {
		dialogMode = 'edit';
		editLinkId = linkId;
		label = linkLabel;
		url = linkUrl;
		openDialog = true;
	}

	function handleRemove() {
		links = links.filter((l) => l.id !== editLinkId);
		editLinkId = '';
		label = '';
		url = '';

		openDialog = false;
	}

	function handleClick() {
		if (dialogMode === 'edit') {
			const linkToEdit = links.find((l) => l.id === editLinkId);
			if (!linkToEdit) return;
			links = links.map((link) => (link.id === editLinkId ? { ...link, label, url } : link));
		} else if (dialogMode === 'create') {
			if (links.length >= 2) return;
			links = [...links, { id: generateRandomId(), label, url }];
		}

		editLinkId = '';
		label = '';
		url = '';

		openDialog = false;
	}
</script>

<div class="flex flex-col">
	<label for="links" class="text-main-500">Links</label>
	<ul class="bg-main-900 border-main-800 mt-1.5 flex w-full flex-col gap-y-1 border p-1">
		{#if links && links.length > 0}
			{#each links as link (link.id)}
				<button
					type="button"
					class="bg-main-900 inner-main-800 hocus:bg-accent-100/20 hocus:inner-accent-no-shadow/25 text-main-300 hocus:text-accent-50 flex w-full items-center gap-x-2 px-3 py-1.5 transition duration-100 hover:cursor-pointer"
					onclick={() => openEditMode(link.id, link.label, link.url)}
				>
					{link.label}
				</button>
			{/each}
		{/if}
		<li>
			<button
				type="button"
				class={[
					' flex w-full items-center justify-center gap-x-2 py-1.5 transition duration-100',
					links?.length < 2
						? 'bg-main-900 inner-main-800 hocus:bg-accent-100/20 hocus:inner-accent-no-shadow/25 text-main-300 hocus:text-accent-50 hover:cursor-pointer'
						: 'inner-red-400/20 hocus:bg-red-400/25 hocus:inner-red-400/40 bg-red-400/15 text-red-400 hover:cursor-not-allowed'
				]}
				onclick={() => {
					if (links.length >= 2) return;
					dialogMode = 'create';
					openDialog = true;
				}}
			>
				{links?.length || 0} / 2
				<PlusSimple height={14} width={14} />
			</button>
		</li>
	</ul>
</div>

<Dialog.Root open={openDialog} onOpenChange={(s) => (openDialog = s)}>
	<Dialog.Overlay class="fixed inset-0 z-50 bg-black/20" />
	<CustomDialogContent>
		<CloseDialogButton />
		<div class="flex items-center justify-between px-8">
			<div>
				<Dialog.Title class="text-lg font-semibold">Add a link</Dialog.Title>
				<Dialog.Description class="text-main-400 max-w-[24rem] text-sm">
					Your birthday, your first job, tell us about yourself.
				</Dialog.Description>
			</div>
		</div>

		<div class="mt-6 flex flex-col px-8">
			<label for="link-label" class="text-main-500 text-sm">Link</label>
			<input
				id="link-label"
				bind:value={label}
				placeholder="Portfolio"
				class="bg-main-900 placeholder:text-main-400 hocus:bg-main-800/50 border-main-800 hocus:border-main-700 mt-1.5 transition-colors duration-100 focus:ring-0"
			/>
		</div>

		<div class="mt-4 flex flex-col px-8">
			<label for="link-url" class="text-main-500 text-sm">Url</label>
			<input
				id="link-value"
				bind:value={url}
				placeholder="https://google.com"
				class="bg-main-900 placeholder:text-main-400 hocus:bg-main-800/50 border-main-800 hocus:border-main-700 mt-1.5 transition-colors duration-100 focus:ring-0"
			/>
		</div>

		<FooterDialog>
			{#if dialogMode === 'edit'}
				<Button variants="danger" onclick={handleRemove}>Remove Link</Button>
			{/if}
			<SubmitButton
				type="button"
				class="relative px-3"
				onclick={links.length >= 2 ? undefined : handleClick}
			>
				{dialogMode === 'create' ? 'Add Link' : 'Edit Link'}
			</SubmitButton>
		</FooterDialog>
	</CustomDialogContent>
</Dialog.Root>
