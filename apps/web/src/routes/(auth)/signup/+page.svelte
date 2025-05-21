<script lang="ts">
	import { valibot } from 'sveltekit-superforms/adapters';
	import AuthForm from '../../../components/auth/AuthForm.svelte';
	import { SignUpSchema } from '../../../types/schemas';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { goto } from '$app/navigation';

	let globalError = $state<string | undefined>();

	const { form, errors, enhance } = superForm(defaults(valibot(SignUpSchema)), {
		SPA: true,
		validators: valibot(SignUpSchema),
		async onUpdate({ form }) {
			if (form.valid) {
				try {
					const res = await fetch(`${import.meta.env.VITE_API_URL}/signup`, {
						method: 'post',
						credentials: 'include',
						headers: {
							'Content-Type': 'application/json'
						},
						body: JSON.stringify({
							email: form.data.email,
							username: form.data.username,
							display_name: form.data.display_name,
							password: form.data.password
						})
					});

					if (!res.ok) {
						const data = await res.json();
						console.error('signup failed', res.status, data);
						return;
					}

					return goto('/');
				} catch (err) {
					console.error(err);
					globalError = 'Signup failed';
				}
			}
		}
	});
</script>

<AuthForm type="signup" {form} {errors} {enhance} {globalError} />
