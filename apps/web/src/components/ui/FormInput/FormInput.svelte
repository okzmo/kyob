<script lang="ts">
	import AboutInput from '../../settings/AboutInput.svelte';

	interface Props {
		id: string;
		inputValue: any;
		error: any;
		placeholder: string;
		title: string;
		type: 'text' | 'password' | 'textarea' | 'rich';
		class?: string;
	}

	let {
		id,
		inputValue = $bindable(),
		error = $bindable(),
		placeholder,
		title,
		type,
		class: classes
	}: Props = $props();
</script>

<div class={['flex flex-col', classes]}>
	<div class="flex items-center gap-x-1">
		<label for={id} class={['text-sm', error ? 'text-red-400' : 'text-main-500']}>{title}</label>
		{#if error}
			<p class="text-sm text-red-400">- {error}</p>
		{/if}
	</div>
	{#if type !== 'textarea' && type !== 'rich'}
		<input
			{id}
			{type}
			bind:value={inputValue}
			{placeholder}
			class={[
				'bg-main-900 placeholder:text-main-400 hocus:bg-main-800/50 mt-1.5 transition-colors duration-100 focus:ring-0',
				error ? 'border-red-400' : 'border-main-800 hocus:border-main-700'
			]}
		/>
	{:else if type === 'textarea'}
		<textarea
			{id}
			bind:value={inputValue}
			{placeholder}
			class={[
				'bg-main-900 placeholder:text-main-400 hocus:bg-main-800/50 mt-1.5 h-30 resize-none transition-colors duration-100 focus:ring-0',
				error ? 'border-red-400' : 'border-main-800 hocus:border-main-700'
			]}
		></textarea>
	{:else if type === 'rich'}
		<AboutInput bind:content={inputValue} {placeholder} />
	{/if}
</div>
