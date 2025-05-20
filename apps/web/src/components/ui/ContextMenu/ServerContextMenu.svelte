<script lang="ts">
	import { ContextMenu } from 'bits-ui';
	import UserInvite from '../icons/UserInvite.svelte';
	import LogoutIcon from '../icons/LogoutIcon.svelte';
	import Bin from '../icons/Bin.svelte';
	import { serversStore } from '../../../stores/servers.svelte';
	import { userStore } from '../../../stores/user.svelte';
	import { backend } from '../../../stores/backend.svelte';

	interface Props {
		targetId: number;
	}

	let { targetId }: Props = $props();

	let isOwner = $derived(serversStore.isOwner(userStore.user?.id || -1, targetId));

	async function deleteServer(serverId: number) {
		const res = await backend.deleteServer(serverId);
		if (res.isErr()) {
			console.error(res.error);
		}

		if (res.isOk()) {
			serversStore.removeServer(serverId);
		}
	}

	async function leaveServer(serverId: number) {
		const res = await backend.leaveServer(serverId);
		if (res.isErr()) {
			console.error(res.error);
		}

		if (res.isOk()) {
			serversStore.removeServer(serverId);
		}
	}

	async function createServerInvite(serverId: number) {
		const res = await backend.createInvite(serverId);
		if (res.isErr()) {
			console.error(res.error);
		}

		if (res.isOk()) {
			navigator.clipboard.writeText(res.value);
		}
	}
</script>

<ContextMenu.Item
	class="rounded-button data-highlighted:bg-main-800 flex h-10 items-center justify-between rounded-lg py-3 pr-1.5  pl-3 font-medium select-none hover:cursor-pointer focus-visible:outline-none"
	onclick={() => createServerInvite(targetId)}
>
	<p class="flex items-center">Invite people</p>
	<UserInvite height={20} width={20} />
</ContextMenu.Item>
{#if isOwner}
	<ContextMenu.Item
		class="rounded-button flex h-10 items-center justify-between rounded-lg py-3 pr-1.5 pl-3 font-medium  text-red-400 select-none hover:cursor-pointer focus-visible:outline-none  data-highlighted:bg-red-400/20"
		onclick={() => deleteServer(targetId)}
	>
		<p class="flex items-center">Delete server</p>
		<Bin height={20} width={20} />
	</ContextMenu.Item>
{:else}
	<ContextMenu.Item
		class="rounded-button flex h-10 items-center justify-between rounded-lg py-3 pr-1.5 pl-3 font-medium  text-red-400 select-none hover:cursor-pointer focus-visible:outline-none  data-highlighted:bg-red-400/20"
		onclick={() => leaveServer(targetId)}
	>
		<p class="flex items-center">Leave server</p>
		<LogoutIcon height={20} width={20} />
	</ContextMenu.Item>
{/if}
