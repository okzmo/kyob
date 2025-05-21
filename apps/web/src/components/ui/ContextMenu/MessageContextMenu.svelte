<script lang="ts">
	import { ContextMenu } from 'bits-ui';
	import Bin from '../icons/Bin.svelte';
	import { userStore } from '../../../stores/user.svelte';
	import Pen from '../icons/Pen.svelte';
	import CopyIcon from '../icons/CopyIcon.svelte';
	import { backend } from '../../../stores/backend.svelte';
	import { windows } from '../../../stores/windows.svelte';
	import { core } from '../../../stores/core.svelte';
	import { serversStore } from '../../../stores/servers.svelte';
	import { generateText } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { CustomMention } from '../../ChatWindow/chatWindowInput/mentions';

	let { authorId, targetId } = $props();

	function handleDelete(messageId: string) {
		const w = windows.getActiveWindow();
		if (!w) return;
		backend.deleteMessage(w?.serverId, w?.channelId, messageId);
	}

	function handleEdit(messageId: string) {
		core.startEditingMessage(messageId);
	}

	function handleCopyMessage(messageId: string) {
		const w = windows.getActiveWindow();
		if (!w) return;
		const message = serversStore.getMessage(w?.serverId, w?.channelId, messageId);

		const text = generateText(message?.content, [
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

		console.log(text);
	}
</script>

<ContextMenu.Item
	class="rounded-button data-highlighted:bg-main-800 flex h-10 items-center rounded-lg py-3 pr-1.5 pl-3  font-medium select-none hover:cursor-pointer focus-visible:outline-none"
	onclick={() => handleCopyMessage(targetId)}
>
	<div class="flex w-full items-center justify-between">
		Copy text
		<CopyIcon height={20} width={20} />
	</div>
</ContextMenu.Item>
{#if authorId === userStore.user?.id}
	<ContextMenu.Item
		class="rounded-button data-highlighted:bg-main-800 flex h-10 items-center rounded-lg py-3 pr-1.5 pl-3  font-medium select-none hover:cursor-pointer focus-visible:outline-none"
		onclick={() => handleEdit(targetId)}
	>
		<div class="flex w-full items-center justify-between">
			Edit Message
			<Pen height={20} width={20} />
		</div>
	</ContextMenu.Item>
	<ContextMenu.Item
		class="rounded-button flex h-10 items-center justify-between rounded-lg py-3 pr-1.5 pl-3 font-medium  text-red-400 select-none hover:cursor-pointer focus-visible:outline-none  data-highlighted:bg-red-400/20"
		onclick={() => handleDelete(targetId)}
	>
		<p class="flex items-center">Delete Message</p>
		<Bin height={20} width={20} />
	</ContextMenu.Item>
{/if}
