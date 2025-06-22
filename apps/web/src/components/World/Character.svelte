<script lang="ts">
	import { useTask } from '@threlte/core';
	import { GLTF } from '@threlte/extras';
	import { userStore } from 'stores/user.svelte';
	import { onMount } from 'svelte';

	let { avatarPos = $bindable() } = $props();
	let keys = $state({ z: false, s: false, q: false, d: false });
	let speed = $state(0.05);

	useTask(() => {
		if (keys.z) avatarPos.z += speed;
		if (keys.s) avatarPos.z -= speed;
		if (keys.q) avatarPos.x += speed;
		if (keys.d) avatarPos.x -= speed;
	});

	function handleKeyDown(e: KeyboardEvent) {
		if (e.key in keys) keys[e.key] = true;
	}

	function handleKeyUp(e: KeyboardEvent) {
		if (e.key in keys) keys[e.key] = false;
	}

	onMount(() => {
		window.addEventListener('keydown', handleKeyDown);
		window.addEventListener('keyup', handleKeyUp);

		return () => {
			window.removeEventListener('keydown', handleKeyDown);
			window.removeEventListener('keyup', handleKeyUp);
		};
	});
</script>

<GLTF
	url={`https://models.readyplayer.me/${userStore.user!.rpm_avatar_id}.glb`}
	castShadow
	position={[avatarPos.x, avatarPos.y, avatarPos.z]}
	scale={1}
/>
