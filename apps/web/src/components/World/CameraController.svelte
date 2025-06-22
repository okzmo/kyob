<script lang="ts">
	import { T, useTask } from '@threlte/core';
	import { OrbitControls } from '@threlte/extras';

	let { avatarPos } = $props();
	let camera = $state<any>();
	let orbitControls = $state<any>();
	let isMoving = $state(false);
	let lastAvatarPos = $state({ x: avatarPos.x, y: avatarPos.y, z: avatarPos.z });

	$effect(() => {
		let moved =
			avatarPos.x !== lastAvatarPos.x ||
			avatarPos.y !== lastAvatarPos.y ||
			avatarPos.z !== lastAvatarPos.z;

		if (moved) {
			isMoving = true;

			lastAvatarPos = {
				x: avatarPos.x,
				y: avatarPos.y,
				z: avatarPos.z
			};
		} else {
			setTimeout(() => {
				isMoving = false;
			}, 100);
		}
	});

	useTask(() => {
		if (!camera || !orbitControls) return;

		orbitControls.target.set(avatarPos.x, avatarPos.y + 1, avatarPos.z);
	});
</script>

<T.PerspectiveCamera
	bind:ref={camera}
	makeDefault
	position={[avatarPos.x, avatarPos.y + 3, avatarPos.z - 5]}
	fov={40}
>
	<OrbitControls
		bind:ref={orbitControls}
		enableDamping
		dampingFactor={0.05}
		maxDistance={10}
		minDistance={2}
	/>
</T.PerspectiveCamera>
