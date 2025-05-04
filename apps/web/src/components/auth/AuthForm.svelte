<script lang="ts">
	import { page } from '$app/state';
	import type { SuperFormData, SuperFormErrors } from 'sveltekit-superforms/client';
	import Envelope from '../ui/icons/Envelope.svelte';
	import KyobIcon from '../ui/icons/KyobIcon.svelte';
	import Lock from '../ui/icons/Lock.svelte';
	import User from '../ui/icons/User.svelte';

	interface Props {
		type: 'signin' | 'signup';
		form: SuperFormData<any>;
		errors: SuperFormErrors<any>;
		enhance: any;
		globalError?: string;
	}

	let { type, form, errors, enhance, globalError }: Props = $props();
</script>

{#if globalError}
	<p
		class="fixed top-[22%] left-1/2 mt-4 w-[25rem] -translate-x-1/2 rounded-lg border border-red-400 bg-red-400/10 py-2 text-center text-red-400"
	>
		{globalError}
	</p>
{/if}
<div
	class="bg-main-900 border-main-800 fixed top-[29%] left-1/2 flex w-[25rem] -translate-x-1/2 flex-col items-center rounded-2xl border px-4 pt-8 pb-4"
>
	<KyobIcon />

	<nav class="mt-8">
		<ul class="bg-main-800 border-main-600 flex w-fit items-center gap-x-1 rounded-lg border p-1">
			<li>
				<a
					class={[
						'flex items-center justify-center rounded-md px-8 py-1',
						page.url.pathname === '/signup' && 'bg-main-600/50'
					]}
					href="/signup">Sign Up</a
				>
			</li>
			<li>
				<a
					class={[
						'flex items-center justify-center rounded-md px-8 py-1',
						page.url.pathname === '/signin' && 'bg-main-600/50'
					]}
					href="/signin"
				>
					Sign In
				</a>
			</li>
		</ul>
	</nav>

	<form method="POST" class="mt-8 flex w-full flex-col gap-y-5" use:enhance>
		<div class="flex flex-col gap-y-2.5">
			{#if type === 'signup'}
				<div
					class={[
						' flex w-full items-center rounded-xl border pl-3.5 transition-colors duration-100',
						$errors.email ? 'border-red-400' : 'border-main-800 focus-within:border-main-500'
					]}
				>
					<Envelope width={20} height={20} class="text-main-500" />
					<input
						bind:value={$form.email}
						class="placeholder:text-main-500 w-full border-none bg-transparent py-3 focus:ring-0"
						type="email"
						autocomplete="email"
						placeholder="Email Address"
					/>
				</div>
				{#if $errors.email}
					<p class="text-sm leading-none text-red-400">{$errors.email?.[0]}</p>
				{/if}
			{/if}

			<div
				class={[
					' flex w-full items-center rounded-xl border pl-3.5 transition-colors duration-100',
					$errors.username ? 'border-red-400' : 'border-main-800 focus-within:border-main-500'
				]}
			>
				<User width={20} height={20} class="text-main-500" />
				<input
					bind:value={$form.username}
					class="placeholder:text-main-500 w-full border-none bg-transparent py-3 focus:ring-0"
					type="text"
					placeholder={type === 'signin' ? 'Email or username' : 'Username'}
				/>
			</div>
			{#if $errors.username}
				<p class="text-sm leading-none text-red-400">{$errors.username?.[0]}</p>
			{/if}

			{#if type === 'signup'}
				<div
					class={[
						' flex w-full items-center rounded-xl border pl-3.5 transition-colors duration-100',
						$errors.display_name ? 'border-red-400' : 'border-main-800 focus-within:border-main-500'
					]}
				>
					<User width={20} height={20} class="text-main-500" />
					<input
						bind:value={$form.display_name}
						class="placeholder:text-main-500 w-full border-none bg-transparent py-3 focus:ring-0"
						type="text"
						placeholder="Display name"
					/>
				</div>
				{#if $errors.display_name}
					<p class="text-sm leading-none text-red-400">{$errors.display_name?.[0]}</p>
				{/if}
			{/if}

			<div
				class={[
					'flex w-full items-center rounded-xl border pl-3.5 transition-colors duration-100',
					$errors.password ? 'border-red-400' : 'border-main-800 focus-within:border-main-500 '
				]}
			>
				<Lock width={20} height={20} class="text-main-500" />
				<input
					class="placeholder:text-main-500 w-full border-none bg-transparent py-3 focus:ring-0"
					bind:value={$form.password}
					type="password"
					autocomplete="new-password"
					placeholder="Password"
				/>
			</div>
			{#if $errors.password}
				<p class="text-sm leading-none text-red-400">{$errors.password?.[0]}</p>
			{/if}
		</div>
		{#if type === 'signin'}
			<button
				type="submit"
				class="bg-main-800 hocus:bg-accent-100/35 hocus:text-accent-50 w-full rounded-xl py-3 transition-colors duration-100 hover:cursor-pointer"
				>Sign in</button
			>
		{:else}
			<button
				type="submit"
				class="bg-main-800 hocus:bg-accent-100/35 hocus:text-accent-50 w-full rounded-xl py-3 transition-colors duration-100 hover:cursor-pointer"
				>Create an account</button
			>
		{/if}
	</form>
</div>
