<script lang="ts">
	import Corners from '../ui/Corners/Corners.svelte';
	import MagnifyingGlass from '../ui/icons/MagnifyingGlass.svelte';
</script>

<div
	class="bg-main-900/80 inner-shadow-input group fixed bottom-5 left-1/2 z-50 flex w-[18.5rem] -translate-x-1/2 items-center px-4 py-1 backdrop-blur-2xl transition-all duration-100 focus-within:w-[22rem]"
>
	<Corners color="border-main-700" class="group-hocus:border-main-300" />
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
