<script lang="ts">
	import { ContextMenu } from 'bits-ui';
	import UserInvite from '../icons/UserInvite.svelte';
	import Bin from '../icons/Bin.svelte';
	import { serversStore } from 'stores/servers.svelte';
	import { userStore } from 'stores/user.svelte';
	import { page } from '$app/state';
	import { backend } from 'stores/backend.svelte';

	interface Props {
		targetId: string;
	}

	let { targetId }: Props = $props();

	let isOwner = $derived(serversStore.isOwner(userStore.user?.id || '', page.params.server_id));

	async function deleteChannel(serverId: string, channelId: string) {
		const res = await backend.deleteChannel(serverId, channelId);
		if (res.isErr()) {
			console.error(res.error);
		}
	}
</script>

<ContextMenu.Item class="context-menu-item">
	<div class="flex w-full items-center justify-between">Mark as read</div>
</ContextMenu.Item>
{#if isOwner}
	<ContextMenu.Item class="context-menu-item">
		<div class="flex w-full items-center justify-between">
			Edit Channel
			<UserInvite height={20} width={20} />
		</div>
	</ContextMenu.Item>
	<ContextMenu.Item
		class="context-menu-item-danger"
		onclick={() => deleteChannel(page.params.server_id, targetId)}
	>
		<p class="flex items-center">Delete Channel</p>
		<Bin height={20} width={20} />
	</ContextMenu.Item>
{/if}
