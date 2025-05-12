<script lang="ts">
	import MagnifyingGlass from '../ui/icons/MagnifyingGlass.svelte';
</script>

<div
	class="bg-main-900/80 group border-main-800 focus-within:border-main-500 fixed bottom-5 left-1/2 z-50 flex w-[18.5rem] -translate-x-1/2 items-center rounded-2xl border px-4 py-1 backdrop-blur-2xl transition-[width,color,border-color] duration-100 focus-within:w-[22rem]"
>
	<MagnifyingGlass
		height={20}
		width={20}
		class="text-main-500 group-focus-within:text-main-50 transition-colors duration-100"
	/>
	<input
		placeholder="What should we search?"
		type="text"
		autocomplete="off"
		class="text-main-50 placeholder:text-main-500 w-full border-none bg-transparent focus:ring-0"
	/>
</div>
