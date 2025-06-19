<script lang="ts">
	import { windows } from 'stores/windows.svelte';
	import { flyBlur } from 'utils/transition';

	interface Props {
		channelId: string;
		friendId: string;
		username: string;
		avatar: string;
	}

	let { friendId, channelId, avatar, username }: Props = $props();
</script>

<button
	in:flyBlur={{ duration: 125, y: 5 }}
	out:flyBlur={{ duration: 50, y: 5 }}
	onclick={() => {
		windows.createWindow({
			id: `window-${friendId}`,
			serverId: 'global',
			channelId: channelId,
			friendId: friendId
		});
	}}
	class="attachment relative h-[2.25rem] w-[2.25rem] hover:cursor-pointer"
>
	<img src={avatar} alt="{username} dm" />
</button>

<style>
</style>
