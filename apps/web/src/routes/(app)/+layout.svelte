<script lang="ts">
	import { page } from '$app/state';
	import { goback } from '../../stores/goback.svelte';
	import Topbar from '../../components/Topbar/Topbar.svelte';
	import Userbar from '../../components/Userbar/Userbar.svelte';
	import Searchbar from '../../components/Searchbar/Searchbar.svelte';
	import { windows } from '../../stores/windows.svelte';
	import ChatWindow from '../../components/ChatWindow/ChatWindow.svelte';

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

{@render children()}

{#each windows.openWindows as chatWindow (chatWindow.id)}
	<ChatWindow id={chatWindow.id} serverId={chatWindow.serverId} channelId={chatWindow.channelId} />
{/each}
