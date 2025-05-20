<script lang="ts">
	import { ContextMenu } from 'bits-ui';
	import { core } from '../../../../stores/core.svelte';
	import Plus from '../../icons/Plus.svelte';
	import JoinIcon from '../../icons/JoinIcon.svelte';
</script>

<ContextMenu.Item
	class="rounded-button data-highlighted:bg-main-800 flex h-10 items-center rounded-lg py-3 pr-1.5 pl-3  font-medium select-none hover:cursor-pointer focus-visible:outline-none"
	onclick={() => {
		core.deactivateMapDragging();
		core.openCreateServerModal.status = true;
	}}
>
	<div class="flex w-full items-center justify-between">
		Create realm
		<Plus height={20} width={20} />
	</div>
</ContextMenu.Item>
<ContextMenu.Item
	class="rounded-button data-highlighted:bg-main-800 flex h-10 items-center rounded-lg py-3 pr-1.5 pl-3  font-medium select-none hover:cursor-pointer focus-visible:outline-none"
	onclick={() => {
		core.deactivateMapDragging();
		core.openJoinServerModal.status = true;
	}}
>
	<div class="flex w-full items-center justify-between">
		Join realm
		<JoinIcon height={20} width={20} />
	</div>
</ContextMenu.Item>
