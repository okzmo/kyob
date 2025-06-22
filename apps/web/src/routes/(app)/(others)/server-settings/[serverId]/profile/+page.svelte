<script lang="ts">
	import { defaults, superForm } from 'sveltekit-superforms';
	import { UpdateServerProfileSchema } from 'types/schemas';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { backend } from 'stores/backend.svelte';
	import { delay } from 'utils/time';
	import { generateText } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import ServerProfileForm from 'components/settings/ServerProfileForm.svelte';
	import { serversStore } from 'stores/servers.svelte';
	import { page } from '$app/state';
	import ServerProfileSettings from 'components/ServerProfile/ServerProfileSettings.svelte';

	const server = $derived(serversStore.getServer(page.params.serverId));

	let isSubmitting = $state(false);
	let isSubmitted = $state(false);
	let buttonWidth = $derived(isSubmitted || isSubmitting ? 40 : 110);

	const { form, errors, enhance } = superForm(defaults(valibot(UpdateServerProfileSchema)), {
		SPA: true,
		dataType: 'json',
		validators: valibot(UpdateServerProfileSchema),
		validationMethod: 'onsubmit',
		resetForm: false,
		async onUpdate({ form }) {
			if (form.valid) {
				isSubmitting = true;

				const res = await backend.updateServerProfile(server.id, {
					name: form.data.name,
					description: form.data.description
				});

				if (res.isErr()) {
					if (res.error.code === 'ERR_UNKNOWN') {
						console.error(res.error.error);
					}
					isSubmitting = false;
					return;
				}

				if (res.isOk()) {
					await delay(400);
					isSubmitting = false;
					isSubmitted = true;
					await delay(800);

					isSubmitted = false;

					const serverEl = serversStore.getServer(server.id);
					if (form.data.name) serverEl.name = form.data.name;
					if (form.data.description) serverEl.description = form.data.description;
				}
			}
		}
	});

	$effect(() => {
		document.documentElement.style.setProperty(
			'--server-color-85',
			`rgba(${server?.main_color}, 0.85)`
		);
		document.documentElement.style.setProperty(
			'--server-color-95',
			`rgba(${server?.main_color}, 0.95)`
		);
		document.documentElement.style.setProperty('--server-color', `rgba(${server?.main_color}, 1)`);
	});

	$effect(() => {
		$form.description = server?.description || '';
		$form.name = server?.name || '';
	});

	function checkChanges() {
		if (!server) return true;

		if ($form.name !== server?.name) return false;

		if ($form.description && !server.description) return false;

		if ($form.description && server.description) {
			const formDescription = generateText($form.description, [
				StarterKit.configure({
					gapcursor: false,
					dropcursor: false,
					heading: false,
					orderedList: false,
					bulletList: false,
					blockquote: false
				})
			]);
			const serverDescription = generateText(server.description || {}, [
				StarterKit.configure({
					gapcursor: false,
					dropcursor: false,
					heading: false,
					orderedList: false,
					bulletList: false,
					blockquote: false
				})
			]);

			if (formDescription !== serverDescription) return false;
		}

		return true;
	}

	let isEmpty = $derived.by(checkChanges);
</script>

<h1 class="text-2xl font-bold select-none">Server Profile</h1>

<div class="mt-10 flex gap-x-15">
	<ServerProfileForm
		{server}
		{enhance}
		bind:serverName={$form.name}
		bind:about={$form.description}
		bind:errors={$errors}
		bind:isSubmitted
		bind:isSubmitting
		{isEmpty}
		{buttonWidth}
	/>

	{#if server}
		<ServerProfileSettings
			{server}
			bind:serverName={$form.name}
			bind:description={$form.description}
		/>
	{/if}
</div>
