<script lang="ts">
	import '@fontsource/nunito/400.css';
	import '@fontsource/nunito/500.css';
	import '@fontsource/nunito/600.css';
	import '@fontsource/nunito/700.css';
	import '@fontsource/vollkorn/400-italic.css';
	import '@fontsource/outfit/700.css';
	import '@fontsource/outfit/800.css';
	import '@fontsource/outfit/900.css';
	import '../app.css';
	import { goback } from '../stores/goback.svelte';
	import { page } from '$app/state';
	import Searchbar from '../components/Searchbar/Searchbar.svelte';
	import Userbar from '../components/Userbar/Userbar.svelte';
	import Topbar from '../components/Topbar/Topbar.svelte';
	import GridDots from '../components/GridDots.svelte';
	import { windows } from '../stores/windows.svelte';
	import ChatWindow from '../components/ChatWindow/ChatWindow.svelte';

	let { children } = $props();

	$effect(() => {
		if (page.url.pathname === '/') {
			goback.off();
		} else {
			goback.on();
		}
	});
</script>

<Topbar canGoBack={goback.active} />
<Userbar />
<Searchbar />

<GridDots />

{@render children()}

{#each windows.openWindows as chatWindow (chatWindow.id)}
	<ChatWindow id={chatWindow.id} serverId={chatWindow.serverId} channelId={chatWindow.channelId} />
{/each}
