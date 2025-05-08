<script lang="ts">
	import { ContextMenu } from 'bits-ui';
	import UserInvite from '../icons/UserInvite.svelte';
	import Bin from '../icons/Bin.svelte';
	import { serversStore } from '../../../stores/servers.svelte';
	import { userStore } from '../../../stores/user.svelte';
	import { page } from '$app/state';
	import { backend } from '../../../stores/backend.svelte';

	interface Props {
		targetId: number;
	}

	let { targetId }: Props = $props();

	let isOwner = $derived(
		serversStore.isOwner(userStore.user?.id || -1, Number(page.params.server_id))
	);
	let isMember = $derived(serversStore.isMember(Number(page.params.server_id)));

	async function deleteChannel(serverId: number, channelId: number) {
		const res = await backend.deleteChannel(serverId, channelId);
		if (res.isErr()) {
			console.error(res.error);
		}

		if (res.isOk()) {
			serversStore.removeChannel(serverId, channelId);
		}
	}
</script>

<ContextMenu.Content class="bg-main-900 border-main-800 w-[225px] rounded-xl border p-2">
	{#if isMember}
		<ContextMenu.Item
			class="rounded-button data-highlighted:bg-main-800 flex h-10 items-center rounded-lg py-3 pr-1.5 pl-3  font-medium select-none hover:cursor-pointer focus-visible:outline-none"
		>
			<div class="flex w-full items-center justify-between">Mark as read</div>
		</ContextMenu.Item>
		{#if isOwner}
			<ContextMenu.Item
				class="rounded-button data-highlighted:bg-main-800 flex h-10 items-center rounded-lg py-3 pr-1.5 pl-3  font-medium select-none hover:cursor-pointer focus-visible:outline-none"
			>
				<div class="flex w-full items-center justify-between">
					Edit Channel
					<UserInvite height={20} width={20} />
				</div>
			</ContextMenu.Item>
			<ContextMenu.Item
				class="rounded-button flex h-10 items-center justify-between rounded-lg py-3 pr-1.5 pl-3 font-medium  text-red-400 select-none hover:cursor-pointer focus-visible:outline-none  data-highlighted:bg-red-400/20"
				onclick={() => deleteChannel(Number(page.params.server_id), targetId)}
			>
				<p class="flex items-center">Delete Channel</p>
				<Bin height={20} width={20} />
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
