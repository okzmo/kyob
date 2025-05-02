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
		opacity="0.5"
		d="M16.9277 9.297C16.9277 13.5113 13.5114 16.9276 9.29713 16.9276C5.08293 16.9276 1.66663 13.5113 1.66663 9.297C1.66663 5.0828 5.08293 1.6665 9.29713 1.6665C13.5114 1.6665 16.9277 5.0828 16.9277 9.297Z"
		fill="currentColor"
	/>
	<path
		fill-rule="evenodd"
		clip-rule="evenodd"
		d="M15.6989 15.6989C15.9341 15.4636 16.3156 15.4636 16.5508 15.6989L18.1572 17.3053C18.3925 17.5406 18.3925 17.922 18.1572 18.1572C17.922 18.3925 17.5406 18.3925 17.3053 18.1572L15.6989 16.5508C15.4636 16.3156 15.4636 15.9341 15.6989 15.6989Z"
		fill="currentColor"
	/>
</svg>
