<script lang="ts">
	import { valibot } from 'sveltekit-superforms/adapters';
	import Corners from '../ui/Corners/Corners.svelte';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { UpdatePasswordSchema } from '../../types/schemas';
	import { backend } from '../../stores/backend.svelte';

	const { form, errors, enhance } = superForm(defaults(valibot(UpdatePasswordSchema)), {
		SPA: true,
		validators: valibot(UpdatePasswordSchema),
		validationMethod: 'onsubmit',
		async onUpdate({ form }) {
			if (form.valid) {
			}
		}
	});
</script>

<form class="mt-2 flex flex-col gap-y-3" use:enhance>
	<div class="flex flex-col gap-y-1">
		<div class="flex items-center">
			<label for="currentPassword" class={['w-fit', $errors.current_password && 'text-red-400']}
				>Current password</label
			>
			{#if $errors.current_password}
				<p class="text-red-400">- {$errors.current_password}</p>
			{/if}
		</div>
		<input
			id="currentPassword"
			name="currentPassword"
			autocomplete="current-password"
			type="password"
			class={[
				'bg-main-900 placeholder:text-main-400 hocus:bg-main-800/50 w-full transition-colors duration-100 focus:ring-0',
				$errors.current_password ? 'border-red-400' : 'border-main-800 hocus:border-main-700'
			]}
			placeholder="Current password"
			bind:value={$form.current_password}
		/>
	</div>

	<div class="flex flex-col gap-y-1">
		<div class="flex items-center">
			<label for="newPassword" class={['w-fit', $errors.new_password && 'text-red-400']}>
				New password
			</label>
			{#if $errors.new_password}
				<p class="text-red-400">- {$errors.new_password}</p>
			{/if}
		</div>
		<input
			id="newPassword"
			name="newPassword"
			autocomplete="off"
			type="new-password"
			class={[
				'bg-main-900 placeholder:text-main-400 hocus:bg-main-800/50 w-full transition-colors duration-100 focus:ring-0',
				$errors.new_password ? 'border-red-400' : 'border-main-800 hocus:border-main-700'
			]}
			placeholder="New password"
			bind:value={$form.new_password}
		/>
	</div>

	<div class="flex flex-col gap-y-1">
		<div class="flex items-center">
			<label for="confirm" class={['w-fit', $errors.confirm && 'text-red-400']}>
				Confirm password
			</label>
			{#if $errors.confirm}
				<p class="text-red-400">- {$errors.confirm}</p>
			{/if}
		</div>
		<input
			id="confirm"
			name="confirm"
			type="confirm"
			autocomplete="new-password"
			class={[
				'bg-main-900 placeholder:text-main-400 hocus:bg-main-800/50 max-w-[20rem] transition-colors duration-100 focus:ring-0',
				$errors.confirm ? 'border-red-400' : 'border-main-800 hocus:border-main-700'
			]}
			placeholder="Confirm new password"
			bind:value={$form.confirm}
		/>
	</div>

	<button
		class="group inner-accent/15 hocus:bg-accent-100/25 hocus:inner-accent-no-shadow/25 bg-accent-100/15 text-accent-50 relative mt-2 block w-fit px-4 py-1 text-left transition duration-100 hover:cursor-pointer"
	>
		<Corners color="border-accent-100" />
		Save my informations
	</button>
</form>
