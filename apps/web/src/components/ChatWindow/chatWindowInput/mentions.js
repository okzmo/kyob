import { PasteRule } from '@tiptap/core';
import Mention from '@tiptap/extension-mention';
import { windows } from 'stores/windows.svelte';
import { serversStore } from 'stores/servers.svelte';

export const CustomMention = Mention.extend({
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
	}
});
