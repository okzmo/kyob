<script lang="ts">
	import { goto } from '$app/navigation';
	import IconButton from '../ui/IconButton/IconButton.svelte';
	import ArrowLeft from '../ui/icons/ArrowLeft.svelte';
	import Bell from '../ui/icons/Bell.svelte';
	import People from '../ui/icons/People.svelte';
	import FriendsButton from './friends/FriendsButton.svelte';

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
		<button
			class="text-main-400 group hocus:bg-accent-100/15 flex items-center gap-x-4 rounded-xl py-1.5 pr-4 pl-3 transition-colors duration-100 hover:cursor-pointer"
			onclick={() => goto('/')}
		>
			<ArrowLeft
				height={22}
				width={22}
				class="group-hocus:text-accent-50 transition-colors duration-100"
			/>
			<p class="group-hocus:text-accent-50 text-lg transition-colors duration-100">Go back</p>
		</button>
	{/if}
	<div class="flex items-center gap-x-4">
		<IconButton Icon={Bell} label="Notifications" />
		<FriendsButton />
	</div>
</div>
