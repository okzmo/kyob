<script lang="ts">
	interface Props {
		class?: string[];
		color?: string;
		hide?: boolean;
	}

	let { color = 'border-white', class: classes, hide }: Props = $props();
</script>

<div class="pointer-events-none absolute inset-0 z-10">
	<span
		class={[
			'absolute top-0 left-0 block h-[7px] w-[7px] border-t border-l transition',
			color,
			hide && 'group-hocus:opacity-100 opacity-0',
			classes
		]}
	></span>
	<span
		class={[
			'absolute top-0 right-0 block h-[7px] w-[7px] border-t border-r transition',
			color,
			hide && 'group-hocus:opacity-100 opacity-0',
			classes
		]}
	></span>
	<span
		class={[
			'absolute bottom-0 left-0 block h-[7px] w-[7px] border-b border-l transition',
			color,
			hide && 'group-hocus:opacity-100 opacity-0',
			classes
		]}
	></span>
	<span
		class={[
			'absolute right-0 bottom-0 block h-[7px] w-[7px] border-r border-b transition',
			color,
			hide && 'group-hocus:opacity-100 opacity-0',
			classes
		]}
	></span>
</div>
