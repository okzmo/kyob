<script lang="ts">
	import { ContextMenu } from 'bits-ui';
	import { serversStore } from 'stores/servers.svelte';
	import { userStore } from 'stores/user.svelte';
	import { page } from '$app/state';
	import { core } from 'stores/core.svelte';
	import UserInvite from '../../icons/UserInvite.svelte';

	let isOwner = $derived(serversStore.isOwner(userStore.user?.id || '', page.params.server_id));
</script>

{#if isOwner}
	<ContextMenu.Item
		class="context-menu-item"
		onclick={() => {
			core.deactivateMapDragging();
			core.openCreateChannelModal.status = true;
		}}
	>
		<div class="flex w-full items-center justify-between">
			Create Channel
			<UserInvite height={20} width={20} />
		</div>
	</ContextMenu.Item>
{/if}
