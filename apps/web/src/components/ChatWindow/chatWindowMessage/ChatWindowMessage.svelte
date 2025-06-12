<script lang="ts">
	import { onMount } from 'svelte';
	import { core } from 'stores/core.svelte';
	import type { Attachment, Channel, Server, User } from 'types/types';
	import Corners from '../../ui/Corners/Corners.svelte';
	import ChatWindowMessageUser from './ChatWindowMessageUser.svelte';
	import ChatWindowMessageContent from './ChatWindowMessageContent.svelte';
	import UserProfileWithTriggerAndFetch from 'components/UserProfile/UserProfileWithTriggerAndFetch.svelte';

	interface Props {
		id: string;
		author: Partial<User>;
		time: string;
		content: any;
		isUserMentioned: boolean;
		isEdited: boolean;
		server: Server;
		channel: Channel;
		attachments: Attachment[];
	}
	let {
		id,
		author,
		content,
		time,
		isUserMentioned,
		isEdited,
		server,
		channel,
		attachments
	}: Props = $props();

	let message = $state<HTMLElement>();

	function handleMention(e: MouseEvent) {
		const target = e.target as HTMLButtonElement;
		const userId = target.attributes.getNamedItem('user-id')?.value;
		if (userId && core.profileOpen.userId !== userId) {
			core.openProfile(userId, target);
		}
	}

	onMount(() => {
		const mentions = message?.querySelectorAll('[data-type="mention"]');

		if (!mentions) return;
		for (const mention of mentions) {
			(mention as HTMLButtonElement).addEventListener('click', handleMention);
		}

		return () => {
			for (const mention of mentions) {
				(mention as HTMLButtonElement).removeEventListener('click', handleMention);
			}
		};
	});
</script>

<div
	id="message-{id}"
	data-author-id={author.id}
	bind:this={message}
	class={[
		'@container relative flex items-start gap-x-3 px-4 py-2 transition-colors duration-100',
		isUserMentioned
			? 'message-mention'
			: core.editingMessage.id === id
				? 'bg-accent-100/10'
				: 'hocus:bg-main-800/50'
	]}
>
	{#if isUserMentioned}
		<Corners color="border-mention-100" />
	{/if}
	<UserProfileWithTriggerAndFetch userId={author.id!} side="right" align="center" x={-10}>
		<img
			src={author.avatar}
			alt="{author.username}'s avatar"
			class="mt-[2.5px] h-[3rem] w-[3rem] object-cover select-none hover:cursor-pointer active:translate-y-[1px]"
			draggable="false"
		/>
	</UserProfileWithTriggerAndFetch>
	<div class="pointer-events-none w-full pt-1">
		<ChatWindowMessageUser {id} {author} {time} {isUserMentioned} {isEdited} />
		<ChatWindowMessageContent {id} {server} {channel} {content} {attachments} />
	</div>
</div>
