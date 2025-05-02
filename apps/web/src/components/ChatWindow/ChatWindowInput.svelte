<script lang="ts">
	import EmojiIcon from '../ui/icons/EmojiIcon.svelte';
	import Plus from '../ui/icons/Plus.svelte';

	let { channel, server } = $props();
</script>

<div
	class="bg-main-900 border-t-main-800 absolute bottom-0 left-0 flex h-[3.75rem] w-full items-center border-t px-4"
>
	<button
		class="text-main-600 hocus:text-main-200 transition-colors duration-100 hover:cursor-pointer"
	>
		<Plus height={20} width={20} />
	</button>
	<input
		type="text"
		class="bg-main-900 placeholder:text-main-600 h-full w-full border-none placeholder:truncate focus:ring-0"
		placeholder="Message #{channel?.name} in {server?.name}"
	/>
	<button
		class="text-main-600 hocus:text-main-200 transition-colors duration-100 hover:cursor-pointer"
	>
		<EmojiIcon height={20} width={20} />
	</button>
</div>
