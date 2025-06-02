<script lang="ts">
	import { Popover } from 'bits-ui';
	import People from '../../ui/icons/People.svelte';
	import CustomPopoverContent from '../../ui/CustomPopoverContent/CustomPopoverContent.svelte';
	import FriendsList from './FriendsList.svelte';
	import AddFriend from './AddFriend.svelte';
	import { userStore } from 'stores/user.svelte';
	import Corners from '../../ui/Corners/Corners.svelte';

	let isOpen = $state(false);
	let friends = $derived(userStore?.friends?.filter((f) => !f.sender || f.accepted) || []);
</script>

<Popover.Root open={isOpen} onOpenChange={(s) => (isOpen = s)}>
	<Popover.Trigger
		aria-label="Friends"
		class="top-bar-button group text-main-400 hocus:text-accent-50 hocus:bg-accent-100/15 relative flex h-[2.25rem] w-[2.25rem] items-center justify-center transition-colors hover:cursor-pointer"
	>
		<Corners color="border-main-300" class="group-hocus:border-accent-100 duration-100" />
		<People height={22} width={22} />
	</Popover.Trigger>
	<CustomPopoverContent
		class="bg-main-900 border-main-800 inner-shadow-main-800 relative z-30 w-[20rem] p-2"
		align="end"
		side="bottom"
		sideOffset={10}
		y={-10}
	>
		<Corners color="border-main-700" />
		<FriendsList {friends} bind:isOpen />
		<AddFriend bind:isOpen />
	</CustomPopoverContent>
</Popover.Root>
