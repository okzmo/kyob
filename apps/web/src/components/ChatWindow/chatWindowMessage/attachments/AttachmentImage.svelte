<script lang="ts">
	import { core } from 'stores/core.svelte';
	import { windows } from 'stores/windows.svelte';

	let { image, images, idx, medias } = $props();
</script>

<button
	class={[
		'attachment pointer-events-auto relative select-none',
		medias.length > 1 && 'aspect-square'
	]}
	onclick={() => {
		core.openAttachmentsModal(images, idx);
		windows.setActiveWindow(null);
	}}
>
	<img src={image.url} alt="Attachment" class="h-full w-full object-cover" />
</button>
