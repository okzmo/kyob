<script lang="ts">
	import Close from 'components/ui/icons/Close.svelte';

	interface Props {
		id: string;
		url: string;
		shortcode: string;
		deleteFunction: (id: string) => void;
	}

	let { id, url, shortcode, deleteFunction }: Props = $props();
</script>

<li
	class="group hocus:bg-main-800 hocus:border-main-700 focus-within:bg-main-800 focus-within:border-main-700 border-b-main-800 flex items-center justify-between border border-t-transparent border-r-transparent border-l-transparent p-4 transition-colors"
>
	<div class="flex items-center">
		<img src={url} alt={shortcode} class="h-[32px] w-[32px] object-contain" />
		<div
			class="hover:bg-main-900 hover:border-main-700 focus-within:bg-main-900 focus-within:border-main-700 flex h-[32px] items-center border border-transparent p-2 transition-colors"
		>
			<span class="text-main-50/25 font-bold">:</span>
			<input
				type="text"
				bind:value={shortcode}
				class="placeholder:text-main-50/50 h-full border-none bg-transparent px-1 py-0 focus-visible:ring-0 focus-visible:outline-none"
				placeholder="emoji_name"
			/>
			<span class="text-main-50/25 font-bold">:</span>
		</div>
	</div>

	<button
		type="button"
		class="inner-red-400/20 hocus:inner-red-400/40 hocus:bg-red-400/25 group-hocus:opacity-100 bg-red-400/15 p-1 text-red-400 opacity-0 transition hover:cursor-pointer"
		aria-label="Delete emoji"
		onclick={() => deleteFunction(id)}
	>
		<Close height={16} width={16} />
	</button>
</li>
