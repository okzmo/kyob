<script lang="ts">
	import { windows } from 'stores/windows.svelte';
	import Corners from '../ui/Corners/Corners.svelte';
	import Close from '../ui/icons/Close.svelte';
	import Phone from 'components/ui/icons/Phone.svelte';
	import HashChat from 'components/ui/icons/HashChat.svelte';
	import Button from 'components/ui/Button/Button.svelte';
	import { rtc } from 'stores/rtc.svelte';
	import { backend } from 'stores/backend.svelte';
	import { sounds } from 'stores/audio.svelte';
	import { core } from 'stores/core.svelte';
	import { setControllableTimeout } from 'utils/timeout';
	import type { Channel, Friend, Server } from 'types/types';

	interface Props {
		id: string;
		tab: 'chat' | 'call';
		type: 'default' | 'world';
		server: Server;
		channel?: Channel;
		friend?: Friend;
	}

	let { id, tab, type, server, channel, friend }: Props = $props();

	let isFriendChannel = $derived(server.id === 'global' && friend);
	let isFriendCallActive = $derived(channel && channel.voice_users.length > 0);

	async function handleCallTab() {
		if (!channel) return;

		if (isFriendChannel && !isFriendCallActive) {
			if (rtc.currentVC) {
				await backend.disconnectFromCall(rtc.currentVC.serverId, rtc.currentVC.channelId);
				await rtc.quitRoom();
			}

			const res = await backend.connectToCall(server.id, channel.id);

			if (res.isErr()) {
				console.error(res.error.error);
				return;
			}

			if (res.isOk()) {
				await rtc.connectToRoom(res.value.token, server.id, channel.id);
				windows.toggleCallTab(channel.id);

				core.callTimeout = setControllableTimeout(() => {
					backend.disconnectFromCall(server.id, channel.id);
					rtc.quitRoom();
					sounds.playSound('call-off');
					windows.goToChatTab(channel.id);
				}, 10000);
			}
		} else {
			windows.toggleCallTab(channel.id);
		}
	}
</script>

<div id={`top-bar-${id}`} class="flex gap-x-0.5 hover:cursor-grab active:cursor-grabbing">
	<div
		class="inner-main-800 bg-main-900 relative flex h-[2.375rem] flex-grow items-center justify-between px-2.5 transition duration-100"
	>
		<Corners color="border-main-700" />
		<div
			class="text-main-400 flex items-center gap-x-2 text-sm transition-colors duration-100 select-none"
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

	{#if type === 'default' && channel}
		<Button
			variants="icon"
			class={[
				'inner-main-800 !w-auto gap-x-2 px-[0.7rem]',
				channel.voice_users.length > 0 && tab === 'chat' ? '!aspect-auto' : '',
				tab !== 'chat'
					? 'hocus:inner-main-700-shadow'
					: 'hocus:inner-green-400/40 hocus:text-green-400'
			]}
			onclick={handleCallTab}
			tooltip={tab !== 'chat'
				? 'Go to chat'
				: isFriendChannel && !isFriendCallActive
					? 'Start call'
					: 'Go to voice chat'}
			corners
			cornerClass={tab !== 'chat' ? 'group-hocus:border-main-600' : 'group-hocus:border-green-400'}
		>
			{#if tab === 'chat'}
				<Phone height={16} width={16} />
			{:else if tab === 'call'}
				<HashChat height={16} width={16} />
			{/if}

			{#if channel.voice_users.length > 0 && tab === 'chat' && server.id !== 'global'}
				{channel.voice_users.length}
			{/if}
		</Button>
	{/if}
	<Button
		class="inner-main-800 hocus:inner-main-700-shadow px-[0.7rem]"
		variants="icon"
		onclick={() => windows.closeWindow(id)}
		tooltip="Close chat"
		corners
		cornerClass="group-hocus:border-main-600"
	>
		<Close height={16} width={16} />
	</Button>
</div>
