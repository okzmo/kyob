<script lang="ts">
	import { core } from 'stores/core.svelte';
	import Corners from '../../ui/Corners/Corners.svelte';
	import PlusSimple from '../../ui/icons/PlusSimple.svelte';

	let { isOpen = $bindable() } = $props();
</script>

<button
	class="group bg-accent-100/15 hocus:bg-accent-100/25 text-accent-50 inner-accent/15 hocus:inner-accent/25 relative mt-4 flex w-full items-center justify-center gap-x-2 py-2 font-semibold transition duration-100 hover:cursor-pointer"
	onclick={() => {
		core.openAddFriendModal.status = true;
		isOpen = false;
	}}
>
	Add friend
	<PlusSimple height={16} width={16} />
	<Corners color="border-accent-200" class="group-hocus:border-accent-100" />
</button>
