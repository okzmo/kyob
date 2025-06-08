<script lang="ts">
	import { core } from 'stores/core.svelte';
	import { generateHTML } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { EmojisSuggestion } from '../chatWindowInput/extensions/emojis/emojis';
	import { CustomMention } from '../chatWindowInput/extensions/mentions/mentions';
	import EditMessageInput from '../editMessageInput/editMessageInput.svelte';
	import ChatWindowMessagesAttachments from './ChatWindowMessagesAttachments.svelte';

	let { id, server, channel, content, attachments } = $props();
</script>

<div class="flex w-full flex-col gap-y-1">
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
		<ChatWindowMessagesAttachments {attachments} />
	{/if}
</div>
