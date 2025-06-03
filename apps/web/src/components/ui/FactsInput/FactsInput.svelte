<script lang="ts">
	import type { Fact } from 'types/types';
	import PlusSimple from '../icons/PlusSimple.svelte';
	import { Dialog } from 'bits-ui';
	import CustomDialogContent from '../CustomDialogContent/CustomDialogContent.svelte';
	import CloseDialogButton from '../CustomDialogContent/CloseDialogButton.svelte';
	import FooterDialog from '../CustomDialogContent/FooterDialog.svelte';
	import SubmitButton from '../SubmitButton/SubmitButton.svelte';
	import { generateRandomId } from 'utils/randomId';
	import Button from '../Button/Button.svelte';

	interface Props {
		facts: Fact[];
	}

	let { facts = $bindable() }: Props = $props();

	let openDialog = $state(false);
	let dialogMode = $state<'create' | 'edit'>('create');
	let editFactId = $state('');
	let label = $state('');
	let value = $state('');

	function openEditMode(factId: string, factLabel: string, factUrl: string) {
		dialogMode = 'edit';
		editFactId = factId;
		label = factLabel;
		value = factUrl;
		openDialog = true;
	}

	function handleRemove() {
		facts = facts.filter((l) => l.id !== editFactId);
		editFactId = '';
		label = '';
		value = '';

		openDialog = false;
	}

	function handleClick() {
		if (dialogMode === 'edit') {
			const factToEdit = facts.find((l) => l.id === editFactId);
			if (!factToEdit) return;
			facts = facts.map((fact) => (fact.id === editFactId ? { ...fact, label, value } : fact));
		} else if (dialogMode === 'create') {
			if (facts?.length >= 3) return;
			facts = [...facts, { id: generateRandomId(), label, value }];
		}

		label = '';
		value = '';
		openDialog = false;
	}
</script>

<div class="flex flex-col">
	<label for="facts" class="text-main-500">Facts</label>
	<ul class="bg-main-900 border-main-800 mt-1.5 flex w-full flex-col gap-y-1 border p-1">
		{#if facts && facts.length > 0}
			{#each facts as fact (fact.id)}
				<button
					type="button"
					class="group bg-main-900 inner-main-800 hocus:bg-accent-100/20 hocus:inner-accent-no-shadow/25 text-main-300 hocus:text-accent-50/80 w-full items-center px-3 py-2 text-left transition duration-100 hover:cursor-pointer"
					onclick={() => openEditMode(fact.id, fact.label, fact.value)}
				>
					<p class="max-w-[20rem] truncate overflow-hidden">
						{fact.label}
						<span class="text-main-50 group-hocus:text-accent-50 transition-colors duration-100">
							{fact.value}
						</span>
					</p>
				</button>
			{/each}
		{/if}
		<li>
			<button
				type="button"
				class={[
					'flex w-full items-center justify-center gap-x-2 py-1.5 transition duration-100',
					facts?.length < 3
						? 'bg-main-900 inner-main-800 hocus:bg-accent-100/20 hocus:inner-accent-no-shadow/25 text-main-300 hocus:text-accent-50 hover:cursor-pointer'
						: 'inner-red-400/20 hocus:bg-red-400/25 hocus:inner-red-400/40 bg-red-400/15 text-red-400 hover:cursor-not-allowed'
				]}
				onclick={() => {
					if (facts.length >= 3) return;
					dialogMode = 'create';
					openDialog = true;
				}}
			>
				{facts?.length || 0} / 3
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
				<Dialog.Title class="text-lg font-semibold">
					{dialogMode === 'create' ? 'Add a fact' : 'Edit your fact'}
				</Dialog.Title>
				<Dialog.Description class="text-main-400 max-w-[24rem] text-sm">
					Your birthday, your first job, tell us about yourself.
				</Dialog.Description>
			</div>
		</div>

		<div class="mt-6 flex flex-col px-8">
			<label for="link-label" class="text-main-500 text-sm">Fact</label>
			<input
				id="fact-label"
				bind:value={label}
				placeholder="My birthday is on the"
				class="bg-main-900 placeholder:text-main-400 hocus:bg-main-800/50 border-main-800 hocus:border-main-700 mt-1.5 transition-colors duration-100 focus:ring-0"
			/>
		</div>

		<div class="mt-4 flex flex-col px-8">
			<label for="link-url" class="text-main-500 text-sm">Answer</label>
			<input
				id="fact-value"
				bind:value
				placeholder="7th of july"
				class="bg-main-900 placeholder:text-main-400 hocus:bg-main-800/50 border-main-800 hocus:border-main-700 mt-1.5 transition-colors duration-100 focus:ring-0"
			/>
		</div>

		<FooterDialog>
			{#if dialogMode === 'edit'}
				<Button variants="danger" onclick={handleRemove}>Remove Fact</Button>
			{/if}
			<SubmitButton type="button" class="relative px-3" onclick={handleClick}>
				{dialogMode === 'create' ? 'Add Fact' : 'Edit Fact'}
			</SubmitButton>
		</FooterDialog>
	</CustomDialogContent>
</Dialog.Root>
