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
	<circle cx="12" cy="6" r="4" fill="currentColor" />
	<path
		fill="currentColor"
		d="M20 17.5c0 2.485 0 4.5-8 4.5s-8-2.015-8-4.5S7.582 13 12 13s8 2.015 8 4.5"
		opacity=".5"
	/>
</svg>
