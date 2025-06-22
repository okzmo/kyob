<script lang="ts">
	import { Canvas, T } from '@threlte/core';
	import { Sky } from '@threlte/extras';
	import { onMount } from 'svelte';
	import CameraController from './CameraController.svelte';
	import Character from './Character.svelte';

	let avatarPos = { x: 0, y: 0, z: 0 };
</script>

<Canvas shadows>
	<CameraController {avatarPos} />

	<T.AmbientLight intensity={0.4} />
	<T.DirectionalLight
		position={[10, 10, 5]}
		intensity={1}
		castShadow
		shadow.camera.left={-10}
		shadow.camera.right={10}
		shadow.camera.top={10}
		shadow.camera.bottom={-10}
	/>

	<Sky elevation={0.5} azimuth={180} />

	<T.Mesh receiveShadow rotation={[-Math.PI / 2, 0, 0]}>
		<T.PlaneGeometry args={[20, 20]} />
		<T.MeshStandardMaterial color="#90EE90" />
	</T.Mesh>

	<Character bind:avatarPos />
</Canvas>
