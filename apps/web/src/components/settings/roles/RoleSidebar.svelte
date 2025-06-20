<script lang="ts">
	import PlusSimple from 'components/ui/icons/PlusSimple.svelte';

	let { activeRole = $bindable(), creatingRole = $bindable() } = $props();
</script>

<ul class="flex w-[15rem] flex-col gap-y-2">
	<li>
		<button
			class="border-main-400 bg-main-400/15 text-main-400 hocus:bg-main-400/25 flex w-full items-center justify-center gap-x-2 border px-2 py-2 font-semibold transition-colors duration-100 hover:cursor-pointer"
			onclick={() => (creatingRole = true)}
		>
			<PlusSimple height={20} width={20} />
		</button>
	</li>
</ul>
