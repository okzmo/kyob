<script lang="ts">
	import { onMount } from 'svelte';
	import ServerButton from '../ui/ServerButton/ServerButton.svelte';
	import type { Server } from '../../types/types';
	import { windows } from '../../stores/windows.svelte';
	import { core } from '../../stores/core.svelte';
	import CreateServerModal from '../ui/ContextMenu/ServerMapContextMenu/CreateServerModal.svelte';
	import JoinServerModal from '../ui/ContextMenu/ServerMapContextMenu/JoinServerModal.svelte';

	let dragging = $state(false);
	let startPos = $state({ x: 0, y: 0 });
	let velocity = $state({ x: 0, y: 0 });
	let lastMousePos = $state({ x: 0, y: 0 });
	let lastTimestamp = $state(0);
	let animationFrameId = $state<number | null>(null);
	let dragDistance = $state(0);
	let dragStartTime = $state(0);

	function handleMouseDown(e: MouseEvent) {
		if (e.buttons !== 1 || !core.canDragMap) return;
		dragging = true;
		windows.activeWindow = null;
		startPos = { x: e.clientX, y: e.clientY };
		lastMousePos = { x: e.clientX, y: e.clientY };
		lastTimestamp = Date.now();
		dragStartTime = Date.now();
		dragDistance = 0;

		if (animationFrameId !== null) {
			cancelAnimationFrame(animationFrameId);
			animationFrameId = null;
		}

		velocity = { x: 0, y: 0 };
	}

	function handleMouseMove(e: MouseEvent) {
		if (!dragging || e.buttons !== 1) return;

		const currentTime = Date.now();
		const deltaTime = currentTime - lastTimestamp;

		const dx = e.clientX - lastMousePos.x;
		const dy = e.clientY - lastMousePos.y;
		dragDistance += Math.sqrt(dx * dx + dy * dy);

		if (deltaTime > 10) {
			const maxVelocity = 30;
			velocity = {
				x: Math.min(
					Math.max(((e.clientX - lastMousePos.x) / deltaTime) * 8, -maxVelocity),
					maxVelocity
				),
				y: Math.min(
					Math.max(((e.clientY - lastMousePos.y) / deltaTime) * 8, -maxVelocity),
					maxVelocity
				)
			};

			lastMousePos = { x: e.clientX, y: e.clientY };
			lastTimestamp = currentTime;
		}

		const totalDx = e.clientX - startPos.x;
		const totalDy = e.clientY - startPos.y;

		core.offsetServerMap = {
			x: core.totalOffsetServerMap.x + totalDx,
			y: core.totalOffsetServerMap.y + totalDy
		};
	}

	function handleMouseUp() {
		if (dragging) {
			core.totalOffsetServerMap = { ...core.offsetServerMap };
			dragging = false;

			const dragDuration = Date.now() - dragStartTime;
			const velocityMagnitude = Math.sqrt(velocity.x * velocity.x + velocity.y * velocity.y);
			const shouldApplyInertia = velocityMagnitude > 2.0 && dragDistance > 10 && dragDuration < 300;

			if (shouldApplyInertia) {
				applyInertia();
			}
		}
	}

	function applyInertia() {
		const resistance = 0.95;

		const animate = () => {
			velocity.x *= resistance;
			velocity.y *= resistance;

			core.offsetServerMap = {
				x: core.offsetServerMap.x + velocity.x,
				y: core.offsetServerMap.y + velocity.y
			};

			core.totalOffsetServerMap = { ...core.offsetServerMap };

			if (Math.abs(velocity.x) < 0.1 && Math.abs(velocity.y) < 0.1) {
				animationFrameId = null;
				return;
			}

			animationFrameId = requestAnimationFrame(animate);
		};

		animationFrameId = requestAnimationFrame(animate);
	}

	onMount(() => {
		window.addEventListener('mousedown', handleMouseDown);
		window.addEventListener('mouseup', handleMouseUp);
		window.addEventListener('mousemove', handleMouseMove);
		window.addEventListener('mouseleave', handleMouseUp);

		return () => {
			window.removeEventListener('mousedown', handleMouseDown);
			window.removeEventListener('mouseup', handleMouseUp);
			window.removeEventListener('mousemove', handleMouseMove);
			window.removeEventListener('mouseleave', handleMouseUp);

			if (animationFrameId !== null) {
				cancelAnimationFrame(animationFrameId);
			}
		};
	});

	interface Props {
		servers: Server[];
	}

	let { servers }: Props = $props();
</script>

{#each servers as server (server.id)}
	<ServerButton
		id={server.id}
		name={server.name}
		avatar={server.avatar}
		href={String(server.id)}
		x={server.x + core.offsetServerMap.x}
		y={server.y + core.offsetServerMap.y}
	/>
{/each}

<CreateServerModal />
<JoinServerModal />
