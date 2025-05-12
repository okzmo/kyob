<script lang="ts">
	import { goto } from '$app/navigation';
	import { elasticInOut } from 'svelte/easing';
	import { scale } from 'svelte/transition';

	interface Props {
		id: number;
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
	class="group server-button absolute h-[4rem] w-[4rem] hover:cursor-pointer"
	aria-label={`${name} server background`}
	style="transform: translate({x}px, {y}px);"
	onclick={() => {
		goto(`/${href}`);
	}}
>
	<img
		src={avatar}
		alt={name.slice(0, 2).toUpperCase()}
		class="transition-radius group-hocus:rounded-2xl h-full w-full rounded-[50%] object-cover"
	/>
</button>

<style>
	button::before {
		content: '';
		height: 100%;
		width: 100%;
		position: absolute;
		left: 0;
		top: 0;
		border-radius: 50%;
		box-shadow: 0px 0px 0px 1px #6a6a7c;
		transition:
			border-radius 350ms cubic-bezier(0.65, 0.05, 0, 1),
			box-shadow 350ms cubic-bezier(0.65, 0.05, 0, 1);
	}

	button:hover::before {
		border-radius: 16px;
		box-shadow: 0px 0px 0px 2.5px #fafafa;
	}

	button:focus-visible::before {
		border-radius: 16px;
		box-shadow: 0px 0px 0px 2.5px #fafafa;
	}

	.transition-radius {
		transition: border-radius 350ms cubic-bezier(0.65, 0.05, 0, 1);
	}
</style>
