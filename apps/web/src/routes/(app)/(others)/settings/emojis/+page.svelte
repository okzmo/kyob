<script lang="ts">
	import SubmitButton from 'components/ui/SubmitButton/SubmitButton.svelte';
	import type { Emoji } from 'types/types';
	import { generateRandomId } from 'utils/randomId';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { AddEmojisSchema } from 'types/schemas';
	import { backend } from 'stores/backend.svelte';
	import { userStore } from 'stores/user.svelte';
	import { delay } from 'utils/delay';
	import Warning from 'components/ui/icons/Warning.svelte';
	import EmojiLine from 'components/settings/EmojiLine.svelte';
	import { fly } from 'svelte/transition';

	let emojis = $state<Emoji[]>([]);
	let isSubmitting = $state(false);
	let isSubmitted = $state(false);
	let isDeleting = $state(false);
	let buttonWidth = $derived(isSubmitted || isSubmitting ? 40 : 135);

	const { form, errors, enhance } = superForm(defaults(valibot(AddEmojisSchema)), {
		SPA: true,
		dataType: 'json',
		validators: valibot(AddEmojisSchema),
		validationMethod: 'onsubmit',
		resetForm: false,
		async onUpdate({ form }) {
			if (form.valid) {
				isSubmitting = true;
				emojis.forEach((e) => form.data.shortcodes.push(e.shortcode));

				const res = await backend.uploadEmojis(form.data);

				for (const emoji of emojis) {
					URL.revokeObjectURL(emoji.url);
				}

				emojis = [];

				if (res.isErr()) {
					console.error(res.error.error);
					isSubmitting = false;
				}

				if (res.isOk()) {
					const emojis = res.value;
					userStore.emojis.push(...emojis);

					await delay(400);
					isSubmitting = false;
					isSubmitted = true;
					await delay(800);

					isSubmitted = false;
				}
			}
		}
	});

	function onFile(e: Event) {
		const target = e.target as HTMLInputElement;
		const image = target.files?.[0];

		if (image) {
			const dataUrl = URL.createObjectURL(image);

			emojis.push({
				id: generateRandomId(),
				url: dataUrl,
				shortcode: ''
			});

			$form.emojis = [...$form.emojis, image];
		}
	}

	async function onExistingEmojiDelete(id: string) {
		isDeleting = true;
		const res = await backend.deleteEmoji(id);

		if (res.isErr()) {
			console.error(res.error.error);
		}

		if (res.isOk()) {
			userStore.emojis = userStore.emojis.filter((emoji) => emoji.id !== id);
		}
		isDeleting = false;
	}

	function onEmojiDelete(id: string) {
		emojis = emojis.filter((emoji) => emoji.id !== id);
	}
</script>

<h1 class="text-2xl font-bold select-none">Emojis</h1>
<p class="text-main-400 mt-3">
	Add emojis to express yourself! The limit increase as your level increases.
</p>

<h2 class="text-main-400 mt-4 text-sm uppercase">Requirements</h2>
<ul class="text-main-400 mt-1 flex flex-col">
	<li>- File type: JPEG, PNG, GIF, WEBP, AVIF</li>
	<li>- Recommended emoji dimensions: 128x128</li>
</ul>

<div class="mt-4 flex items-end gap-x-3">
	<label
		for="avatar-profile"
		class="group inner-accent/15 hocus:inner-accent-no-shadow/25 bg-accent-100/15 hover:bg-accent-100/25 text-accent-50 relative flex w-fit items-center justify-center overflow-hidden px-2 py-1 whitespace-nowrap transition duration-100"
	>
		<input
			type="file"
			id="avatar-profile"
			name="avatar-profile"
			aria-label="Profile avatar and banner"
			class="absolute h-full w-full text-transparent hover:cursor-pointer"
			accept="image/png, image/jpeg, image/gif, image/webp, image/avif"
			onchange={onFile}
		/>
		<p>Upload an emoji</p>
	</label>
	{#if isDeleting}
		<p class="text-accent-50" transition:fly={{ duration: 100, y: 5 }}>
			Deleting an existing emoji...
		</p>
	{/if}
</div>

<hr class="mt-5 w-full border-none" style="height: 1px; background-color: var(--color-main-800);" />

{#if emojis.length > 0 || userStore.emojis.length > 0}
	<form use:enhance>
		<ul>
			{#each userStore.emojis as emoji (emoji.id)}
				<EmojiLine {...emoji} deleteFunction={onExistingEmojiDelete} />
			{/each}
			{#each emojis as emoji (emoji.id)}
				<EmojiLine {...emoji} deleteFunction={onEmojiDelete} />
			{/each}
		</ul>

		<SubmitButton {isSubmitted} {isSubmitting} {buttonWidth} type="submit" class="relative mt-5">
			Save my emojis
		</SubmitButton>
	</form>

	{#if $errors.emojis || $errors.shortcodes}
		<p class="mt-5 flex items-center gap-x-2 text-red-400">
			<Warning height={20} width={20} />
			{$errors.emojis?.[0] || $errors.shortcodes?.[0]}
		</p>
	{/if}
{/if}
