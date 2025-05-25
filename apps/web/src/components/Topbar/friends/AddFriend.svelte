<script lang="ts">
	import { core } from '../../../stores/core.svelte';
	import PlusSimple from '../../ui/icons/PlusSimple.svelte';
</script>

<button
	class="bg-accent-100/80 text-main-50 hocus:bg-accent-100/60 mt-4 flex w-full items-center justify-center gap-x-2 rounded-xl py-2 font-semibold transition-colors duration-100 hover:cursor-pointer"
	onclick={() => (core.openAddFriendModal.status = true)}
>
	Add friend
	<PlusSimple height={16} width={16} />
</button>
