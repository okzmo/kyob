<script lang="ts">
	import { defaults, superForm } from 'sveltekit-superforms';
	import { CreateOrUpdateRoleSchema } from 'types/schemas';
	import { valibot } from 'sveltekit-superforms/adapters';
	import FormInput from 'components/ui/FormInput/FormInput.svelte';
	import { backend } from 'stores/backend.svelte';
	import PermissionLine from './PermissionLine.svelte';
	import SubmitButton from 'components/ui/SubmitButton/SubmitButton.svelte';
	import { delay } from 'utils/time';
	import Button from 'components/ui/Button/Button.svelte';
	import type { Role } from 'types/types';

	let {
		roles = $bindable(),
		serverId,
		activeTab,
		PERMISSIONS,
		creatingRole = $bindable(),
		activeRole
	} = $props();

	let isSubmitting = $state(false);
	let isSubmitted = $state(false);
	let buttonWidth = $derived(isSubmitted || isSubmitting ? 40 : 100);
	let rgbColor = $state({ r: 0, g: 0, b: 0, a: 1 });

	const { form, errors, enhance } = superForm(defaults(valibot(CreateOrUpdateRoleSchema)), {
		SPA: true,
		dataType: 'json',
		validators: valibot(CreateOrUpdateRoleSchema),
		validationMethod: 'onsubmit',
		resetForm: false,
		async onUpdate({ form }) {
			if (form.valid) {
				isSubmitting = true;

				if (activeRole) {
				} else if (creatingRole) {
					const rgbaColor = Object.values(rgbColor).join(',');

					form.data.color = rgbaColor.slice(0, rgbaColor.length - 2);
					form.data.index = roles ? roles.length : 0;
					console.log(form.data.index);

					const res = await backend.createRole(serverId, form.data);

					if (res.isErr()) {
						console.error(res.error);
						isSubmitting = false;
					}

					if (res.isOk()) {
						await delay(400);
						isSubmitting = false;
						isSubmitted = true;
						await delay(800);

						roles = roles ? [...roles, res.value] : [res.value];

						creatingRole = false;
					}

					isSubmitted = false;
				}
			}
		}
	});

	$effect(() => {
		if (activeRole) {
			$form.name = activeRole.name;
			$form.abilities = activeRole.abilities;
			$form.color = activeRole.color;

			const [r, g, b] = activeRole.color.split(',');
			rgbColor = { r: Number(r), g: Number(g), b: Number(b), a: 1 };
		} else if (!activeRole && creatingRole) {
			$form.name = '';
			$form.abilities = [];
			$form.color = '';
			rgbColor = { r: 150, g: 150, b: 150, a: 1 };
		}
	});

	async function deleteRole(id: string) {
		const res = await backend.deleteRole(serverId, id);

		if (res.isErr()) {
			console.error(res.error);
		}

		if (res.isOk()) {
			roles = roles.filter((role: Role) => role.id !== id);
			activeRole = undefined;
		}
	}
</script>

<form use:enhance class="h-full w-full">
	{#if activeTab === 'display'}
		<FormInput
			title="Role name"
			id="role-name"
			type="text"
			bind:error={$errors.name}
			bind:inputValue={$form.name}
			placeholder="Admin"
			class="mt-4 max-w-2/3"
			inputClass="w-full"
		/>

		<FormInput
			title="Role color"
			id="role-color"
			type="color-picker"
			bind:error={$errors.color}
			bind:inputValue={rgbColor}
			placeholder="Color"
			class="mt-4 max-w-2/3"
			inputClass="w-full"
		/>
	{:else if activeTab === 'permissions'}
		<ul class="flex h-full flex-col overflow-y-auto">
			{#each PERMISSIONS as permission, idx (idx)}
				<PermissionLine
					active={$form.abilities.includes(permission.ability)}
					label={permission.label}
					description={permission.description}
					ability={permission.ability}
					bind:input={$form.abilities}
				/>
			{/each}
		</ul>
	{/if}

	{#if activeTab === 'display'}
		<div class="mt-3 flex gap-x-2">
			<SubmitButton type="submit" class="relative" {isSubmitted} {isSubmitting} {buttonWidth}>
				{activeRole ? 'Save role' : 'Create role'}
			</SubmitButton>
			{#if activeRole}
				<Button variants="danger" onclick={() => deleteRole(activeRole.id)} corners
					>Delete role</Button
				>
			{/if}
		</div>
	{/if}
</form>
