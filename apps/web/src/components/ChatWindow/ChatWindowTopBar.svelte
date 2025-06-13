<script lang="ts">
	import { userStore } from 'stores/user.svelte';
	import { windows } from 'stores/windows.svelte';
	import Corners from '../ui/Corners/Corners.svelte';
	import Close from '../ui/icons/Close.svelte';
	import Phone from 'components/ui/icons/Phone.svelte';
	import HashChat from 'components/ui/icons/HashChat.svelte';
	import { backend } from 'stores/backend.svelte';
	import Button from 'components/ui/Button/Button.svelte';
	import { rtc } from 'stores/rtc.svelte';
	import { sounds } from 'stores/audio.svelte';

	let { id, tab, server, channel, friend } = $props();

	$effect(() => {
		if (userStore.mention) {
			setTimeout(() => {
				userStore.mention = false;
			}, 500);
		}
	});

	async function joinCall() {
		const res = await backend.connectToCall(server.id, channel.id);

		if (res.isErr()) {
			console.error(res.error.error);
			return;
		}

		if (res.isOk()) {
			rtc.connectToRoom(res.value.token);
			sounds.playSound('call-on');
		}
	}
</script>

<div id={`window-top-bar-${id}`} class="flex gap-x-0.5 hover:cursor-grab active:cursor-grabbing">
	<div
		class={[
			'inner-main-800 relative flex h-[2.375rem] w-full items-center justify-between px-2.5 transition duration-100',
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
						class="pointer-events-none h-[1.25rem] w-[1.25rem] rounded-full"
					/>
					<p class="select-none">{friend.display_name}</p>
				</div>
			{:else if server}
				<div class="flex items-center gap-x-1.5">
					<img
						src={server.avatar}
						alt="{server.name} server background"
						class="pointer-events-none h-[1.25rem] w-[1.25rem] rounded-full"
					/>
					<p class="select-none">{server.name}</p>
				</div>
			{/if}
			{#if channel && !friend}
				|
				<p class="select-none">#{channel.name}</p>
			{/if}
		</div>
	</div>

	<Button
		variants="icon"
		class={[
			'inner-main-800',
			tab !== 'chat'
				? 'hocus:inner-main-700-shadow'
				: 'hocus:inner-green-400/40 hocus:text-green-400'
		]}
		onclick={() => {
			if (!userStore.currentVoiceChannel) {
				joinCall();
			}
			windows.toggleCallTab();
		}}
		tooltip={tab !== 'chat' ? 'Go to chat' : 'Join call'}
		corners
		cornerClass={tab !== 'chat' ? 'group-hocus:border-main-600' : 'group-hocus:border-green-400'}
	>
		{#if tab === 'chat'}
			<Phone height={14} width={14} />
		{:else if tab === 'call'}
			<HashChat height={14} width={14} />
		{/if}
	</Button>
	<Button
		class="inner-main-800 hocus:inner-main-700-shadow"
		variants="icon"
		onclick={() => windows.closeWindow(id)}
		tooltip="Close chat"
		corners
		cornerClass="group-hocus:border-main-600"
	>
		<Close height={16} width={16} />
	</Button>
</div>
