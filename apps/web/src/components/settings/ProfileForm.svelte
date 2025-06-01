<script lang="ts">
	import FormInput from '../ui/FormInput/FormInput.svelte';

	let {
		user,
		about = $bindable(),
		displayName = $bindable(),
		errors = $bindable(),
		enhance
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
</form>
