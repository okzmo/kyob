<script lang="ts">
	import Button from 'components/ui/Button/Button.svelte';
	import Phone from 'components/ui/icons/Phone.svelte';
	import PhoneDown from 'components/ui/icons/PhoneDown.svelte';
	import { sounds } from 'stores/audio.svelte';
	import { backend } from 'stores/backend.svelte';
	import { core } from 'stores/core.svelte';
	import { rtc } from 'stores/rtc.svelte';
	import { userStore } from 'stores/user.svelte';
	import { windows } from 'stores/windows.svelte';
	import { fly } from 'svelte/transition';

	async function acceptCall() {
		if (!core.friendCalling) return;
		const channelId = core.friendCalling.channelId;
		const friendId = core.friendCalling.friendId;
		core.friendCalling = undefined;
		sounds.stopSound('ring-tone');

		if (rtc.currentVC) {
			await backend.disconnectFromCall('global', rtc.currentVC.channelId);
			await rtc.quitRoom();
		}

		const res = await backend.connectToCall('global', channelId);

		if (res.isOk()) {
			await rtc.connectToRoom(res.value.token, 'global', channelId);
			windows.createWindow({
				id: `window-${friendId}`,
				serverId: 'global',
				channelId: channelId,
				friendId: friendId,
				tab: 'call'
			});
		}
	}

	async function refuseCall() {
		if (!core.friendCalling) return;
		const channelId = core.friendCalling.channelId;
		core.friendCalling = undefined;
		sounds.stopSound('ring-tone');

		await backend.disconnectFromCall('global', channelId);
	}
</script>

{#if core.friendCalling}
	{@const friend = userStore.getFriend(core.friendCalling.friendId)}
	<div
		class="mb-3 flex w-full flex-col items-center justify-center gap-y-3 px-2"
		in:fly={{ delay: 150, y: -5, duration: 200 }}
		out:fly={{ y: -5, duration: 100 }}
	>
		<div class="flex items-center gap-x-2">
			<img src={friend?.avatar} class="h-[1.75rem] w-[1.75rem]" alt="{friend?.username}'s avatar" />
			<p>{friend?.display_name} is calling</p>
		</div>
		<div class="flex w-full gap-x-2">
			<Button
				variants="green"
				class="flex flex-grow items-center justify-center gap-x-2"
				onclick={acceptCall}
			>
				<Phone height={16} width={16} />
				Accept
			</Button>
			<Button
				variants="danger"
				class="flex flex-grow items-center justify-center gap-x-2"
				onclick={refuseCall}
			>
				<PhoneDown height={16} width={16} />
				Refuse
			</Button>
		</div>
	</div>
{/if}
