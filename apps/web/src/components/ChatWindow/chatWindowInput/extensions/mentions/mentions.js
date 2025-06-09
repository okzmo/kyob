import { PasteRule } from '@tiptap/core';
import Mention from '@tiptap/extension-mention';
import { windows } from 'stores/windows.svelte';
import { serversStore } from 'stores/servers.svelte';
import { PluginKey } from '@tiptap/pm/state';
import { editorStore } from 'stores/editor.svelte';

const MentionExtended = Mention.extend({
	name: 'customMention',

	addPasteRules() {
		return [
			new PasteRule({
				find: /<@(\d+)>/g,
				handler: ({ state, range, match }) => {
					const userId = match[1];
					const activeWindow = windows.getActiveWindow();

					if (activeWindow?.serverId) {
						const user = serversStore.getMemberById(activeWindow.serverId, userId);

						const attributes = {
							'user-id': userId,
							label: user?.display_name || 'Unknown User',
							mentionSuggestionChar: '@'
						};

						const { tr } = state;
						tr.replaceWith(range.from, range.to, this.type.create(attributes));
					}
				}
			})
		];
	},

	addAttributes() {
		return {
			'user-id': {
				default: null
			},
			label: {
				default: null
			},
			mentionSuggestionChar: {
				default: '@'
			}
		};
	},

	addStorage() {
		return {
			mentionProps: null,
			mentionsListEl: null
		};
	}
});

export const CustomMention = MentionExtended.configure({
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

				let users = [];

				if (activeWindow.serverId === 'global' && activeWindow.channelId) {
					users =
						serversStore.getChannel(activeWindow.serverId, activeWindow.channelId).users || [];
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
						editorStore.mentionProps = props;
					},
					onUpdate: (props) => {
						editorStore.mentionProps = props;
					},
					onExit: () => {
						editorStore.mentionProps = null;
					},
					onKeyDown: (props) => {
						if (props.event.key === 'Escape') {
							editorStore.mentionProps = null;
							return true;
						}

						return editorStore.mentionsListEl?.handleKeyDown(props);
					}
				};
			}
		}
	]
});
