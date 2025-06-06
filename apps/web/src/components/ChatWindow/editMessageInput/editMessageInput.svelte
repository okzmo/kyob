<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import type { SuggestionProps } from '@tiptap/suggestion';
	import MentionsList from 'components/ChatWindow/chatWindowInput/extensions/mentions/MentionsList.svelte';
	import { CustomMention } from 'components/ChatWindow/chatWindowInput/extensions/mentions/mentions';
	import { core } from 'stores/core.svelte';
	import { backend } from 'stores/backend.svelte';

	let element: Element;
	let editor: Editor;
	let mentionProps = $state<SuggestionProps | null>();
	let mentionsListEl = $state<any>();

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
					class: 'editor-message'
				},
				handleKeyDown: (_, ev) => {
					if (
						ev.key === 'Enter' &&
						!ev.shiftKey &&
						(!mentionProps || mentionProps.items.length === 0)
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
		if (editor) {
			editor.destroy();
		}
	});
</script>

{#if mentionProps}
	<MentionsList
		props={mentionProps}
		bind:this={mentionsListEl}
		class="absolute -top-[2.25rem] left-2 w-[calc(100%-1rem)]"
	/>
{/if}
<div class="pointer-events-auto w-full" bind:this={element}></div>
