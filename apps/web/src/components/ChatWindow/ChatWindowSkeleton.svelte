<script lang="ts">
	import { fly } from 'svelte/transition';
	import { onMount, type Snippet } from 'svelte';
	import { windows } from 'stores/windows.svelte';
	import type { Channel, Friend, Server } from 'types/types';
	import ChatWindowTopBar from './ChatWindowTopBar.svelte';
	import Corners from '../ui/Corners/Corners.svelte';

	interface Props {
		id: string;
		tab: 'chat' | 'call';
		type: 'default' | 'world';
		server: Server;
		channel: Channel;
		friend?: Friend;
		children: Snippet;
	}

	type Positions = 'br' | 'bl' | 'tl' | 'tr' | 'l' | 'r' | 'b' | 't' | '';

	let { id, children, tab, type, server, channel, friend }: Props = $props();

	let windowState = $state(windows.openWindows.find((w) => w.id === id)!);
	let startPos = $state({ x: 0, y: 0 });
	let offset = $state({ x: windowState.x, y: windowState.y });
	let dragging = $state(false);
	let resizingPos = $state<Positions>('');

	let resizing = $state(false);
	let startPosResizing = $state({ x: 0, y: 0 });
	let initialSize = $state({ width: 0, height: 0 });

	function chatTopBarMouseDown(e: MouseEvent) {
		dragging = true;
		startPos = { x: e.clientX, y: e.clientY };

		document.addEventListener('mouseup', chatTopBarMouseUp);
		document.addEventListener('mousemove', chatTopBarMouseMove);
	}

	function chatTopBarMouseUp() {
		dragging = false;
		windowState.x = offset.x;
		windowState.y = offset.y;

		document.removeEventListener('mouseup', chatTopBarMouseUp);
		document.removeEventListener('mousemove', chatTopBarMouseMove);
	}

	function chatTopBarMouseMove(e: MouseEvent) {
		if (!dragging) return;
		e.preventDefault();

		const dx = e.clientX - startPos.x;
		const dy = e.clientY - startPos.y;

		offset = {
			x: Math.max(0, Math.min(windowState.x + dx, window.innerWidth - windowState.width)),
			y: Math.max(0, Math.min(windowState.y + dy, window.innerHeight - windowState.height - 38))
		};
	}

	function chatResizeMouseDown(e: MouseEvent, pos: Positions) {
		resizing = true;
		startPosResizing = { x: e.clientX, y: e.clientY };
		initialSize = { height: windowState?.height, width: windowState?.width };
		resizingPos = pos;

		document.addEventListener('mouseup', stopResizing);
		document.addEventListener('mousemove', chatResizing);
	}

	function chatResizing(e: MouseEvent) {
		if (!resizing) return;
		e.preventDefault();

		const dx = e.clientX - startPosResizing.x;
		const dy = e.clientY - startPosResizing.y;

		switch (resizingPos) {
			case 'br':
				windowState.width = Math.max(450, initialSize.width + dx);
				windowState.height = Math.max(350, initialSize.height + dy);
				break;
			case 'bl':
				{
					const newWidth = Math.max(450, initialSize.width - dx);
					const diffWidth = newWidth - windowState.width;

					windowState.width = newWidth;
					windowState.height = Math.max(350, initialSize.height + dy);
					offset.x = offset.x - diffWidth;
				}
				break;
			case 'tl':
				{
					const newWidth = Math.max(450, initialSize.width - dx);
					const diffWidth = newWidth - windowState.width;
					const newHeight = Math.max(350, initialSize.height - dy);
					const diffHeight = newHeight - windowState.height;

					windowState.width = newWidth;
					windowState.height = newHeight;
					offset.x = offset.x - diffWidth;
					offset.y = offset.y - diffHeight;
				}
				break;
			case 'tr':
				{
					const newHeight = Math.max(350, initialSize.height - dy);
					const diffHeight = newHeight - windowState.height;

					windowState.width = Math.max(450, initialSize.width + dx);
					windowState.height = newHeight;
					offset.y = offset.y - diffHeight;
				}
				break;
			case 't':
				{
					const newHeight = Math.max(350, initialSize.height - dy);
					const diffHeight = newHeight - windowState.height;

					windowState.height = newHeight;
					offset.y = offset.y - diffHeight;
				}
				break;
			case 'b':
				windowState.height = Math.max(350, initialSize.height + dy);
				break;
			case 'r':
				windowState.width = Math.max(450, initialSize.width + dx);
				break;
			case 'l':
				{
					const newWidth = Math.max(450, initialSize.width - dx);
					const diffWidth = newWidth - windowState.width;

					windowState.width = newWidth;
					offset.x = offset.x - diffWidth;
				}
				break;
		}
	}

	function stopResizing() {
		resizing = false;
		windowState.x = offset.x;
		windowState.y = offset.y;
		resizingPos = '';

		document.removeEventListener('mouseup', stopResizing);
		document.removeEventListener('mousemove', chatResizing);
	}

	function windowMouseDown(e: MouseEvent) {
		e.stopImmediatePropagation();
		windows.setActiveWindow(id);
	}

	onMount(() => {
		const fullWindow = document.getElementById(id);
		const windowBar = document.getElementById(`top-bar-${id}`);
		const windowResizeBR = document.getElementById(`resize-br-${id}`);
		const windowResizeBL = document.getElementById(`resize-bl-${id}`);
		const windowResizeTL = document.getElementById(`resize-tl-${id}`);
		const windowResizeTR = document.getElementById(`resize-tr-${id}`);
		const windowResizeT = document.getElementById(`resize-t-${id}`);
		const windowResizeB = document.getElementById(`resize-b-${id}`);
		const windowResizeL = document.getElementById(`resize-l-${id}`);
		const windowResizeR = document.getElementById(`resize-r-${id}`);

		if (!windowBar || !fullWindow) return;
		if (!windowResizeT || !windowResizeB || !windowResizeL || !windowResizeR) return;
		if (!windowResizeTL || !windowResizeTR || !windowResizeBL || !windowResizeBR) return;

		fullWindow?.addEventListener('mousedown', windowMouseDown);
		windowBar.addEventListener('mousedown', chatTopBarMouseDown);
		windowResizeBR.addEventListener('mousedown', (e) => chatResizeMouseDown(e, 'br'));
		windowResizeBL.addEventListener('mousedown', (e) => chatResizeMouseDown(e, 'bl'));
		windowResizeTL.addEventListener('mousedown', (e) => chatResizeMouseDown(e, 'tl'));
		windowResizeTR.addEventListener('mousedown', (e) => chatResizeMouseDown(e, 'tr'));
		windowResizeT.addEventListener('mousedown', (e) => chatResizeMouseDown(e, 't'));
		windowResizeB.addEventListener('mousedown', (e) => chatResizeMouseDown(e, 'b'));
		windowResizeL.addEventListener('mousedown', (e) => chatResizeMouseDown(e, 'l'));
		windowResizeR.addEventListener('mousedown', (e) => chatResizeMouseDown(e, 'r'));

		return () => {
			windowBar.removeEventListener('mousedown', chatTopBarMouseDown);
			windowResizeBR.removeEventListener('mousedown', (e) => chatResizeMouseDown(e, 'br'));
			windowResizeBL.removeEventListener('mousedown', (e) => chatResizeMouseDown(e, 'bl'));
			windowResizeTL.removeEventListener('mousedown', (e) => chatResizeMouseDown(e, 'tl'));
			windowResizeTR.removeEventListener('mousedown', (e) => chatResizeMouseDown(e, 'tr'));
			windowResizeT.removeEventListener('mousedown', (e) => chatResizeMouseDown(e, 't'));
			windowResizeB.removeEventListener('mousedown', (e) => chatResizeMouseDown(e, 'b'));
			windowResizeL.removeEventListener('mousedown', (e) => chatResizeMouseDown(e, 'l'));
			windowResizeR.removeEventListener('mousedown', (e) => chatResizeMouseDown(e, 'r'));
		};
	});
