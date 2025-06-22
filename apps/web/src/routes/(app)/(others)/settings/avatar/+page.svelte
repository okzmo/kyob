<script lang="ts">
	import { userStore } from 'stores/user.svelte';
</script>

<h1 class="text-2xl font-bold select-none">Avatar</h1>
<p class="text-main-400 mt-3">
	Powered by Ready Player Me, your avatar is your way to express yourself on Kyob. You can change it
	at any time. When you're done simply click on <span class="font-bold">Next</span> and you're done!
	Go back to your realms and show your slick avatar.
</p>

{#if userStore.user}
	<iframe
		title=""
		src={`https://kyob.readyplayer.me/avatar?frameApi&id=${userStore.user.rpm_avatar_id}&token=${userStore.user.rpm_token}`}
		id="frame"
		class="border-main-500 mt-10 h-[40rem] w-full border"
		allow="camera *; microphone *; clipboard-write"
	></iframe>
{:else}
	Loading...
{/if}
