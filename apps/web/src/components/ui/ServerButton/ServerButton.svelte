<script lang="ts">
	import { goto } from '$app/navigation';

	interface Props {
		name: string;
		background: string;
		href: string;
		x: number;
		y: number;
	}

	let { name, background, href, x, y }: Props = $props();
</script>

<button
	class="group absolute h-[4rem] w-[4rem] hover:cursor-pointer"
	aria-label={`${name} server background`}
	style="transform: translate({x}px, {y}px);"
	onclick={() => {
		goto(`/${href}`);
	}}
>
	<img
		src={background}
		alt="Server background"
		class="transition-radius group-hocus:rounded-2xl rounded-[50%]"
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
		box-shadow: 0px 0px 0px 0px #6a6a7c;
		transition:
			border-radius 350ms cubic-bezier(0.65, 0.05, 0, 1),
			box-shadow 350ms cubic-bezier(0.65, 0.05, 0, 1);
	}

	button:hover::before {
		border-radius: 16px;
		box-shadow: 0px 0px 0px 2.5px #6a6a7c;
	}

	button:focus-visible::before {
		border-radius: 16px;
		box-shadow: 0px 0px 0px 2.5px #6a6a7c;
	}

	.transition-radius {
		transition: border-radius 350ms cubic-bezier(0.65, 0.05, 0, 1);
	}
</style>
