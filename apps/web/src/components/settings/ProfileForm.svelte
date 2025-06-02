<script lang="ts">
	import SubmitButton from 'components/ui/SubmitButton/SubmitButton.svelte';
	import FormInput from '../ui/FormInput/FormInput.svelte';
	import FactsInput from 'components/ui/FactsInput/FactsInput.svelte';
	import LinksInput from 'components/ui/LinksInput/LinksInput.svelte';

	let {
		user,
		about = $bindable(),
		displayName = $bindable(),
		facts = $bindable(),
		links = $bindable(),
		errors = $bindable(),
		enhance,
		isSubmitted = $bindable(),
		isSubmitting = $bindable(),
		buttonWidth
	} = $props();
</script>

<form use:enhance class="flex min-w-[35%] flex-col gap-y-4">
	<FormInput
		type="text"
		title="Display name"
		id="display-name"
		bind:error={errors.display_name}
		bind:inputValue={displayName}
		placeholder={user?.display_name || 'Display name'}
	/>
	<FormInput
		type="rich"
		title="About me"
		id="about-me"
		bind:error={errors.about}
		bind:inputValue={about}
		placeholder="I like potatoes"
	/>

	<LinksInput bind:links />
	<FactsInput bind:facts />

	<SubmitButton type="submit" {buttonWidth} {isSubmitting} {isSubmitted} class="relative mt-2">
		Save my profile
	</SubmitButton>
</form>
