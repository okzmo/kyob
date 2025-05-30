<script lang="ts">
	import FriendsButton from './friends/FriendsButton.svelte';
	import GoBack from './goback/GoBack.svelte';
	import NotificationsButton from './notifications/NotificationsButton.svelte';

	interface Props {
		canGoBack: boolean;
	}

	let { canGoBack = false }: Props = $props();
</script>

<div
	class={[
		'fixed top-5 z-50 flex w-screen items-center px-5',
		canGoBack ? 'justify-between' : 'justify-end'
	]}
>
	{#if canGoBack}
		<GoBack />
	{/if}
	<div class="flex items-center gap-x-4">
		<NotificationsButton />
		<FriendsButton />
	</div>
</div>
