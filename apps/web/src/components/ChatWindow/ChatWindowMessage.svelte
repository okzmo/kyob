<script lang="ts">
	import { generateHTML } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { formatMessageTime } from '../../utils/date';

	interface Props {
		id: number;
		avatar: string;
		username: string;
		displayName: string;
		time: string;
		content: any;
	}
	let { id, avatar, displayName, content, username, time }: Props = $props();
</script>

<div
	id="message-{id}"
	class="hocus:bg-main-800/50 flex gap-x-3 px-4 py-2 transition-colors duration-100"
>
	<img src={avatar} alt="{username}'s avatar" class="h-[3rem] w-[3rem] rounded-full select-none" />
	<div class="pt-1">
		<div class="flex items-baseline gap-x-3 select-none">
			<p class="text-sm font-semibold">{displayName}</p>
			<time class="text-main-600 text-xs">{formatMessageTime(time)}</time>
		</div>
		<div class="mt-1 flex flex-col gap-y-1">
			{@html generateHTML(content, [
				StarterKit.configure({
					gapcursor: false,
					dropcursor: false,
					heading: false,
					orderedList: false,
					bulletList: false,
					blockquote: false
				})
			])}
		</div>
	</div>
</div>
