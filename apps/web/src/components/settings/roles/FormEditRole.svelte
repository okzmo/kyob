<script lang="ts">
	import { defaults, superForm } from 'sveltekit-superforms';
	import { UpdateServerProfileSchema } from 'types/schemas';
	import { valibot } from 'sveltekit-superforms/adapters';

	let { activeTab } = $props();

	const { form, errors, enhance } = superForm(defaults(valibot(UpdateServerProfileSchema)), {
		SPA: true,
		dataType: 'json',
		validators: valibot(UpdateServerProfileSchema),
		validationMethod: 'onsubmit',
		resetForm: false,
		async onUpdate({ form }) {
			if (form.valid) {
			}
		}
	});
</script>

<form use:enhance>
	{#if activeTab === 'display'}
		display
	{:else if activeTab === 'permissions'}
		permissions
	{/if}
</form>
