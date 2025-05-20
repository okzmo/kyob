import Mention from '@tiptap/extension-mention';

export const CustomMention = Mention.extend({
	addAttributes() {
		return {
			'user-id': {
				default: null
			},
			avatar: {
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
