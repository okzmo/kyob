<script lang="ts">
	import { userStore } from '../../stores/user.svelte';
	import { windows } from '../../stores/windows.svelte';
	import Close from '../ui/icons/Close.svelte';

	let { id, server, channel } = $props();

	$effect(() => {
		if (userStore.mention) {
			setTimeout(() => {
				userStore.mention = false;
			}, 500);
		}
	});
</script>

<div
	id={`window-top-bar-${id}`}
	class={[
		'flex h-[2.375rem] w-full items-center justify-between rounded-t-[14px] border px-2.5 transition-colors duration-100',
		userStore.mention ? 'border-red-500 bg-red-800' : 'bg-main-800 border-main-600'
	]}
>
	<div
		class={[
			'flex items-center gap-x-2 text-sm transition-colors duration-100 select-none',
			userStore.mention ? 'text-red-300 ' : 'text-main-400 '
		]}
	>
		{#if server}
			<div class="flex items-center gap-x-1.5">
				<img
					src={server.avatar}
					alt="{server.name} server background"
					class="h-[1.25rem] w-[1.25rem] rounded-full"
				/>
				<p class="select-none">{server.name}</p>
			</div>
		{/if}
		{#if channel}
			|
			<p class="select-none">#{channel.name}</p>
		{/if}
	</div>
	<button
		class={[
			'flex items-center justify-center rounded-md p-0.5 transition-colors duration-100 hover:cursor-pointer',
			userStore.mention ? 'hocus:bg-red-600' : 'hocus:bg-main-600'
		]}
		onclick={() => windows.closeWindow(id)}
	>
		<Close height={16} width={16} class="" />
	</button>
</div>
