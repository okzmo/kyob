<script lang="ts">
	import { Popover } from 'bits-ui';
	import { backend } from 'stores/backend.svelte';
	import { userStore } from 'stores/user.svelte';
	import Check from '../../ui/icons/Check.svelte';
	import Close from '../../ui/icons/Close.svelte';
	import MoreIcon from '../../ui/icons/MoreIcon.svelte';
	import CustomPopoverContent from '../../ui/CustomPopoverContent/CustomPopoverContent.svelte';
	import Bin from '../../ui/icons/Bin.svelte';
	import { windows } from 'stores/windows.svelte';
	import Corners from 'components/ui/Corners/Corners.svelte';

	let {
		id,
		friendshipId,
		channelId,
		displayName,
		avatar,
		accepted = false,
		isOpen = $bindable()
	} = $props();

	async function acceptFriend() {
		const res = await backend.acceptFriend({
			friendship_id: friendshipId,
			user_id: userStore.user!.id,
			friend_id: id
		});
		if (res.isErr()) {
			console.error(res.error);
			return;
		}

		if (res.isOk()) {
			userStore.acceptFriend({ friendshipId });
		}
	}

	async function deleteFriend() {
		const res = await backend.deleteFriend({
			friendship_id: friendshipId,
			friend_id: id,
			user_id: userStore.user!.id
		});
		if (res.isErr()) {
			console.error(res.error);
			return;
		}

		if (res.isOk()) {
			userStore.deleteFriend(friendshipId);
		}
	}
</script>

<button
	class="group hocus:bg-main-800 hocus:inner-main-700 relative flex w-full items-center gap-x-2.5 py-1.5 pr-4 pl-1.5 text-left transition hover:cursor-pointer"
	onclick={() => {
		windows.createWindow({
			id: `window-${id}`,
			serverId: 'global',
			channelId: channelId,
			friendId: id
		});
		isOpen = false;
	}}
>
	<img src={avatar} alt="avatar" class="h-[2.75rem] w-[2.75rem] object-cover" />
	<div class="flex flex-col">
		<p class="leading-[1.15rem] font-medium transition-colors">
			{displayName}
		</p>
		<p class="text-main-400 group-hocus:text-main-300 text-sm leading-[1.15rem] transition-colors">
			Connected
		</p>
	</div>
</button>

{#if accepted}
	<Popover.Root>
		<Popover.Trigger
			class="bg-main-900 hocus:bg-main-800 hocus:text-main-50 text-main-600 absolute top-1/2 right-2 -translate-y-1/2 rounded-full p-1.5 transition-colors duration-100 hover:cursor-pointer"
		>
			<MoreIcon height={20} width={20} />
		</Popover.Trigger>
		<CustomPopoverContent
			class="bg-main-900 border-main-800 inner-shadow-main-800 relative z-30 w-[14rem] p-2 select-none"
			align="end"
			side="bottom"
			sideOffset={10}
			y={-10}
		>
			<Corners color="border-main-700" />
			<button class="context-menu-item-danger w-full" onclick={deleteFriend}>
				<p class="flex items-center">Remove Friend</p>
				<Bin height={20} width={20} />
			</button>
		</CustomPopoverContent>
	</Popover.Root>
{:else}
	<div class="absolute top-1/2 right-2 flex items-center gap-x-2">
		<button
			onclick={acceptFriend}
			class="hocus:bg-green-400/50 flex h-[2rem] w-[2rem] -translate-y-1/2 items-center justify-center border border-green-400 bg-green-400/20 text-green-400 transition-colors duration-100 hover:cursor-pointer"
		>
			<Check height={16} width={16} />
		</button>
		<button
			onclick={deleteFriend}
			class="hocus:bg-red-400/50 flex h-[2rem] w-[2rem] -translate-y-1/2 items-center justify-center border border-red-400 bg-red-400/20 text-red-400 transition-colors duration-100 hover:cursor-pointer"
		>
			<Close height={16} width={16} />
		</button>
	</div>
{/if}
