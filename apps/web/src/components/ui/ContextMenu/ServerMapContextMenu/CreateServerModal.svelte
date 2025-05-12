<script lang="ts">
	import { Dialog } from 'bits-ui';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { CreateServerSchema } from '../../../../types/schemas';
	import Close from '../../icons/Close.svelte';
	import Cropper from 'svelte-easy-crop';
	import { core } from '../../../../stores/core.svelte';
	import { backend } from '../../../../stores/backend.svelte';
	import { serversStore } from '../../../../stores/servers.svelte';
	import { animateCoordinates } from '../../../../utils/transition';
	import type { Server } from '../../../../types/types';

	let avatar = $state<string | undefined>();
	let crop = $state({ x: 0, y: 0 });
	let zoom = $state(1);
	let minZoom = $state(3);
	let maxZoom = $state(5);

	const { form, errors, enhance } = superForm(defaults(valibot(CreateServerSchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(CreateServerSchema),
		async onUpdate({ form }) {
			if (form.valid) {
				form.data.x = core.openCreateServerModal.x;
				form.data.y = core.openCreateServerModal.y;

				const res = await backend.createServer(form.data);
				if (res.isErr()) {
					if (res.error.code === 'ERR_VALIDATION_FAILED') {
						console.log(res.error.error);
					}
				}

				if (res.isOk()) {
					const server: Server = {
						...res.value,
						is_member: true,
						channels: {}
					};
					serversStore.addServer(server);
					core.openCreateServerModal.status = false;

					const targetX = -(server.x - window.innerWidth / 2 + 32);
					const targetY = -(server.y - window.innerHeight / 2 + 32);

					await animateCoordinates(
						core.offsetServerMap,
						{ x: core.totalOffsetServerMap.x, y: core.totalOffsetServerMap.y },
						{ x: targetX, y: targetY }
					);

					core.totalOffsetServerMap = {
						x: targetX,
						y: targetY
					};
				}
			}
		}
	});

	function onFile(e: Event) {
		const target = e.target as HTMLInputElement;
		const image = target.files?.[0];

		if (image) {
			const dataUrl = URL.createObjectURL(image);
			const img = new Image();

			img.onload = () => {
				const aspectAvatar = 1;
				const aspectImage = img.naturalWidth / img.naturalHeight;

				minZoom = aspectImage > 1 ? aspectImage / aspectAvatar : aspectAvatar / aspectImage;
				zoom = minZoom;

				URL.revokeObjectURL(dataUrl);
			};
			img.src = dataUrl;
			avatar = dataUrl;
			$form.avatar = image;
		}
	}
</script>

<Dialog.Root
	onOpenChange={(s) => {
		if (!s) core.activateMapDragging();
		core.openCreateServerModal.status = s;
	}}
	open={core.openCreateServerModal.status}
>
	<Dialog.Portal>
		<Dialog.Overlay class="fixed inset-0 bg-black/20" />
		<Dialog.Content
			class="bg-main-900 border-main-800 fixed top-1/2 left-1/2 w-[550px] -translate-1/2 rounded-2xl border"
		>
			<form method="post" use:enhance enctype="multipart/form-data">
				<div class="border-b-main-800 relative mb-8 w-full border-b py-7">
					<Dialog.Close
						type="button"
						class="text-main-400 hocus:text-main-50 absolute top-1/2 right-5 -translate-y-1/2 transition-colors hover:cursor-pointer"
					>
						<Close width={18} height={18} />
					</Dialog.Close>
				</div>
				<div class="flex items-center justify-between px-8">
					<div>
						<Dialog.Title class="text-lg font-semibold">Create a new realm</Dialog.Title>
						<Dialog.Description class="text-main-400 max-w-[24rem]  text-sm">
							Realms are what you see on the canvas, communities to share what you love or simply
							interact with people!
						</Dialog.Description>
					</div>
					<div
						class={[
							'relative h-[85px] w-[85px] overflow-hidden rounded-[50%] border-2 text-transparent transition-colors hover:cursor-pointer',
							$errors.avatar
								? 'hocus:bg-red-400/35 border-red-400 bg-red-400/15'
								: 'border-accent-100 bg-accent-100/15 hocus:bg-accent-100/35'
						]}
					>
						<input
							type="file"
							id="avatar"
							name="avatar"
							onchange={onFile}
							aria-label="Realm avatar"
							class="absolute h-full w-full text-transparent hover:cursor-pointer"
						/>
						{#if $form.avatar}
							<Cropper
								image={avatar}
								cropSize={{ height: 85, width: 85 }}
								cropShape="round"
								showGrid={false}
								bind:crop
								bind:zoom
								{minZoom}
								{maxZoom}
								oncropcomplete={(e) => {
									$form.crop = e.pixels;
								}}
							/>
						{/if}
					</div>
				</div>

				<div class="mt-4 flex flex-col px-8">
					<div class="flex items-center gap-x-1">
						<label
							for="realm-name"
							class={['text-sm', $errors.name ? 'text-red-400 ' : 'text-main-500']}
							>Realm name</label
						>
						{#if $errors.name}
							<p class="text-sm text-red-400">- {$errors.name}</p>
						{/if}
					</div>
					<input
						id="realm-name"
						type="text"
						bind:value={$form.name}
						placeholder="My cool community"
						class={[
							'bg-main-800 border-main-600 placeholder:text-main-400 mt-1.5 rounded-xl border py-2.5 focus-visible:ring-0',
							$errors.name ? 'border-red-400' : 'border-main-600'
						]}
					/>
				</div>

				<div class="mt-4 flex flex-col px-8">
					<div class="flex items-center gap-x-1">
						<label
							for="realm-description"
							class={['text-sm', $errors.description ? 'text-red-400 ' : 'text-main-500']}
							>Realm description</label
						>
						{#if $errors.description}
							<p class="text-sm text-red-400">- {$errors.description}</p>
						{/if}
					</div>
					<textarea
						id="realm-description"
						bind:value={$form.description}
						placeholder="Here we do..."
						class={[
							'bg-main-800 border-main-600 placeholder:text-main-400 mt-1.5 min-h-[8rem] rounded-xl border py-2.5 focus-visible:ring-0',
							$errors.description ? 'border-red-400' : 'border-main-600'
						]}
					></textarea>
				</div>

				<div class="border-t-main-800 relative mt-8 w-full border-t py-9">
					<button
						type="submit"
						class="hocus:text-main-50 bg-accent-100/15 text-accent-50 hocus:bg-accent-100/75 absolute top-1/2 right-3 -translate-y-1/2 rounded-lg px-3.5 py-2 transition-colors hover:cursor-pointer"
					>
						Create your realm
					</button>
				</div>
			</form>
		</Dialog.Content>
	</Dialog.Portal>
</Dialog.Root>
