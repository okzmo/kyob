<script lang="ts">
	interface Props {
		offset: { x: number; y: number };
	}

	let { offset }: Props = $props();

	// The size of one grid cell remains constant.
	const gridSize = 38;
</script>

<div
	class="pointer-events-none fixed top-0 left-0 z-[-1] h-screen w-screen
  bg-[radial-gradient(circle,rgba(255,255,255,0.05)_1.5px,transparent_1.5px)]"
	style:background-size="{gridSize}px {gridSize}px"
	style:background-position="{offset.x % gridSize}px {offset.y % gridSize}px"
></div>
