<script lang="ts">
	import Corners from 'components/ui/Corners/Corners.svelte';
	import Microphone from 'components/ui/icons/Microphone.svelte';
	import Camera from 'components/ui/icons/Camera.svelte';
	import ScreenShare from 'components/ui/icons/ScreenShare.svelte';
	import PhoneDown from 'components/ui/icons/PhoneDown.svelte';
	import { backend } from 'stores/backend.svelte';
	import { windows } from 'stores/windows.svelte';
	import Button from 'components/ui/Button/Button.svelte';
	import { rtc } from 'stores/rtc.svelte';
	import { userStore } from 'stores/user.svelte';
	import { sounds } from 'stores/audio.svelte';

	let { server, channel } = $props();

	async function joinCall() {
		if (rtc.currentVC) return;

		const res = await backend.connectToCall(server.id, channel.id);

		if (res.isErr()) {
			console.error(res.error.error);
			return;
		}

		if (res.isOk()) {
			await rtc.connectToRoom(res.value.token);
		}
	}

	async function leaveCall() {
		await rtc.quitRoom();
		backend.disconnectFromCall(server.id, channel.id);
		sounds.playSound('call-off');
		windows.toggleCallTab();
	}
</script>

<div class="absolute bottom-4 left-1/2 flex -translate-x-1/2 gap-x-2">
	{#if !rtc.currentVC}
		<Button
			variants="nostyle"
			class="inner-green-400/20 hocus:inner-green-400/40 hocus:bg-green-400/30 relative bg-green-400/20 px-6.5 py-2 text-green-400 backdrop-blur-lg transition hover:cursor-pointer"
			onclick={joinCall}
			corners
			cornerColor="border-green-400/50"
			cornerClass="hocus:border-green-400/80"
		>
			Join channel
		</Button>
	{:else}
		<div class="bg-main-900/50 inner-main-800 relative flex items-center p-1.5 backdrop-blur-lg">
			<Corners color="border-main-700" />
			<Button
				variants="nostyle"
				class={[
					'relative px-3 py-2 transition-colors duration-100 hover:cursor-pointer',
					userStore.mute
						? 'inner-red-400/20 hocus:bg-red-400/30 hocus:inner-red-400/40 bg-red-400/20 text-red-400'
						: 'text-main-400 hocus:bg-main-800 hocus:inner-main-700-shadow hocus:text-main-50'
				]}
				onclick={() => userStore.toggleMute()}
				tooltip="Mute"
			>
				<Microphone height={24} width={24} mute={userStore.mute} />
			</Button>
			<Button
				variants="nostyle"
				class="hocus:bg-main-800 hocus:inner-main-700-shadow text-main-400 hocus:text-main-50 relative px-3 py-2 transition-colors duration-100 hover:cursor-pointer"
				onclick={() => {}}
				tooltip="Turn On Camera"
			>
				<Camera height={24} width={24} />
			</Button>
			<Button
				variants="nostyle"
				class="hocus:bg-main-800 hocus:inner-main-700-shadow text-main-400 hocus:text-main-50 relative px-3 py-2 transition-colors duration-100 hover:cursor-pointer"
				onclick={() => {}}
				tooltip="Share Your Screen"
			>
				<ScreenShare height={24} width={24} />
			</Button>
		</div>

		<Button
			variants="nostyle"
			class="inner-red-400/20 hocus:inner-red-400/40 hocus:bg-red-400/30 relative bg-red-400/20 px-6.5 py-2 text-red-400 backdrop-blur-lg transition hover:cursor-pointer"
			onclick={leaveCall}
			tooltip="Disconnect"
			corners
			cornerColor="border-red-400/50"
			cornerClass="hocus:border-red-400/80"
		>
			<PhoneDown height={24} width={24} />
		</Button>
	{/if}
</div>
