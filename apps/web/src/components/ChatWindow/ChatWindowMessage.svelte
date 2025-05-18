<script lang="ts">
	import { generateHTML } from '@tiptap/core';
	import StarterKit from '@tiptap/starter-kit';
	import { formatMessageTime } from '../../utils/date';
	import { CustomMention } from './chatWindowInput/mentions';

	interface Props {
		id: number;
		avatar: string;
		username: string;
		displayName: string;
		time: string;
		content: any;
		isUserMentioned: boolean;
	}
	let { id, avatar, displayName, content, username, time, isUserMentioned }: Props = $props();
</script>

<div
	id="message-{id}"
	class={[
		' flex gap-x-3 px-4 py-2 transition-colors duration-100',
		isUserMentioned ? 'message-mention' : 'hocus:bg-main-800/50'
	]}
>
	<img
		src={avatar}
		alt="{username}'s avatar"
		class="h-[3rem] w-[3rem] rounded-full object-cover select-none"
	/>
	<div class="pt-1">
		<div class="flex items-baseline gap-x-2.5 select-none">
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
				}),
				CustomMention.configure({
					HTMLAttributes: {
						class: 'mention'
					},
					renderHTML({ options, node }) {
						return [
							'button',
							options.HTMLAttributes,
							[
								'img',
								{
									src: node.attrs.avatar || '',
									alt: `${node.attrs.label || ''} avatar`
								}
							],
							`${node.attrs.label}`
						];
					}
				})
			])}
		</div>
	</div>
</div>
