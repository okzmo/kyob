<script lang="ts">
	import { useTask, T } from '@threlte/core';
	import { OrbitControls } from '@threlte/extras';
	import { Vector3 } from 'three';

	let { avatarPos } = $props();
	let camera = $state<any>();
	let orbitControls = $state<any>();

	let isMoving = $state(false);
	let lastAvatarPos = $state({ x: avatarPos.x, y: avatarPos.y, z: avatarPos.z });
	let movementDirection = $state(new Vector3(0, 0, 1)); // Default behind avatar

	// Camera positioning
	const cameraDistance = 5;
	const cameraHeight = 3;
	const lookAtOffset = { x: 0, y: 1, z: 0 };

	let targetCameraPos = $state(new Vector3());
	let isTransitioning = $state(false);
	let transitionStartTime = $state(0);
	let lastMoveTime = $state(0);

	// Smooth easing function
	function easeInOutCubic(t: number): number {
		return t < 0.5 ? 4 * t * t * t : 1 - Math.pow(-2 * t + 2, 3) / 2;
	}

	$effect(() => {
		const moved =
			lastAvatarPos.x !== avatarPos.x ||
			lastAvatarPos.z !== avatarPos.z ||
			lastAvatarPos.y !== avatarPos.y;

		if (moved) {
			// Calculate movement direction
			const moveDir = new Vector3(
				avatarPos.x - lastAvatarPos.x,
				0,
				avatarPos.z - lastAvatarPos.z
			).normalize();

			if (moveDir.length() > 0) {
				movementDirection = moveDir;
			}

			isMoving = true;
			lastMoveTime = Date.now();

			// Update last position
			lastAvatarPos.x = avatarPos.x;
			lastAvatarPos.y = avatarPos.y;
			lastAvatarPos.z = avatarPos.z;
		} else if (isMoving) {
			// Stop moving after 150ms of no movement
			setTimeout(() => {
				if (Date.now() - lastMoveTime >= 150) {
					isMoving = false;
					startCameraTransition();
				}
			}, 150);
		}
	});

	function startCameraTransition() {
		if (!orbitControls || !camera) return;

		isTransitioning = true;
		transitionStartTime = Date.now();

		// Calculate target position behind avatar based on movement direction
		const behindOffset = movementDirection.clone().multiplyScalar(-cameraDistance);

		targetCameraPos.set(
			avatarPos.x + behindOffset.x,
			avatarPos.y + cameraHeight,
			avatarPos.z + behindOffset.z
		);

		// Temporarily disable orbit controls during transition
		orbitControls.enabled = false;
	}

	useTask(() => {
		if (!camera || !orbitControls) return;

		// Always update orbit controls target
		orbitControls.target.set(
			avatarPos.x + lookAtOffset.x,
			avatarPos.y + lookAtOffset.y,
			avatarPos.z + lookAtOffset.z
		);

		if (isTransitioning) {
			const elapsed = Date.now() - transitionStartTime;
			const duration = 1000; // 1 second transition
			const progress = Math.min(elapsed / duration, 1);
			const easedProgress = easeInOutCubic(progress);

			// Smoothly interpolate camera position
			const startPos = camera.position.clone();
			camera.position.lerpVectors(startPos, targetCameraPos, easedProgress * 0.1);

			// Smoothly look at avatar
			camera.lookAt(
				avatarPos.x + lookAtOffset.x,
				avatarPos.y + lookAtOffset.y,
				avatarPos.z + lookAtOffset.z
			);

			// End transition
			if (progress >= 1) {
				isTransitioning = false;
				orbitControls.enabled = true;
				orbitControls.update();
			}
		} else if (!isMoving && !isTransitioning) {
			// Free camera mode - let orbit controls handle everything
			orbitControls.update();
		} else if (isMoving) {
			// While moving, keep camera behind avatar
			const behindOffset = movementDirection.clone().multiplyScalar(-cameraDistance);

			camera.position.set(
				avatarPos.x + behindOffset.x,
				avatarPos.y + cameraHeight,
				avatarPos.z + behindOffset.z
			);

			camera.lookAt(
				avatarPos.x + lookAtOffset.x,
				avatarPos.y + lookAtOffset.y,
				avatarPos.z + lookAtOffset.z
			);
		}
	});

	function handleControlStart() {
		// User started orbiting - stop any transitions
		isTransitioning = false;
	}

	function handleControlEnd() {
		// Control ended - no immediate action needed
		// Transition will trigger when movement stops
	}
</script>

<T.PerspectiveCamera
	bind:ref={camera}
	makeDefault
	position={[avatarPos.x, avatarPos.y + cameraHeight, avatarPos.z - cameraDistance]}
	fov={40}
>
	<OrbitControls
		bind:ref={orbitControls}
		enableDamping
		dampingFactor={0.05}
		maxDistance={10}
		minDistance={2}
		on:start={handleControlStart}
		on:end={handleControlEnd}
	/>
</T.PerspectiveCamera>
