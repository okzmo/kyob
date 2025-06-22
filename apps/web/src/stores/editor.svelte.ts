import type { SuggestionProps } from '@tiptap/suggestion';

class Editor {
	currentChannel = $state('');
	currentInput = $state<'main' | 'edit'>('main');
	mentionProps = $state<SuggestionProps | null>();
	mentionsListEl = $state<any>(null);
	emojiProps = $state<SuggestionProps | null>();
	emojisListEl = $state<any>(null);
}

export const editorStore = new Editor();
