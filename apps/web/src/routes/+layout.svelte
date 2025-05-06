<script lang="ts">
	import '@fontsource/nunito/400.css';
	import '@fontsource/nunito/500.css';
	import '@fontsource/nunito/600.css';
	import '@fontsource/nunito/700.css';
	import '@fontsource/vollkorn/400-italic.css';
	import '@fontsource/outfit/700.css';
	import '@fontsource/outfit/800.css';
	import '@fontsource/outfit/900.css';
	import '../app.css';
	import GridDots from '../components/GridDots.svelte';

	let { children } = $props();
</script>

{@render children()}

<GridDots />
