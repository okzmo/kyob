<script lang="ts">
	import { userStore } from '../../../../stores/user.svelte';
	import SettingsSection from '../../../../components/ui/SettingsSection/SettingsSection.svelte';
	import { onMount } from 'svelte';
	import UserProfileSettings from '../../../../components/UserProfile/UserProfileSettings.svelte';
	import ProfileForm from '../../../../components/settings/ProfileForm.svelte';
	import { defaults, setError, superForm } from 'sveltekit-superforms';
	import { UpdateProfileSchema } from '../../../../types/schemas';
	import { valibot } from 'sveltekit-superforms/adapters';

	const user = $derived(userStore.user);

	let isSubmitting = $state(false);
	let isSubmitted = $state(false);
	let buttonWidth = $derived(isSubmitted || isSubmitting ? 40 : 184);

	const { form, errors, enhance } = superForm(defaults(valibot(UpdateProfileSchema)), {
		SPA: true,
		dataType: 'json',
		validators: valibot(UpdateProfileSchema),
		validationMethod: 'onsubmit',
		multipleSubmits: 'prevent',
		async onUpdate({ form }) {
			if (form.valid) {
			}
		}
	});

	onMount(() => {
		document.documentElement.style.setProperty('--user-color-85', '#153c45d9');
		document.documentElement.style.setProperty('--user-color-95', '#153c45f2');
		document.documentElement.style.setProperty('--user-color', '#153c45');
	});

	$effect(() => {
		$form.about = user?.about || '';
		$form.display_name = user?.display_name || '';
	});
</script>

<h1 class="text-2xl font-bold">Profile</h1>

<SettingsSection
	title="Personal profile"
	description="This is what everyone see when clicking on your avatar or display name."
	class="mt-10"
>
	<div class="mt-5 flex gap-x-15">
		<ProfileForm
			{user}
			{enhance}
			bind:about={$form.about}
			bind:displayName={$form.display_name}
			bind:errors={$errors}
		/>

		{#if user}
			<UserProfileSettings {user} bind:about={$form.about} bind:displayName={$form.display_name} />
		{/if}
	</div>
</SettingsSection>
