<script lang="ts">
	import { ContextMenu } from 'bits-ui';
	import { serversStore } from '../../../../stores/servers.svelte';
	import { userStore } from '../../../../stores/user.svelte';
	import { page } from '$app/state';
	import { core } from '../../../../stores/core.svelte';
	import UserInvite from '../../icons/UserInvite.svelte';

	let isOwner = $derived(
		serversStore.isOwner(userStore.user?.id || -1, Number(page.params.server_id))
	);
</script>

<ContextMenu.Content class="bg-main-900 border-main-800 w-[225px] rounded-xl border p-2">
	{#if isOwner}
		<ContextMenu.Item
			class="rounded-button data-highlighted:bg-main-800 flex h-10 items-center rounded-lg py-3 pr-1.5 pl-3  font-medium select-none hover:cursor-pointer focus-visible:outline-none"
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
</ContextMenu.Content>
