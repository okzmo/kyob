<script lang="ts">
	import { ContextMenu } from 'bits-ui';
	import Bin from '../icons/Bin.svelte';
	import { userStore } from 'stores/user.svelte';
	import Pen from '../icons/Pen.svelte';
	import CopyIcon from '../icons/CopyIcon.svelte';
	import { backend } from 'stores/backend.svelte';
	import { windows } from 'stores/windows.svelte';
	import { core } from 'stores/core.svelte';
	import { serversStore } from 'stores/servers.svelte';
	import { generateText } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { CustomMention } from 'components/ChatWindow/chatWindowInput/extensions/mentions/mentions';

	let { authorId, targetId } = $props();

	let haveSelection = $derived(window.getSelection()?.type === 'Range');

	function handleDelete(messageId: string) {
		const window = windows.getActiveWindow();
		if (!window?.serverId || !window.channelId) return;
		const channel = serversStore.getChannel(window.serverId, window.channelId);
		if (channel.messages?.[0].id === messageId && channel.messages.length > 1) {
			channel.last_message_read = channel.messages?.[1].id;
			channel.last_message_sent = channel.messages?.[1].id;
		} else {
			channel.last_message_read = '';
			channel.last_message_sent = '';
		}

		backend.deleteMessage(window.serverId, window.channelId, messageId);
	}

	function handleEdit(messageId: string) {
		core.startEditingMessage(messageId);
	}

	function handleCopyMessage(messageId: string) {
		const activeWindow = windows.getActiveWindow();
		if (!activeWindow?.serverId || !activeWindow.channelId) return;
		const message = serversStore.getMessage(
			activeWindow.serverId,
			activeWindow.channelId,
			messageId
		);

		let text = generateText(message?.content, [
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
				renderText({ node }) {
					return `<@${node.attrs['user-id']}>`;
				}
			})
		]);

		if (haveSelection) {
			const selection = window.getSelection();
			const start = selection?.anchorOffset;
			const finish = selection?.focusOffset;

			if (finish !== 1) text = text.slice(start, finish);
		}

		navigator.clipboard.writeText(text);
	}
</script>

<ContextMenu.Item class="context-menu-item" onclick={() => handleCopyMessage(targetId)}>
	<div class="flex w-full items-center justify-between">
		{haveSelection ? 'Copy selection' : 'Copy text'}
		<CopyIcon height={20} width={20} />
	</div>
</ContextMenu.Item>
{#if authorId === userStore.user?.id}
	{#if !haveSelection}
		<ContextMenu.Item class="context-menu-item" onclick={() => handleEdit(targetId)}>
			<div class="flex w-full items-center justify-between">
				Edit Message
				<Pen height={20} width={20} />
			</div>
		</ContextMenu.Item>
		<ContextMenu.Item
			class="context-menu-item-danger text-red-400"
			onclick={() => handleDelete(targetId)}
		>
			<p class="flex items-center">Delete Message</p>
			<Bin height={20} width={20} />
		</ContextMenu.Item>
	{/if}
{/if}
