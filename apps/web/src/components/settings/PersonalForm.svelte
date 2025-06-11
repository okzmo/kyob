<script lang="ts">
	import { valibot } from 'sveltekit-superforms/adapters';
	import { userStore } from 'stores/user.svelte';
	import { defaults, setError, superForm } from 'sveltekit-superforms';
	import { UpdateAccountSchema } from 'types/schemas';
	import { backend } from 'stores/backend.svelte';
	import { delay } from '../../utils/delay';
	import SubmitButton from '../ui/SubmitButton/SubmitButton.svelte';

	let isSubmitting = $state(false);
	let isSubmitted = $state(false);
	let buttonWidth = $derived(isSubmitted || isSubmitting ? 40 : 184);

	const { form, errors, enhance } = superForm(defaults(valibot(UpdateAccountSchema)), {
		SPA: true,
		validators: valibot(UpdateAccountSchema),
		validationMethod: 'onsubmit',
		async onUpdate({ form }) {
			if (form.valid) {
				isSubmitting = true;
				const res = await backend.updateAccount({
					username: form.data.username?.trim(),
					email: form.data.email?.trim()
				});

				if (res.isErr()) {
					if (res.error.code === 'ERR_EMAIL_IN_USE') {
						setError(form, 'email', 'Already in use.');
					}

					if (res.error.code === 'ERR_USERNAME_IN_USE') {
						setError(form, 'username', 'Already in use.');
					}

					if (res.error.code === 'ERR_UNKNOWN') {
						console.log(res.error.error);
					}
					isSubmitting = false;
					return;
				}

				if (res.isOk()) {
					if (form.data.username) userStore.user!.username = form.data.username;
					if (form.data.email) userStore.user!.email = form.data.email;

					await delay(400);
					isSubmitting = false;
					isSubmitted = true;
					await delay(800);

					isSubmitted = false;
				}
			}
		}
	});

	let isEmpty = $derived(!$form.username && !$form.email);
</script>

<form class="mt-2 flex flex-col gap-y-3" use:enhance>
	<div class="flex max-w-[40.75rem] gap-x-3">
		<div class="flex w-[50%] flex-col gap-y-1">
			<div class="flex items-center gap-x-1">
				<label for="username" class={['w-fit', $errors.username && 'text-red-400']}>Username</label>
				{#if $errors.username}
					<p class="text-red-400">- {$errors.username}</p>
				{/if}
			</div>
			<input
				id="username"
				name="username"
				autocomplete="off"
				type="text"
				class={[
					'bg-main-900 placeholder:text-main-400 hocus:bg-main-800/50 w-full transition-colors duration-100 focus:ring-0',
					$errors.username ? 'border-red-400' : 'border-main-800 hocus:border-main-700'
				]}
				placeholder={userStore.user?.username}
				bind:value={$form.username}
			/>
		</div>

		<div class="flex w-[50%] flex-col gap-y-1">
			<div class="flex items-center gap-x-1">
				<label for="email" class={['w-fit', $errors.email && 'text-red-400']}>Email</label>
				{#if $errors.email}
					<p class="text-red-400">- {$errors.email}</p>
				{/if}
			</div>
			<input
				id="email"
				name="email"
				type="email"
				autocomplete="email"
				class={[
					'bg-main-900 placeholder:text-main-400 hocus:bg-main-800/50 max-w-[20rem] transition-colors duration-100 focus:ring-0',
					$errors.email ? 'border-red-400' : 'border-main-800 hocus:border-main-700'
				]}
				placeholder={userStore.user?.email}
				bind:value={$form.email}
			/>
		</div>
	</div>

	<SubmitButton
		type="submit"
		{buttonWidth}
		{isEmpty}
		{isSubmitting}
		{isSubmitted}
		class="relative mt-2"
	>
		Save your informations
	</SubmitButton>
</form>
