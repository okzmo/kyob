<script lang="ts">
	import { userStore } from 'stores/user.svelte';
	import FriendsButton from './friends/FriendsButton.svelte';
	import GoBack from './goback/GoBack.svelte';
	import DmButton from './friends/DmButton.svelte';
	// import NotificationsButton from './notifications/NotificationsButton.svelte';

	interface Props {
		canGoBack: boolean;
	}

	let { canGoBack = false }: Props = $props();

	let dms = $derived(userStore?.getDms());
</script>

<div
	class={[
		'fixed top-5 z-50 flex w-screen items-start px-5',
		canGoBack ? 'justify-between' : 'justify-end'
	]}
>
	{#if canGoBack}
		<GoBack />
	{/if}
	<div class="flex items-center gap-x-4">
		<!-- <NotificationsButton /> -->
		<div class="flex flex-col gap-y-3">
			<FriendsButton />
			{#each dms as dm (dm.friendId)}
				<DmButton {...dm} />
			{/each}
		</div>
	</div>
</div>
