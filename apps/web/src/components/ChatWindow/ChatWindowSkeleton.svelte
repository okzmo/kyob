<script lang="ts">
	import { fly } from 'svelte/transition';
	import { onMount, type Snippet } from 'svelte';
	import { windows } from 'stores/windows.svelte';
	import type { Channel, Friend, Server } from 'types/types';
	import ChatWindowTopBar from './ChatWindowTopBar.svelte';
	import Corners from '../ui/Corners/Corners.svelte';

	let width = $state(550);
	let height = $state(400);
	let startPos = $state({ x: 0, y: 0 });
	let offset = $state({ x: 100, y: 150 });
	let totalOffset = $state({ x: 100, y: 150 });
	let dragging = $state(false);

	let resizing = $state(false);
	let startPosResizing = $state({ x: 0, y: 0 });
	let initialSize = $state({ width: 0, height: 0 });

	interface Props {
		id: string;
		server?: Server;
		channel?: Channel;
		friend?: Friend;
		children: Snippet;
	}

	let { id, children, server, channel, friend }: Props = $props();

	function chatTopBarMouseDown(e: MouseEvent) {
		dragging = true;
		startPos = { x: e.clientX, y: e.clientY };

		document.addEventListener('mouseup', chatTopBarMouseUp);
		document.addEventListener('mousemove', chatTopBarMouseMove);
	}

	function chatTopBarMouseUp() {
		dragging = false;
		totalOffset = { ...offset };

		document.removeEventListener('mouseup', chatTopBarMouseUp);
		document.removeEventListener('mousemove', chatTopBarMouseMove);
	}

	function chatTopBarMouseMove(e: MouseEvent) {
		if (!dragging) return;
		e.preventDefault();

		const dx = e.clientX - startPos.x;
		const dy = e.clientY - startPos.y;

		offset = {
			x: Math.max(0, Math.min(totalOffset.x + dx, window.innerWidth - width)),
			y: Math.max(0, Math.min(totalOffset.y + dy, window.innerHeight - height - 38))
		};
	}

	function chatResizeMouseDown(e: MouseEvent) {
		resizing = true;
		startPosResizing = { x: e.clientX, y: e.clientY };
		initialSize = { width, height };

		document.addEventListener('mouseup', stopResizing);
		document.addEventListener('mousemove', chatResizing);
	}

	function chatResizing(e: MouseEvent) {
		if (!resizing) return;
		e.preventDefault();

		width = Math.max(450, initialSize.width + (e.clientX - startPosResizing.x));
		height = Math.max(250, initialSize.height + (e.clientY - startPosResizing.y));
	}

	function stopResizing() {
		resizing = false;
		totalOffset = { ...offset };

		document.removeEventListener('mouseup', stopResizing);
		document.removeEventListener('mousemove', chatResizing);
	}

	function windowMouseDown(e: MouseEvent) {
		e.stopImmediatePropagation();
		windows.activeWindow = id;
	}

	onMount(() => {
		const fullWindow = document.getElementById(`window-${id}`);
		const windowBar = document.getElementById(`window-top-bar-${id}`);
		const windowResize = document.getElementById(`window-resize-${id}`);

		if (!windowBar || !windowResize || !fullWindow) return;

		fullWindow?.addEventListener('mousedown', windowMouseDown);
		windowBar.addEventListener('mousedown', chatTopBarMouseDown);
		windowResize.addEventListener('mousedown', chatResizeMouseDown);

		return () => {
			windowBar.removeEventListener('mousedown', chatTopBarMouseDown);
			windowResize.removeEventListener('mousedown', chatResizeMouseDown);
		};
	});
</script>

<div
	transition:fly={{ duration: 75, y: 10 }}
	id={`window-${id}`}
	class={['absolute flex flex-col', windows.activeWindow === id ? 'z-[52]' : 'z-[40] opacity-40']}
	style="transform: translate({offset.x}px, {offset.y}px);"
>
	<ChatWindowTopBar {id} {server} {channel} {friend} />
	<div
		style="width: {width}px; height: {height}px"
		class="bg-main-900 inner-main-800 relative mt-0.5 flex flex-col items-start overflow-hidden"
	>
		<Corners color="border-main-700" />
		{@render children()}
		<div
			id={`window-resize-${id}`}
			class="absolute right-0 bottom-0 h-[1rem] w-[1rem] hover:cursor-se-resize"
		></div>
	</div>
</div>
