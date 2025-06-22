<script lang="ts">
	import SubmitButton from 'components/ui/SubmitButton/SubmitButton.svelte';
	import FormInput from '../ui/FormInput/FormInput.svelte';

	let {
		server,
		serverName = $bindable(),
		about = $bindable(),
		errors = $bindable(),
		enhance,
		isSubmitted = $bindable(),
		isSubmitting = $bindable(),
		isEmpty,
		buttonWidth
	} = $props();
</script>

<form use:enhance class="flex min-w-[35%] flex-col gap-y-4 select-none">
	<FormInput
		type="text"
		title="Name"
		id="server-name"
		bind:error={errors.display_name}
		bind:inputValue={serverName}
		placeholder={server?.name || 'Server name'}
	/>
	<FormInput
		type="rich"
		title="Description"
		id="description"
		bind:error={errors.about}
		bind:inputValue={about}
		placeholder="This community is about..."
		inputClass="max-w-[21.5rem]"
	/>

	<SubmitButton
		type="submit"
		{buttonWidth}
		{isSubmitting}
		{isSubmitted}
		{isEmpty}
		class="relative mt-2"
	>
		Save profile
	</SubmitButton>
</form>
