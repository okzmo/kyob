import type { SuggestionProps } from '@tiptap/suggestion';

class Editor {
	mentionProps = $state<SuggestionProps | null>();
	mentionsListEl = $state<any>(null);
	emojiProps = $state<SuggestionProps | null>();
	emojisListEl = $state<any>(null);
}

export const editorStore = new Editor();
