<script lang="ts">
	import { generateHTML } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { formatMessageTime } from '../../utils/date';
	import { CustomMention } from './chatWindowInput/mentions';
	import { onMount } from 'svelte';
	import { core } from '../../stores/core.svelte';
	import EditMessageInput from './editMessageInput/editMessageInput.svelte';
	import type { Channel, Server } from '../../types/types';

	interface Props {
		id: number;
		avatar: string;
		username: string;
		displayName: string;
		time: string;
		content: any;
		isUserMentioned: boolean;
		isEdited: boolean;
		userId: number;
		server: Server;
		channel: Channel;
	}
	let {
		id,
		userId,
		avatar,
		displayName,
		content,
		username,
		time,
		isUserMentioned,
		isEdited,
		server,
		channel
	}: Props = $props();

	let message = $state<HTMLElement>();

	function handleMention(e: MouseEvent) {
		const target = e.target as HTMLButtonElement;
		const userId = Number(target.attributes['user-id'].value);
		if (core.profileOpen.userId !== userId) {
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
	data-author-id={userId}
	bind:this={message}
	class={[
		' relative flex gap-x-3 px-4 py-2 transition-colors duration-100',
		isUserMentioned
			? 'message-mention'
			: core.editingMessage.id === id
				? 'bg-accent-100/10'
				: 'hocus:bg-main-800/50'
	]}
>
	<img
		src={avatar}
		alt="{username}'s avatar"
		class="h-[3rem] w-[3rem] rounded-full object-cover select-none"
	/>
	<div class="pointer-events-none pt-1">
		<div class="flex items-baseline gap-x-2.5 select-none">
			<p class="text-sm font-semibold">{displayName}</p>
			<time class="text-main-600 text-xs">{formatMessageTime(time)}</time>
			{#if core.editingMessage.id === id || isEdited}
				<p
					class={[
						'absolute  right-3 uppercase',
						core.editingMessage.id !== id ? 'text-main-600 top-3 text-xs' : 'text-accent-50 top-2'
					]}
				>
					[{core.editingMessage.id !== id ? 'Edited' : 'Editing'}]
				</p>
			{/if}
		</div>
		<div class="mt-1 flex flex-col gap-y-1">
			{#if core.editingMessage.id === id}
				<EditMessageInput {server} {channel} {content} messageId={id} />
			{:else}
				{@html generateHTML(content, [
					StarterKit.configure({
						gapcursor: false,
						dropcursor: false,
						heading: false,
						orderedList: false,
						bulletList: false,
						blockquote: false
					}),
					CustomMention.configure({
						HTMLAttributes: {
							class: 'mention'
						},
						renderHTML({ options, node }) {
							return [
								'button',
								options.HTMLAttributes,
								[
									'img',
									{
										src: node.attrs.avatar || '',
										alt: `${node.attrs.label || ''} avatar`
									}
								],
								`${node.attrs.label}`
							];
						}
					})
				])}
			{/if}
		</div>
	</div>
</div>
