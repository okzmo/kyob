<script lang="ts">
	import { Dialog } from 'bits-ui';
	import Corners from 'components/ui/Corners/Corners.svelte';
	import CustomDialogContent from 'components/ui/CustomDialogContent/CustomDialogContent.svelte';
	import FooterDialog from 'components/ui/CustomDialogContent/FooterDialog.svelte';
	import FormInput from 'components/ui/FormInput/FormInput.svelte';
	import SubmitButton from 'components/ui/SubmitButton/SubmitButton.svelte';
	import { backend } from 'stores/backend.svelte';
	import { core } from 'stores/core.svelte';
	import { serversStore } from 'stores/servers.svelte';
	import { userStore } from 'stores/user.svelte';
	import Cropper from 'svelte-easy-crop';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { CreateServerSchema } from 'types/schemas';
	import type { Server } from 'types/types';
	import { delay } from 'utils/delay';
	import { animateCoordinates } from 'utils/transition';

	let avatar = $state<string | undefined>();
	let crop = $state({ x: 0, y: 0 });
	let zoom = $state(1);
	let minZoom = $state(3);
	let maxZoom = $state(5);

	let isSubmitting = $state(false);
	let isSubmitted = $state(false);
	let buttonWidth = $derived(isSubmitted || isSubmitting ? 40 : 154);

	const { form, errors, enhance } = superForm(defaults(valibot(CreateServerSchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(CreateServerSchema),
		async onUpdate({ form }) {
			if (form.valid) {
				form.data.x = Math.round(core.createServerModal.x - core.totalOffsetServerMap.x - 32);
				form.data.y = Math.round(core.createServerModal.y - core.totalOffsetServerMap.y - 32);

				isSubmitting = true;
				const res = await backend.createServer(form.data);
				if (res.isErr()) {
					if (res.error.code === 'ERR_VALIDATION_FAILED') {
						console.error(res.error.error);
					}
					isSubmitting = false;
				}

				if (res.isOk()) {
					const server: Server = {
						...res.value,
						channels: {},
						member_count: 1,
						members: [
							{
								id: userStore.user?.id,
								username: userStore.user?.username,
								display_name: userStore.user?.display_name,
								avatar: userStore.user?.avatar
							}
						],
						active_count: [],
						hidden: false
					};

					serversStore.addServer(server);

					await delay(400);
					isSubmitting = false;
					isSubmitted = true;
					await delay(800);

					core.createServerModal.status = false;
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

					setTimeout(async () => {
						core.activateMapDragging();
					}, 500);

					isSubmitted = false;
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

	let isEmpty = $derived(!$form.name && !$form.avatar);
</script>

<Dialog.Root
	onOpenChange={(s) => {
		if (!s) core.activateMapDragging();
		core.createServerModal.status = s;
	}}
	open={core.createServerModal.status}
>
	<Dialog.Portal>
		<Dialog.Overlay class="fixed inset-0 bg-black/20" />
		<CustomDialogContent>
			<form method="post" use:enhance enctype="multipart/form-data">
				<div class="flex items-center justify-between px-8">
					<div>
						<Dialog.Title class="text-lg font-semibold">Create a new realm</Dialog.Title>
						<Dialog.Description class="text-main-400 max-w-[24rem] text-sm">
							Realms are what you see on the canvas, communities to share what you love or simply
							interact with people!
						</Dialog.Description>
					</div>
					<div
						class={[
							'group relative h-[85px] w-[85px] overflow-hidden text-transparent transition-colors hover:cursor-pointer',
							$errors.avatar
								? 'hocus:bg-red-400/25 inner-red-400/20 hocus:inner-red-400/40 bg-red-400/15'
								: 'inner-accent/15 bg-accent-100/15 hocus:bg-accent-100/35 hocus:inner-accent-no-shadow/25'
						]}
					>
						<Corners
							color={$errors.avatar ? 'border-red-400/50' : 'border-accent-100/50'}
							class={$errors.avatar
								? 'group-hocus:border-red-400'
								: 'group-hocus:border-accent-100'}
						/>
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
								cropShape="rect"
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

				<FormInput
					title="Realm name"
					id="realm-name"
					type="text"
					bind:error={$errors.name}
					bind:inputValue={$form.name}
					placeholder="My cool community"
					class="mt-4 px-8"
				/>

				<FormInput
					title="Realm description"
					id="realm-description"
					type="rich"
					bind:error={$errors.description}
					bind:inputValue={$form.description}
					placeholder="Here we do..."
					class="mt-4 px-8"
					inputClass="w-full"
				/>

				<FooterDialog>
					<SubmitButton
						type="submit"
						{buttonWidth}
						{isEmpty}
						{isSubmitting}
						{isSubmitted}
						class="relative"
					>
						Create your realm
					</SubmitButton>
				</FooterDialog>
			</form>
		</CustomDialogContent>
	</Dialog.Portal>
</Dialog.Root>
