<script lang="ts">
	import Corners from 'components/ui/Corners/Corners.svelte';
	import Microphone from 'components/ui/icons/Microphone.svelte';
	import Camera from 'components/ui/icons/Camera.svelte';
	import ScreenShare from 'components/ui/icons/ScreenShare.svelte';
	import PhoneDown from 'components/ui/icons/PhoneDown.svelte';
	import { backend } from 'stores/backend.svelte';
	import { windows } from 'stores/windows.svelte';

	let { server, channel } = $props();
</script>

<div class="absolute bottom-4 left-1/2 flex -translate-x-1/2 gap-x-2">
	<div class="bg-main-900/50 inner-main-800 relative flex items-center p-1.5 backdrop-blur-lg">
		<Corners color="border-main-700" />
		<button
			class="hocus:bg-main-800 hocus:inner-main-700-shadow text-main-400 hocus:text-main-50 px-3 py-2 transition-colors duration-100 hover:cursor-pointer"
		>
			<Microphone height={24} width={24} />
		</button>
		<button
			class="hocus:bg-main-800 hocus:inner-main-700-shadow text-main-400 hocus:text-main-50 px-3 py-2 transition-colors duration-100 hover:cursor-pointer"
		>
			<Camera height={24} width={24} />
		</button>
		<button
			class="hocus:bg-main-800 hocus:inner-main-700-shadow text-main-400 hocus:text-main-50 px-3 py-2 transition-colors duration-100 hover:cursor-pointer"
		>
			<ScreenShare height={24} width={24} />
		</button>
	</div>
	<button
		class="inner-red-400/20 hocus:inner-red-400/40 hocus:bg-red-400/30 relative bg-red-400/20 px-5.5 py-2 text-red-400 backdrop-blur-lg transition hover:cursor-pointer"
		onclick={() => {
			backend.disconnectFromCall(server.id, channel.id);
			windows.toggleCallTab();
		}}
	>
		<Corners color="border-red-400/50" class="hocus:border-red-400/80" />
		<PhoneDown height={24} width={24} />
	</button>
</div>
