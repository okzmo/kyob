<script lang="ts">
	import { ContextMenu } from 'bits-ui';
	import UserInvite from '../icons/UserInvite.svelte';
	import LogoutIcon from '../icons/LogoutIcon.svelte';
	import Bin from '../icons/Bin.svelte';
	import { serversStore } from '../../../stores/servers.svelte';
	import { userStore } from '../../../stores/user.svelte';
	import { backend } from '../../../stores/backend.svelte';

	interface Props {
		targetId: string;
	}

	let { targetId }: Props = $props();

	let isOwner = $derived(serversStore.isOwner(userStore.user?.id || '', targetId));

	async function deleteServer(serverId: string) {
		const res = await backend.deleteServer(serverId);
		if (res.isErr()) {
			console.error(res.error);
		}

		if (res.isOk()) {
			serversStore.removeServer(serverId);
		}
	}

	async function leaveServer(serverId: string) {
		const res = await backend.leaveServer(serverId);
		if (res.isErr()) {
			console.error(res.error);
		}

		if (res.isOk()) {
			serversStore.removeServer(serverId);
		}
	}

	async function createServerInvite(serverId: string) {
		const res = await backend.createInvite(serverId);
		if (res.isErr()) {
			console.error(res.error);
		}

		if (res.isOk()) {
			navigator.clipboard.writeText(res.value);
		}
	}
</script>

<ContextMenu.Item class="context-menu-item" onclick={() => createServerInvite(targetId)}>
	<p class="flex items-center">Invite people</p>
	<UserInvite height={20} width={20} />
</ContextMenu.Item>
{#if isOwner}
	<ContextMenu.Item class="context-menu-item-danger" onclick={() => deleteServer(targetId)}>
		<p class="flex items-center">Delete server</p>
		<Bin height={20} width={20} />
	</ContextMenu.Item>
{:else}
	<ContextMenu.Item class="context-menu-item-danger" onclick={() => leaveServer(targetId)}>
		<p class="flex items-center">Leave server</p>
		<LogoutIcon height={20} width={20} />
	</ContextMenu.Item>
{/if}
