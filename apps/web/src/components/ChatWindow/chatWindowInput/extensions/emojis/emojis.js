import Mention from '@tiptap/extension-mention';
import { PluginKey } from '@tiptap/pm/state';
import { editorStore } from 'stores/editor.svelte';
import emojis from 'emojibase-data/en/compact.json';

const EmojiExtended = Mention.extend({
	name: 'emojis',
	addAttributes() {
		return {
			emoji: {
				default: null
			},
			label: {
				default: null
			},
			url: {
				default: null
			}
		};
	}
});

export const EmojisSuggestion = EmojiExtended.configure({
	renderText({ node }) {
		return node.attrs['emoji'];
	},

	suggestions: [
		{
			char: ':',
			pluginKey: new PluginKey('emoji-suggestion'),
			items: ({ query }) => {
				const lowerQuery = query.toLocaleLowerCase();

				let res = emojis
					.filter((emoji) => {
						return (
							emoji.unicode === lowerQuery ||
							emoji.label.replaceAll(' ', '').includes(lowerQuery) ||
							emoji.tags?.includes(lowerQuery)
						);
					})
					.slice(0, 8);

				return res;
			},
			render: () => {
				return {
					onStart: (props) => {
						editorStore.emojiProps = props;
					},
					onUpdate: (props) => {
						editorStore.emojiProps = props;
					},
					onExit: () => {
						editorStore.emojiProps = null;
					},
					onKeyDown: (props) => {
						if (props.event.key === 'Escape') {
							editorStore.emojiProps = null;
							return true;
						}

						return editorStore.emojisListEl?.handleKeyDown(props);
					}
				};
			}
		}
	]
});
