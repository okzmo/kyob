<script lang="ts">
	interface Props {
		height: number;
		width: number;
		class?: string;
	}

	let { class: classes, height = 24, width = 24 }: Props = $props();
</script>

<svg
	{width}
	{height}
	class={classes}
	viewBox="0 0 20 20"
	fill="none"
	xmlns="http://www.w3.org/2000/svg"
>
	<path
		d="M16.6666 10H3.33331M3.33331 10L8.33331 5M3.33331 10L8.33331 15"
		stroke="currentColor"
		stroke-width="2"
		stroke-linecap="round"
		stroke-linejoin="round"
	/>
</svg>
