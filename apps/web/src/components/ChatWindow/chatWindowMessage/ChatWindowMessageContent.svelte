<script lang="ts">
	import { core } from 'stores/core.svelte';
	import { generateHTML } from '@tiptap/core';
	import EditMessageInput from '../editMessageInput/editMessageInput.svelte';
	import ChatWindowMessagesAttachments from './ChatWindowMessagesAttachments.svelte';
	import { getMessageExtensions } from '../chatWindowInput/editorConfig';

	let { id, server, channel, content, attachments } = $props();
</script>

<div class="flex w-full flex-col gap-y-1">
	{#if core.editingMessage.id === id}
		<EditMessageInput {server} {channel} {content} messageId={id} />
	{:else}
		<div class="[&>p]:pointer-events-auto">
			{@html generateHTML(content, getMessageExtensions())}
		</div>
	{/if}
	{#if attachments.length > 0}
		<ChatWindowMessagesAttachments {attachments} />
	{/if}
</div>
