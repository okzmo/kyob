<script lang="ts">
	import { valibot } from 'sveltekit-superforms/adapters';
	import AuthForm from '../../../components/auth/AuthForm.svelte';
	import { SignInSchema } from '../../../types/schemas';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { goto } from '$app/navigation';

	let globalError = $state<string | undefined>();

	const { form, errors, enhance } = superForm(defaults(valibot(SignInSchema)), {
		SPA: true,
		validators: valibot(SignInSchema),
		async onUpdate({ form }) {
			if (form.valid) {
				try {
					const res = await fetch(`${import.meta.env.VITE_API_URL}/signin`, {
						method: 'post',
						credentials: 'include',
						headers: {
							'Content-Type': 'application/json'
						},
						body: JSON.stringify({
							email_or_username: form.data.username,
							password: form.data.password
						})
					});

					if (!res.ok) {
						const data = await res.json();
						console.error('signin failed', res.status, data);
						globalError = data.error;
						return;
					}

					return goto('/');
				} catch (err) {
					console.error(err);
					globalError = 'Signin failed';
				}
			}
		}
	});
</script>

<AuthForm type="signin" {form} {errors} {enhance} bind:globalError />
