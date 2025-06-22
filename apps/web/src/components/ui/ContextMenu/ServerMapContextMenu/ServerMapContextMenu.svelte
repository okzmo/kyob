<script lang="ts">
	import { ContextMenu } from 'bits-ui';
	import { core } from 'stores/core.svelte';
	import Plus from '../../icons/Plus.svelte';
	import JoinIcon from '../../icons/JoinIcon.svelte';
</script>

<ContextMenu.Item
	class="context-menu-item"
	onclick={() => {
		core.deactivateMapDragging();
		core.createServerModal.status = true;
	}}
>
	<div class="flex w-full items-center justify-between">
		Create realm
		<Plus height={20} width={20} />
	</div>
</ContextMenu.Item>
<ContextMenu.Item
	class="context-menu-item"
	onclick={() => {
		core.deactivateMapDragging();
		core.joinServerModal.status = true;
	}}
>
	<div class="flex w-full items-center justify-between">
		Join realm
		<JoinIcon height={20} width={20} />
	</div>
</ContextMenu.Item>
