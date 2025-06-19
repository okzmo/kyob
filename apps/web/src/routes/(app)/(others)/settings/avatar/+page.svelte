<script lang="ts">
	import { Canvas, T } from '@threlte/core';
	import { OrbitControls, GLTF } from '@threlte/extras';
	import { userStore } from 'stores/user.svelte';
</script>

<h1 class="text-2xl font-bold select-none">Avatar</h1>
<p class="text-main-400 mt-3">
	Your avatar is your way to express yourself on Kyob. You can change it at any time.
</p>

<div class="flex max-w-[50rem]">
	<div class="flex h-full flex-grow"></div>
	<div class="relative h-[45rem] w-[20rem]">
		{#if userStore.user}
			<Canvas>
				<T.DirectionalLight position={[10, 10, 5]} intensity={1} />
				<T.AmbientLight intensity={0.4} />

				<GLTF url={userStore.user.body} position={[0, -1.25, 0]} scale={1.5} />

				<T.PerspectiveCamera makeDefault position={[5, 2, 5]} fov={25}>
					<OrbitControls enablePan={false} />
				</T.PerspectiveCamera>
			</Canvas>
		{/if}
	</div>
</div>
