<script lang="ts">
	import { generateHTML } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { formatMessageTime } from 'utils/date';
	import { CustomMention } from 'components/ChatWindow/chatWindowInput/extensions/mentions/mentions';
	import { onMount } from 'svelte';
	import { core } from 'stores/core.svelte';
	import EditMessageInput from '../editMessageInput/editMessageInput.svelte';
	import type { Channel, Server, User } from 'types/types';
	import Corners from '../../ui/Corners/Corners.svelte';
	import UserProfileWithTrigger from '../../UserProfile/UserProfileWithTrigger.svelte';
	import { EmojisSuggestion } from '../chatWindowInput/extensions/emojis/emojis';

	interface Props {
		id: string;
		author: Partial<User>;
		time: string;
		content: any;
		isUserMentioned: boolean;
		isEdited: boolean;
		server: Server;
		channel: Channel;
		attachments: string[];
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
		'relative flex items-start gap-x-3 px-4 py-2 transition-colors duration-100',
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
	<UserProfileWithTrigger user={author as User} side="right" align="center" x={-10}>
		<img
			src={author.avatar}
			alt="{author.username}'s avatar"
			class="mt-[2.5px] h-[3rem] w-[3rem] object-cover select-none hover:cursor-pointer active:translate-y-[1px]"
			draggable="false"
		/>
	</UserProfileWithTrigger>
	<div class="pointer-events-none pt-1">
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
		<div class="flex flex-col gap-y-1">
			{#if core.editingMessage.id === id}
				<EditMessageInput {server} {channel} {content} messageId={id} />
			{:else}
				<div class="[&>p]:pointer-events-auto">
					{@html generateHTML(content, [
						StarterKit.configure({
							gapcursor: false,
							dropcursor: false,
							heading: false,
							orderedList: false,
							bulletList: false,
							blockquote: false
						}),
						EmojisSuggestion.configure({
							HTMLAttributes: {
								class: 'emoji'
							},
							renderHTML({ options, node }) {
								return ['span', options.HTMLAttributes, `${node.attrs.emoji}`];
							}
						}),
						CustomMention.configure({
							HTMLAttributes: {
								class: 'mention'
							},
							renderHTML({ options, node }) {
								return ['button', options.HTMLAttributes, `${node.attrs.label}`];
							}
						})
					])}
				</div>
			{/if}
			{#if attachments.length > 0}
				<div
					class={attachments.length > 1
						? 'grid w-fit max-w-[90%] grid-cols-2 gap-2 @xl:max-w-[65%] @3xl:max-w-[45%] @5xl:max-w-[35%]'
						: 'max-w-[75%] @xl:max-w-[50%] @3xl:max-w-[35%] @5xl:max-w-[25%]'}
				>
					{#each attachments as attachment, idx (idx)}
						<figure
							class={['attachment relative select-none', attachments.length > 1 && 'aspect-square']}
						>
							<img src={attachment} alt="Attachment" class="h-full w-full object-cover" />
						</figure>
					{/each}
				</div>
			{/if}
		</div>
	</div>
</div>

<style>
	.attachment::before {
		content: '';
		position: absolute;
		inset: 0;
		box-shadow: inset 0 0 0 1px #fafafa33;
	}
</style>
