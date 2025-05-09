<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import EmojiIcon from '../ui/icons/EmojiIcon.svelte';
	import Plus from '../ui/icons/Plus.svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { Placeholder } from '@tiptap/extensions';
	import { backend } from '../../stores/backend.svelte';
	import type { Channel, Server } from '../../types/types';
	import { userStore } from '../../stores/user.svelte';

	interface Props {
		channel: Channel;
		server: Server;
	}

	let { channel, server }: Props = $props();

	let element: Element;
	let editor: Editor;

	async function prepareMessage(message: any) {
		if (editor.getText().length <= 0 || editor.getText().length > 2500) return;
		const payload = {
			author_id: userStore.user!.id,
			content: message
		};

		const res = await backend.sendMessage(server.id, channel.id, payload);
		if (res.isErr()) {
			console.log(`${res.error.code}: ${res.error.error}`);
		}

		editor.commands.clearContent();
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
				Placeholder.configure({ placeholder: `Message #${channel?.name} in ${server?.name}` })
			],
			onTransaction: () => {
				editor = editor;
			},
			editorProps: {
				handleKeyDown: (_, ev) => {
					if (ev.key === 'Enter' && !ev.shiftKey) {
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

<div class="bg-main-900 border-t-main-800 absolute bottom-0 left-0 flex w-full border-t px-4">
	<button
		class="text-main-600 hocus:text-main-200 absolute top-4.5 left-4 transition-colors duration-100 hover:cursor-pointer"
	>
		<Plus height={20} width={20} />
	</button>
	<!-- <input -->
	<!-- 	type="text" -->
	<!-- 	class="bg-main-900 placeholder:text-main-600 h-full w-full border-none placeholder:truncate focus:ring-0" -->
	<!-- 	placeholder="Message #{channel?.name} in {server?.name}" -->
	<!-- /> -->
	<div class="w-full px-7 py-4" bind:this={element}></div>
	<button
		class="text-main-600 hocus:text-main-200 absolute top-4.5 right-4 transition-colors duration-100 hover:cursor-pointer"
	>
		<EmojiIcon height={20} width={20} />
	</button>
</div>
