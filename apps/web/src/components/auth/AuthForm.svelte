<script lang="ts">
	import { page } from '$app/state';
	import { type ZodSafeParseResult } from 'zod';
	import Envelope from '../ui/icons/Envelope.svelte';
	import KyobIcon from '../ui/icons/KyobIcon.svelte';
	import Lock from '../ui/icons/Lock.svelte';
	import User from '../ui/icons/User.svelte';
	import { UserSchema } from '../../types/types';

	interface Props {
		type: 'signin' | 'signup';
	}

	let { type }: Props = $props();

	let form = $state({
		email: '',
		username: '',
		password: ''
	});

	let validationSignIn = $state<ZodSafeParseResult<{ email: string; password: string }>>();
	let validationSignUp =
		$state<ZodSafeParseResult<{ email: string; username: string; password: string }>>();

	function signIn(e: SubmitEvent) {
		e.preventDefault();
		validationSignIn = UserSchema.pick({ email: true, password: true }).safeParse({
			email: form.email,
			password: form.password
		});
	}

	function signUp(e: SubmitEvent) {
		e.preventDefault();
		validationSignUp = UserSchema.pick({ email: true, username: true, password: true }).safeParse({
			email: form.email,
			password: form.password
		});
	}
</script>

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

	<form
		method="POST"
		class="mt-8 flex w-full flex-col gap-y-5"
		onsubmit={type === 'signin' ? signIn : signUp}
	>
		<div class="flex flex-col gap-y-2.5">
			<div
				class="border-main-800 focus-within:border-main-500 flex w-full items-center rounded-xl border pl-3.5 transition-colors duration-100"
			>
				<Envelope width={20} height={20} class="text-main-500" />
				<input
					class="placeholder:text-main-500 w-full border-none bg-transparent py-3 focus:ring-0"
					type="email"
					autocomplete="email"
					placeholder="Email Address"
				/>
			</div>
			{#if type === 'signup'}
				<div
					class="border-main-800 focus-within:border-main-500 flex w-full items-center rounded-xl border pl-3.5 transition-colors duration-100"
				>
					<User width={20} height={20} class="text-main-500" />
					<input
						class="placeholder:text-main-500 w-full border-none bg-transparent py-3 focus:ring-0"
						type="text"
						placeholder="Username"
					/>
				</div>
			{/if}
			<div
				class="border-main-800 focus-within:border-main-500 flex w-full items-center rounded-xl border pl-3.5 transition-colors duration-100"
			>
				<Lock width={20} height={20} class="text-main-500" />
				<input
					class="placeholder:text-main-500 w-full border-none bg-transparent py-3 focus:ring-0"
					type="password"
					autocomplete="new-password"
					placeholder="Password"
				/>
			</div>
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
