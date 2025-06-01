<script lang="ts">
	import { onDestroy, onMount } from 'svelte';
	import { Editor } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { Placeholder } from '@tiptap/extensions';

	let element: Element;
	let editor: Editor;
	let contentSet = $state(false);

	let { placeholder, content = $bindable() } = $props();

	onMount(() => {
		editor = new Editor({
			element: element,
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
				Placeholder.configure({
					placeholder: placeholder
				})
			],
			onTransaction: () => {
				editor = editor;
			},
			onUpdate: ({ editor }) => {
				content = editor.getText();
			},
			editorProps: {
				attributes: {
					class: 'about-me-input'
				}
			}
		});
	});

	onDestroy(() => {
		if (editor) {
			editor.destroy();
		}
	});

	$effect(() => {
		if (content && !contentSet) {
			editor.commands.setContent(content);
			contentSet = true;
		}
	});
</script>

<div class="pointer-events-auto w-full" bind:this={element}></div>
