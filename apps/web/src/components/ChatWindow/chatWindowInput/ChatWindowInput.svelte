<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import EmojiIcon from '../../ui/icons/EmojiIcon.svelte';
	import Plus from '../../ui/icons/Plus.svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { Placeholder } from '@tiptap/extensions';
	import { backend } from '../../../stores/backend.svelte';
	import type { Channel, Friend, Server, User } from '../../../types/types';
	import { PluginKey } from '@tiptap/pm/state';
	import type { SuggestionProps } from '@tiptap/suggestion';
	import MentionsList from './MentionsList.svelte';
	import { CustomMention } from './mentions';
	import { serversStore } from '../../../stores/servers.svelte';
	import { windows } from '../../../stores/windows.svelte';

	interface Props {
		friend?: Friend;
		channel: Channel;
		server: Server;
	}

	let { channel, server, friend }: Props = $props();

	let element: Element;
	let editor: Editor;
	let mentionProps = $state<SuggestionProps | null>();
	let mentionsListEl = $state<any>();

	async function prepareMessage(message: any) {
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
				Placeholder.configure({
					placeholder: friend
						? `Message ${friend.display_name}`
						: `Message #${channel?.name} in ${server?.name}`
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

								let users: Partial<User>[] = [];

								if (activeWindow.serverId === 'global' && activeWindow.channelId) {
									users =
										serversStore.getChannel(activeWindow.serverId, activeWindow.channelId).users ||
										[];
								} else {
									users = serversStore.getServer(activeWindow?.serverId).members;
								}

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
				handleKeyDown: (_, ev) => {
					if (
						ev.key === 'Enter' &&
						!ev.shiftKey &&
						(!mentionProps || mentionProps.items.length === 0)
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

{#if mentionProps}
	<MentionsList
		props={mentionProps}
		bind:this={mentionsListEl}
		class="bottom-[4rem] left-[0.4rem] w-[calc(100%-0.8rem)]"
	/>
{/if}
<div
	class="bg-main-900 inner-shadow-input absolute bottom-2 left-2 flex w-[calc(100%-1rem)] px-4 transition duration-100"
>
	<button
		class="text-main-600 hocus:text-main-200 absolute top-4.5 left-4 transition-colors duration-100 hover:cursor-pointer"
	>
		<Plus height={20} width={20} />
	</button>
	<div class="max-h-[10rem] w-full px-7 py-4" bind:this={element}></div>
	<button
		class="text-main-600 hocus:text-main-200 absolute top-4.5 right-4 transition-colors duration-100 hover:cursor-pointer"
	>
		<EmojiIcon height={20} width={20} />
	</button>
</div>
