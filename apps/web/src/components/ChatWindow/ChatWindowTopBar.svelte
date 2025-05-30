<script lang="ts">
	import { userStore } from '../../stores/user.svelte';
	import { windows } from '../../stores/windows.svelte';
	import Corners from '../ui/Corners/Corners.svelte';
	import Close from '../ui/icons/Close.svelte';

	let { id, server, channel, friend } = $props();

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
		'inner-main-800 relative flex h-[2.375rem] w-full items-center justify-between px-2.5 transition duration-100 hover:cursor-grab active:cursor-grabbing',
		userStore.mention ? 'bg-accent-200' : 'bg-main-900'
	]}
>
	<Corners color="border-main-700" />
	<div
		class={[
			'flex items-center gap-x-2 text-sm transition-colors duration-100 select-none',
			userStore.mention ? 'text-accent-50 ' : 'text-main-400 '
		]}
	>
		{#if friend}
			<div class="flex items-center gap-x-1.5">
				<img
					src={friend.avatar}
					alt="{friend.display_name} avatar"
					class="h-[1.25rem] w-[1.25rem] rounded-full"
				/>
				<p class="select-none">{friend.display_name}</p>
			</div>
		{:else if server}
			<div class="flex items-center gap-x-1.5">
				<img
					src={server.avatar}
					alt="{server.name} server background"
					class="h-[1.25rem] w-[1.25rem] rounded-full"
				/>
				<p class="select-none">{server.name}</p>
			</div>
		{/if}
		{#if channel && !friend}
			|
			<p class="select-none">#{channel.name}</p>
		{/if}
	</div>
	<button
		class={[
			'flex items-center justify-center p-0.5 transition duration-100 hover:cursor-pointer',
			userStore.mention ? 'hocus:bg-red-600' : 'hocus:bg-main-800 hocus:inner-main-700'
		]}
		onclick={() => windows.closeWindow(id)}
	>
		<Close height={16} width={16} class="" />
	</button>
</div>
