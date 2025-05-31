<script lang="ts">
	import { page } from '$app/state';
	import { windows } from '../../../stores/windows.svelte';
	import Corners from '../Corners/Corners.svelte';
	import HashChat from '../icons/HashChat.svelte';

	interface Props {
		id: string;
		name: string;
		type: 'textual' | 'voice';
		x: number;
		y: number;
		unread?: boolean;
	}

	let { id, name, type, x, y, unread = false }: Props = $props();
</script>

<button
	id="channelButton-{id}"
	class={[
		'group hocus:bg-accent-800 inner-shadow-main-800 hocus:text-accent-50 hocus:inner-accent/25 absolute z-20 flex items-center gap-x-2.5 px-4 py-3 font-medium transition-[color,box-shadow] duration-100 hover:cursor-pointer',
		unread ? 'bg-main-800 text-main-50' : 'bg-main-900 text-main-400'
	]}
	style="transform: translate({x}px, {y}px);"
	onclick={() => {
		windows.createWindow({
			id: `window-${id}`,
			serverId: page.params.server_id,
			channelId: id
		});
	}}
>
	<Corners color="border-main-700" class="group-hocus:border-accent-100 duration-100" />
	{#if type == 'textual'}
		<HashChat height={20} width={20} />
	{:else if type == 'voice'}
		div
	{/if}
	{name}
</button>
