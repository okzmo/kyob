<script lang="ts">
	import { goto } from '$app/navigation';
	import Corners from '../../ui/Corners/Corners.svelte';
	import ArrowLeft from '../../ui/icons/ArrowLeft.svelte';
</script>

<button
	class="top-bar-button text-main-400 group hocus:bg-accent-100/15 relative flex items-center gap-x-4 py-1.5 pr-4 pl-3 transition duration-100 select-none hover:cursor-pointer"
	onclick={() => goto('/')}
>
	<Corners color="border-main-300" class="group-hocus:border-accent-100 duration-100" />
	<ArrowLeft
		height={22}
		width={22}
		class="group-hocus:text-accent-50 transition-colors duration-100"
	/>
	<p class="group-hocus:text-accent-50 text-lg transition-colors duration-100">Go back</p>
</button>
