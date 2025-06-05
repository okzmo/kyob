<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { PluginKey } from '@tiptap/pm/state';
	import type { SuggestionProps } from '@tiptap/suggestion';
	import MentionsList from '../chatWindowInput/MentionsList.svelte';
	import { CustomMention } from '../chatWindowInput/mentions';
	import { serversStore } from 'stores/servers.svelte';
	import { windows } from 'stores/windows.svelte';
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
					},
					renderText({ node }) {
						return `<@${node.attrs['user-id']}>`;
					},
					suggestions: [
						{
							char: '@',
							pluginKey: new PluginKey('at'),
							items: ({ query }) => {
								const res = [];

								const activeWindow = windows.getActiveWindow();
								if (!activeWindow?.serverId) return [];
								const users = serversStore.getServer(activeWindow?.serverId).members;

								for (const user of users) {
									if (
										user?.username?.toLowerCase().includes(query.toLowerCase()) ||
										user?.display_name?.toLowerCase().includes(query.toLowerCase())
									) {
										res.push(user);
									}
								}

								return res;
							},
							render: () => {
								return {
									onStart: (props) => {
										mentionProps = props;
									},
									onUpdate: (props) => {
										mentionProps = props;
									},
									onExit: () => {
										mentionProps = null;
									},
									onKeyDown: (props) => {
										if (props.event.key === 'Escape') {
											mentionProps = null;
											return true;
										}

										return mentionsListEl?.handleKeyDown(props);
									}
								};
							}
						}
					]
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
