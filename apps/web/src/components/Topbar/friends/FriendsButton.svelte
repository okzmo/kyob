<script lang="ts">
	import { Popover } from 'bits-ui';
	import People from '../../ui/icons/People.svelte';
	import CustomPopoverContent from '../../ui/CustomPopoverContent/CustomPopoverContent.svelte';
	import FriendsList from './FriendsList.svelte';
	import AddFriend from './AddFriend.svelte';
	import { userStore } from '../../../stores/user.svelte';

	let isOpen = $state(false);
	let friends = $derived(userStore?.friends?.filter((f) => !f.sender) || []);
</script>

<Popover.Root open={isOpen} onOpenChange={(s) => (isOpen = s)}>
	<Popover.Trigger
		aria-label="Friends"
		class="text-main-400 hocus:text-accent-50 hocus:bg-accent-100/15 flex h-[2.25rem] w-[2.25rem] items-center justify-center rounded-lg transition-colors hover:cursor-pointer"
	>
		<People height={22} width={22} />
	</Popover.Trigger>
	<CustomPopoverContent
		class="bg-main-900 border-main-800 relative z-30 w-[20rem] rounded-2xl border p-2"
		align="end"
		side="bottom"
		sideOffset={10}
		y={-10}
	>
		<FriendsList {friends} />
		<AddFriend bind:isOpen />
	</CustomPopoverContent>
</Popover.Root>
