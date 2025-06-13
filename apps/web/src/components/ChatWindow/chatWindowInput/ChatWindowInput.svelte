<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import EmojiIcon from 'components/ui/icons/EmojiIcon.svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { Placeholder } from '@tiptap/extensions';
	import { backend } from 'stores/backend.svelte';
	import type { Channel, Friend, Server } from 'types/types';
	import MentionsList from './extensions/mentions/MentionsList.svelte';
	import { CustomMention } from './extensions/mentions/mentions';
	import EmojisList from './extensions/emojis/EmojisList.svelte';
	import { editorStore } from 'stores/editor.svelte';
	import { EmojisSuggestion } from './extensions/emojis/emojis';
	import AttachmentButton from './AttachmentButton.svelte';
	import Attachments from './Attachments.svelte';

	interface Props {
		friend?: Friend;
		channel: Channel;
		server: Server;
	}

	let { channel, server, friend }: Props = $props();

	let element: Element;
	let editor: Editor;
	let attachments = $state<File[]>([]);

	async function prepareMessage(message: any) {
		if (editor.getText().length <= 0 || editor.getText().length > 2500) return;
		const ids =
			editor
				.getText()
				.match(/<@(\d+)>/g)
				?.map((match) => match.slice(2, -1)) || [];

		const payload = {
			content: message,
			mentions_users: [...new Set(ids)],
			attachments
		};

		const res = await backend.sendMessage(server.id, channel.id, payload);
		if (res.isErr()) {
			console.error(`${res.error.code}: ${res.error.error}`);
		}

		editor.commands.clearContent();
		attachments = [];
	}

	onMount(() => {
		editor = new Editor({
			element: element,
			extensions: [
				StarterKit.configure({
					gapcursor: false,
					dropcursor: false,
					heading: false,
					orderedList: false,
					bulletList: false,
					blockquote: false
				}),
				Placeholder.configure({
					placeholder: friend
						? `Message ${friend.display_name}`
						: `Message #${channel?.name} in ${server?.name}`
				}),
				EmojisSuggestion.configure({
					HTMLAttributes: {
						class: 'editor-emoji'
					},
					renderHTML({ options, node }) {
						return ['span', options.HTMLAttributes, `${node.attrs.emoji}`];
					}
				}),
				CustomMention.configure({
					HTMLAttributes: {
						class: 'editor-mention'
					},
					renderHTML({ options, node }) {
						return [
							'span',
							options.HTMLAttributes,
							`${node.attrs.mentionSuggestionChar}${node.attrs.label}`
						];
					}
				})
			],
			onTransaction: () => {
				editor = editor;
			},
			editorProps: {
				attributes: {
					class: 'chat-input'
				},
				handleKeyDown: (_, ev) => {
					if (
						ev.key === 'Enter' &&
						!ev.shiftKey &&
						(!editorStore.mentionProps || editorStore.mentionProps.items.length === 0) &&
						(!editorStore.emojiProps || editorStore.emojiProps.items.length === 0)
					) {
						ev.preventDefault();
						prepareMessage(editor.getJSON());
						return true;
					}

					return false;
				}
			}
		});
	});

	onDestroy(() => {
		if (editor) {
			editor.destroy();
		}
	});
</script>

<div class="flex w-full flex-col gap-y-1 px-2 pb-2">
	{#if editorStore.currentInput === 'main' && editorStore.mentionProps}
		<MentionsList
			props={editorStore.mentionProps}
			bind:this={editorStore.mentionsListEl}
			class="w-full"
		/>
	{/if}
	{#if editorStore.currentInput === 'main' && editorStore.emojiProps}
		<EmojisList
			props={editorStore.emojiProps}
			bind:this={editorStore.emojisListEl}
			class="w-full"
		/>
	{/if}
	<div class="bg-main-900 inner-shadow-input relative flex w-full flex-col transition duration-100">
		<div class="flex w-full">
			<AttachmentButton bind:attachments />
			<div class="max-h-[10rem] w-full" bind:this={element}></div>
			<button
				class="text-main-600 hocus:text-main-200 absolute top-4.5 right-4 z-[1] transition-colors duration-100 hover:cursor-pointer"
			>
				<EmojiIcon height={20} width={20} />
			</button>
		</div>

		{#if attachments.length > 0}
			<Attachments bind:attachments />
		{/if}
	</div>
</div>
