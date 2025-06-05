<script lang="ts">
	import { userStore } from 'stores/user.svelte';
	import SettingsSection from 'components/ui/SettingsSection/SettingsSection.svelte';
	import UserProfileSettings from 'components/UserProfile/UserProfileSettings.svelte';
	import ProfileForm from 'components/settings/ProfileForm.svelte';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { UpdateProfileSchema } from 'types/schemas';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { backend } from 'stores/backend.svelte';
	import { delay } from 'utils/delay';
	import { generateText } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';

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
		resetForm: false,
		async onUpdate({ form }) {
			if (form.valid) {
				isSubmitting = true;

				const res = await backend.updateProfile({
					display_name: form.data.display_name?.trim(),
					about: form.data.about,
					facts: form.data.facts,
					links: form.data.links
				});

				if (res.isErr()) {
					if (res.error.code === 'ERR_UNKNOWN') {
						console.log(res.error.error);
					}
					isSubmitting = false;
					return;
				}

				if (res.isOk()) {
					await delay(500);
					isSubmitting = false;
					isSubmitted = true;
					await delay(1000);

					isSubmitted = false;

					if (form.data.display_name) userStore.user!.display_name = form.data.display_name;
					if (form.data.about) userStore.user!.about = form.data.about;
					if (form.data.links) userStore.user!.links = form.data.links;
					if (form.data.facts) userStore.user!.facts = form.data.facts;
				}
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

	function checkChanges() {
		if (!user) return true;

		if ($form.display_name !== user?.display_name) return false;

		if ($form.about && !user.about) return false;

		if ($form.about && user.about) {
			const formAbout = generateText($form.about, [
				StarterKit.configure({
					gapcursor: false,
					dropcursor: false,
					heading: false,
					orderedList: false,
					bulletList: false,
					blockquote: false
				})
			]);
			const userAbout = generateText(user.about || {}, [
				StarterKit.configure({
					gapcursor: false,
					dropcursor: false,
					heading: false,
					orderedList: false,
					bulletList: false,
					blockquote: false
				})
			]);

			if (formAbout !== userAbout) return false;
		}

		if ($form.links.length > user.links.length || $form.links.length < user.links.length)
			return false;
		if (
			$form.links.some(
				(link, idx) => link.label !== user.links[idx].label || link.url !== user.links[idx].url
			)
		)
			return false;

		if ($form.facts.length > user.facts.length || $form.facts.length < user.facts.length)
			return false;
		if (
			$form.facts.some(
				(link, idx) => link.label !== user.facts[idx].label || link.value !== user.facts[idx].value
			)
		)
			return false;

		return true;
	}

	let isEmpty = $derived.by(checkChanges);
</script>

<h1 class="text-2xl font-bold select-none">Profile</h1>

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
			{isEmpty}
			{buttonWidth}
		/>

		{#if user}
			<UserProfileSettings
				{user}
				bind:about={$form.about}
				bind:displayName={$form.display_name}
				bind:links={$form.links!}
				bind:facts={$form.facts!}
			/>
		{/if}
	</div>
</SettingsSection>
