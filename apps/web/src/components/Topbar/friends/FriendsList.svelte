<script lang="ts">
	import Friend from './Friend.svelte';
	import type { Friend as FriendType } from '../../../types/types';

	interface Props {
		friends: FriendType[];
	}

	let { friends }: Props = $props();
</script>

<ul class="flex flex-col gap-y-2">
	{#if friends.length <= 0}
		<p class="text-main-500 flex w-full items-center justify-center py-2">No friends yet ;(</p>
	{:else}
		{#each friends as friend (friend.id)}
			<li class="relative">
				<Friend
					id={friend.id}
					friendshipId={friend.friendship_id}
					displayName={friend.display_name}
					avatar={friend.avatar}
					accepted={friend.accepted}
				/>
			</li>
		{/each}
	{/if}
</ul>
