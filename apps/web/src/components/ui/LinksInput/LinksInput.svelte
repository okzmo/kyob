<script lang="ts">
	import type { Link } from 'types/types';
	import PlusSimple from '../icons/PlusSimple.svelte';
	import { Dialog } from 'bits-ui';
	import CustomDialogContent from '../CustomDialogContent/CustomDialogContent.svelte';
	import CloseDialogButton from '../CustomDialogContent/CloseDialogButton.svelte';
	import FooterDialog from '../CustomDialogContent/FooterDialog.svelte';
	import SubmitButton from '../SubmitButton/SubmitButton.svelte';

	interface Props {
		links: Link[];
	}

	let { links = $bindable() }: Props = $props();

	let openDialog = $state(false);
	let label = $state('');
	let url = $state('');

	function handleClick() {
		links = [...links, { id: '0', label, url }];
		label = '';
		url = '';
	}
</script>

<div class="flex flex-col">
	<label for="links" class="text-main-500">Links</label>
	<ul class="bg-main-900 border-main-800 mt-1.5 flex w-full flex-col gap-y-1 border p-1">
		{#if links && links.length > 0}
			{#each links as link (link.id)}
				<div
					class="bg-main-900 inner-main-800 hocus:bg-accent-100/20 hocus:inner-accent-no-shadow/25 text-main-300 hocus:text-accent-50 flex w-full items-center gap-x-2 px-3 py-1.5 transition duration-100 hover:cursor-pointer"
				>
					{link.label}
				</div>
			{/each}
		{/if}
		<li>
			<button
				type="button"
				class="bg-main-900 inner-main-800 hocus:bg-accent-100/20 hocus:inner-accent-no-shadow/25 text-main-300 hocus:text-accent-50 flex w-full items-center justify-center gap-x-2 py-1.5 transition duration-100 hover:cursor-pointer"
				onclick={() => {
					openDialog = true;
				}}
			>
				0 / 2
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
			<SubmitButton
				type="button"
				class="absolute top-1/2 right-5 -translate-y-1/2 px-3"
				onclick={handleClick}
			>
				Add link
			</SubmitButton>
		</FooterDialog>
	</CustomDialogContent>
</Dialog.Root>
