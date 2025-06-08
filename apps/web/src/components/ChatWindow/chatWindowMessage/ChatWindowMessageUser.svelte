<script lang="ts">
	import UserProfileWithTrigger from 'components/UserProfile/UserProfileWithTrigger.svelte';
	import { core } from 'stores/core.svelte';
	import type { User } from 'types/types';
	import { formatMessageTime } from 'utils/date';

	let { author, isUserMentioned, isEdited, id, time } = $props();
</script>

<div class="flex items-baseline gap-x-2.5">
	<UserProfileWithTrigger user={author as User} side="bottom" sideOffset={5} y={-10}>
		<p
			class="pointer-events-auto text-sm font-semibold decoration-1 hover:cursor-pointer hover:underline"
		>
			{author.display_name}
		</p>
	</UserProfileWithTrigger>
	<time class={['text-xs', isUserMentioned ? 'text-main-300' : 'text-main-600']}>
		{formatMessageTime(time)}
	</time>
	{#if core.editingMessage.id === id || isEdited}
		<p
			class={[
				'absolute  right-3 uppercase',
				core.editingMessage.id !== id && isUserMentioned && '!text-main-300',
				core.editingMessage.id !== id ? 'text-main-600 top-3 text-xs' : 'text-accent-50 top-2'
			]}
		>
			[{core.editingMessage.id !== id ? 'Edited' : 'Editing'}]
		</p>
	{/if}
</div>
