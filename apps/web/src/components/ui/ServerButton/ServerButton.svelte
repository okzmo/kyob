<script lang="ts">
	import { goto } from '$app/navigation';
	import { elasticInOut } from 'svelte/easing';
	import { scale } from 'svelte/transition';
	import Corners from '../Corners/Corners.svelte';

	interface Props {
		id: string;
		name: string;
		avatar: string;
		href: string;
		x: number;
		y: number;
	}

	let { id, name, avatar, href, x, y }: Props = $props();
</script>

<button
	transition:scale={{ start: 0, duration: 800, easing: elasticInOut }}
	id="serverButton-{id}"
	class="group absolute h-[4rem] w-[4rem] hover:cursor-pointer"
	aria-label={`${name} server background`}
	style="transform: translate({x}px, {y}px);"
	onclick={() => {
		goto(`/${href}`);
	}}
>
	<Corners color="border-accent-100" hide class="duration-100" />
	<img src={avatar} alt={name.slice(0, 2).toUpperCase()} class="h-full w-full object-cover" />
</button>

<style>
	button::before {
		content: '';
		height: 100%;
		width: 100%;
		position: absolute;
		left: 0;
		top: 0;
		box-shadow:
			inset 0px 0px 0px 1px var(--color-main-700),
			inset 0px 0px 12px var(--ui-accent-10000);
		background-color: var(--ui-accent-10000);
		transition:
			background-color 100ms ease-out,
			box-shadow 100ms ease-out;
	}

	button:hover::before {
		box-shadow:
			inset 0px 0px 0px 1px var(--ui-accent-10050),
			inset 0px 0px 12px var(--ui-accent-10050);
		background-color: var(--ui-accent-10020);
	}

	button:focus-visible::before {
		box-shadow:
			inset 0px 0px 0px 1px var(--ui-accent-10050),
			inset 0px 0px 12px var(--ui-accent-10050);
		background-color: var(--ui-accent-10020);
	}
</style>