</script>

<div
	in:fly={{ duration: 300, y: 20 }}
	out:fly={{ duration: 150, y: 20 }}
	{id}
	class={[
		'absolute flex flex-col transition-opacity duration-75',
		windows.activeWindow === id ? 'z-[52]' : 'z-[40] opacity-40',
		resizing && 'select-none'
	]}
	style="transform: translate({offset.x}px, {offset.y}px);"
>
	<ChatWindowTopBar {id} {tab} {type} {server} {channel} {friend} />
	<div
		style="width: {windowState?.width}px; height: {windowState?.height}px"
		class="bg-main-900 inner-main-800 relative mt-0.5 flex flex-col items-start overflow-hidden"
	>
		<Corners color="border-main-700" />
		{@render children()}
	</div>
	<div
		id={`resize-br-${id}`}
		class="absolute right-0 bottom-0 h-[0.75rem] w-[0.75rem] hover:cursor-se-resize"
	></div>
	<div
		id={`resize-bl-${id}`}
		class="absolute bottom-0 left-0 h-[0.75rem] w-[0.75rem] hover:cursor-sw-resize"
	></div>
	<div
		id={`resize-tl-${id}`}
		class="absolute top-0 left-0 h-[0.75rem] w-[0.75rem] hover:cursor-nw-resize"
	></div>
	<div
		id={`resize-tr-${id}`}
		class="absolute top-0 right-0 h-[0.5rem] w-[0.5rem] hover:cursor-ne-resize"
	></div>

	<div
		id={`resize-t-${id}`}
		class="absolute top-0 left-[1rem] h-[0.25rem] w-[calc(100%-6.25rem)] hover:cursor-n-resize"
	></div>
	<div
		id={`resize-b-${id}`}
		class="absolute bottom-0 left-[1rem] h-[0.25rem] w-[calc(100%-2rem)] hover:cursor-s-resize"
	></div>
	<div
		id={`resize-l-${id}`}
		class="absolute top-[1rem] left-0 h-[calc(100%-2rem)] w-[0.25rem] hover:cursor-w-resize"
	></div>
	<div
		id={`resize-r-${id}`}
		class="absolute top-[2.5rem] right-0 h-[calc(100%-3.5rem)] w-[0.25rem] hover:cursor-e-resize"
	></div>
</div>
