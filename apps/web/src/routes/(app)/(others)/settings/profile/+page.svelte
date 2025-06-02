<script lang="ts">
	import { userStore } from 'stores/user.svelte';
	import SettingsSection from 'components/ui/SettingsSection/SettingsSection.svelte';
	import UserProfileSettings from 'components/UserProfile/UserProfileSettings.svelte';
	import ProfileForm from 'components/settings/ProfileForm.svelte';
	import { defaults, setError, superForm } from 'sveltekit-superforms';
	import { UpdateProfileSchema } from 'types/schemas';
	import { valibot } from 'sveltekit-superforms/adapters';

	const user = $derived(userStore.user);

	let isSubmitting = $state(false);
	let isSubmitted = $state(false);
	let buttonWidth = $derived(isSubmitted || isSubmitting ? 40 : 135);

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

	$effect(() => {
		document.documentElement.style.setProperty(
			'--user-color-85',
			`rgba(${userStore.user?.main_color}, 0.85)`
		);
		document.documentElement.style.setProperty(
			'--user-color-95',
			`rgba(${userStore.user?.main_color}, 0.95)`
		);
		document.documentElement.style.setProperty(
			'--user-color',
			`rgba(${userStore.user?.main_color}, 1)`
		);
	});

	$effect(() => {
		$form.about = user?.about || '';
		$form.display_name = user?.display_name || '';
		$form.links = user?.links || [];
		$form.facts = user?.facts || [];
	});

	$inspect($form.links);
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
			bind:isSubmitted
			bind:isSubmitting
			bind:facts={$form.facts}
			bind:links={$form.links}
			{buttonWidth}
		/>

		{#if user}
			<UserProfileSettings {user} bind:about={$form.about} bind:displayName={$form.display_name} />
		{/if}
	</div>
</SettingsSection>
