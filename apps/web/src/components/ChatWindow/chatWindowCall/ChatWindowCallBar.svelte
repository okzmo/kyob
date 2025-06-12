<script lang="ts">
	import Corners from 'components/ui/Corners/Corners.svelte';
	import Microphone from 'components/ui/icons/Microphone.svelte';
	import Camera from 'components/ui/icons/Camera.svelte';
	import ScreenShare from 'components/ui/icons/ScreenShare.svelte';
	import PhoneDown from 'components/ui/icons/PhoneDown.svelte';
	import { backend } from 'stores/backend.svelte';
	import { windows } from 'stores/windows.svelte';
	import Button from 'components/ui/Button/Button.svelte';

	let { server, channel } = $props();
</script>

<div class="absolute bottom-4 left-1/2 flex -translate-x-1/2 gap-x-2">
	<div class="bg-main-900/50 inner-main-800 relative flex items-center p-1.5 backdrop-blur-lg">
		<Corners color="border-main-700" />
		<Button
			variants="nostyle"
			class="hocus:bg-main-800 hocus:inner-main-700-shadow text-main-400 hocus:text-main-50 px-3 py-2 transition-colors duration-100 hover:cursor-pointer"
			onclick={() => {}}
			tooltip="Mute"
		>
			<Microphone height={24} width={24} />
		</Button>
		<Button
			variants="nostyle"
			class="hocus:bg-main-800 hocus:inner-main-700-shadow text-main-400 hocus:text-main-50 px-3 py-2 transition-colors duration-100 hover:cursor-pointer"
			onclick={() => {}}
			tooltip="Turn On Camera"
		>
			<Camera height={24} width={24} />
		</Button>
		<Button
			variants="nostyle"
			class="hocus:bg-main-800 hocus:inner-main-700-shadow text-main-400 hocus:text-main-50 px-3 py-2 transition-colors duration-100 hover:cursor-pointer"
			onclick={() => {}}
			tooltip="Share Your Screen"
		>
			<ScreenShare height={24} width={24} />
		</Button>
	</div>

	<Button
		variants="nostyle"
		class="inner-red-400/20 hocus:inner-red-400/40 hocus:bg-red-400/30 relative h-full bg-red-400/20 px-5.5 py-2 text-red-400 backdrop-blur-lg transition hover:cursor-pointer"
		onclick={() => {
			backend.disconnectFromCall(server.id, channel.id);
			windows.toggleCallTab();
		}}
		tooltip="Disconnect"
		corners
		cornerColor="border-red-400/50"
		cornerClass="hocus:border-red-400/80"
	>
		<PhoneDown height={24} width={24} />
	</Button>
</div>
