<script lang="ts">
	import CreateChannelModal from 'components/modals/CreateChannelModal.svelte';
	import { core } from 'stores/core.svelte';
	import { windows } from 'stores/windows.svelte';
	import { onMount } from 'svelte';
	import type { Channel } from 'types/types';
	import ChannelButton from '../ui/ChannelButton/ChannelButton.svelte';

	let dragging = $state(false);
	let startPos = $state({ x: 0, y: 0 });
	let offset = $state({ x: 0, y: 0 });
	let totalOffset = $state({ x: 0, y: 0 });
	let velocity = $state({ x: 0, y: 0 });
	let lastMousePos = $state({ x: 0, y: 0 });
	let lastTimestamp = $state(0);
	let animationFrameId = $state<number | null>(null);
	let dragDistance = $state(0);
	let dragStartTime = $state(0);

	function handleMouseDown(e: MouseEvent) {
		if (e.buttons !== 1 || !core.canDragMap) return;
		dragging = true;
		windows.setActiveWindow(null);
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
		if (!dragging) return;

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

		offset = {
			x: totalOffset.x + totalDx,
			y: totalOffset.y + totalDy
		};
	}

	function handleMouseUp() {
		if (dragging) {
			totalOffset = { ...offset };
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

			offset = {
				x: offset.x + velocity.x,
				y: offset.y + velocity.y
			};

			totalOffset = { ...offset };

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
		channels: Channel[];
	}

	let { channels }: Props = $props();
</script>

{#if channels.length > 0}
	{#each channels as channel (channel.id)}
		<ChannelButton
			id={channel.id}
			name={channel.name}
			type={channel.type}
			x={channel.x + offset.x}
			y={channel.y + offset.y}
			unread={channel.unread}
		/>
	{/each}
{:else}
	<h3
		class="font-outfit text-main-600 fixed top-1/2 left-1/2 -translate-1/2 text-5xl font-bold uppercase select-none"
	>
		No channels yet
	</h3>
{/if}

<CreateChannelModal />
