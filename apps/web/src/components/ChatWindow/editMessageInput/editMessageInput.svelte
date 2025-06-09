<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import MentionsList from 'components/ChatWindow/chatWindowInput/extensions/mentions/MentionsList.svelte';
	import { CustomMention } from 'components/ChatWindow/chatWindowInput/extensions/mentions/mentions';
	import { core } from 'stores/core.svelte';
	import { backend } from 'stores/backend.svelte';
	import { editorStore } from 'stores/editor.svelte';
	import EmojisList from '../chatWindowInput/extensions/emojis/EmojisList.svelte';
	import { EmojisSuggestion } from '../chatWindowInput/extensions/emojis/emojis';

	let element: Element;
	let editor: Editor;

	let { server, channel, messageId, content } = $props();

	async function editMessage(message: any) {
		if (editor.getText().length <= 0 || editor.getText().length > 2500) return;
		const ids =
			editor
				.getText()
				.match(/<@(\d+)>/g)
				?.map((match) => match.slice(2, -1)) || [];

		const payload = {
			content: message,
			mentions_users: [...new Set(ids)]
		};

		const res = await backend.editMessage(server.id, channel.id, messageId, payload);
		if (res.isErr()) {
			console.log(`${res.error.code}: ${res.error.error}`);
			return;
		}

		if (res.isOk()) {
			core.stopEditingMessage();
		}
	}

	onMount(() => {
		editorStore.currentInput = 'edit';

		editor = new Editor({
			element: element,
			autofocus: 'end',
			content: content,
			extensions: [
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
			onBlur: () => {
				core.stopEditingMessage();
			},
			editorProps: {
				attributes: {
					class: 'editor-message'
				},
				handleKeyDown: (_, ev) => {
					if (
						ev.key === 'Enter' &&
						!ev.shiftKey &&
						(!editorStore.mentionProps || editorStore.mentionProps.items.length === 0) &&
						(!editorStore.emojiProps || editorStore.emojiProps.items.length === 0)
					) {
						ev.preventDefault();
						editMessage(editor.getJSON());
						return true;
					}

					if (ev.key === 'Escape') {
						core.stopEditingMessage();
					}

					return false;
				}
			}
		});
	});

	onDestroy(() => {
		editorStore.currentInput = 'main';
		if (editor) {
			editor.destroy();
		}
	});
</script>

<div class="absolute top-2 left-0 w-[calc(100%-1.5rem)]">
	{#if editorStore.currentInput === 'edit' && editorStore.mentionProps}
		<MentionsList
			props={editorStore.mentionProps}
			bind:this={editorStore.mentionsListEl}
			class="absolute bottom-0 left-3.5 w-full"
		/>
	{/if}
	{#if editorStore.currentInput === 'edit' && editorStore.emojiProps}
		<EmojisList
			bind:this={editorStore.emojisListEl}
			props={editorStore.emojiProps}
			class="absolute bottom-0 left-3.5 w-full"
		/>
	{/if}
</div>
<div class="pointer-events-auto w-full" bind:this={element}></div>
