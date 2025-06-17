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
	import { createEditorConfig } from './editorConfig';

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
		const everyone = editor.getText().includes('@everyone');
		const ids =
			editor
				.getText()
				.match(/<@(\d+)>/g)
				?.map((match) => match.slice(2, -1)) || [];

		const payload = {
			content: message,
			mentions_users: [...new Set(ids)],
			everyone: everyone,
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
		editorStore.currentChannel = channel.id;

		editor = new Editor(
			createEditorConfig({
				element: element,
				placeholder: friend
					? `Message ${friend.display_name}`
					: `Message #${channel?.name} in ${server?.name}`,
				onTransaction: () => {
					editor = editor;
				},
				editorProps: {
					attributes: {
						class: 'chat-input'
					}
				},
				onEnterPress: () => prepareMessage(editor.getJSON())
			})
		);
	});

	onDestroy(() => {
		if (editor) {
			editor.destroy();
		}
	});
</script>

<div class="flex w-full flex-col gap-y-1 px-2 pb-2">
	{#if editorStore.currentInput === 'main' && editorStore.currentChannel === channel.id && editorStore.mentionProps}
		<MentionsList
			props={editorStore.mentionProps}
			bind:this={editorStore.mentionsListEl}
			class="w-full"
		/>
	{/if}
	{#if editorStore.currentInput === 'main' && editorStore.currentChannel === channel.id && editorStore.emojiProps}
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
