<script lang="ts">
	import { ContextMenu } from 'bits-ui';
	import UserInvite from '../icons/UserInvite.svelte';
	import Bin from '../icons/Bin.svelte';
	import { serversStore } from '../../../stores/servers.svelte';
	import { userStore } from '../../../stores/user.svelte';
	import { page } from '$app/state';

	let isOwner = $state(
		serversStore.isOwner(userStore.user?.id || -1, Number(page.params.server_id))
	);
	let isMember = $state(serversStore.isMember(Number(page.params.server_id)));

	$effect(() => {
		isOwner = serversStore.isOwner(userStore.user?.id || -1, Number(page.params.server_id));
		isMember = serversStore.isMember(Number(page.params.server_id));
	});
</script>

<ContextMenu.Content class="bg-main-900 border-main-800 w-[225px] rounded-xl border p-2">
	{#if isMember}
		{#if isOwner}
			<ContextMenu.Item
				class="rounded-button data-highlighted:bg-main-800 flex h-10 items-center rounded-lg py-3 pr-1.5 pl-3  font-medium select-none hover:cursor-pointer focus-visible:outline-none"
			>
				<div class="flex w-full items-center justify-between">
					Create Channel
					<UserInvite height={20} width={20} />
				</div>
			</ContextMenu.Item>
		{/if}
	{:else}
		<ContextMenu.Item
			class="rounded-button flex h-10 items-center justify-between rounded-lg py-3 pr-1.5  pl-3 font-medium select-none focus-visible:outline-none"
		>
			<p class="text-main-600 flex items-center">No interactions available</p>
		</ContextMenu.Item>
	{/if}
</ContextMenu.Content>
