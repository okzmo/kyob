<script lang="ts">
	import PlusSimple from '../icons/PlusSimple.svelte';
</script>

<button
	class="button fixed bottom-20 left-1/2 -translate-x-1/2 rounded-xl bg-gradient-to-r from-indigo-400 to-rose-500 py-1.5 pr-4 pl-3.5 hover:cursor-pointer"
>
	<span
		class="flex items-center gap-x-2 bg-gradient-to-r from-indigo-400 to-rose-500 bg-clip-text text-transparent"
	>
		<PlusSimple height={16} width={16} class="text-indigo-400" />
		Create a realm
	</span>
</button>

<style>
	.button::before {
		content: '';
		position: absolute;
		width: calc(100% - 1px);
		height: calc(100% - 1px);
		background-color: var(--ui-main-900);
		left: 50%;
		top: 50%;
		transform: translate(-50%, -50%);
		z-index: -1;
		border-radius: inherit;
		transition: opacity 100ms ease-out;
	}

	.button:hover::before {
		opacity: 92%;
	}
</style>
