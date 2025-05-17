<script lang="ts">
	import { Dialog } from 'bits-ui';
	import { defaults, superForm } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { CreateChannelSchema } from '../../../../types/schemas';
	import Close from '../../icons/Close.svelte';
	import { core } from '../../../../stores/core.svelte';
	import { backend } from '../../../../stores/backend.svelte';
	import { page } from '$app/state';
	import CustomDialogContent from '../../CustomDialogContent/CustomDialogContent.svelte';

	const { form, errors, enhance } = superForm(defaults(valibot(CreateChannelSchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(CreateChannelSchema),
		async onUpdate({ form }) {
			if (form.valid) {
				const serverId = Number(page.params.server_id);
				form.data.x = core.openCreateChannelModal.x;
				form.data.y = core.openCreateChannelModal.y;

				const res = await backend.createChannel(serverId, form.data);
				if (res.isErr()) {
					if (res.error.code === 'ERR_VALIDATION_FAILED') {
						console.log(res.error.error);
					}
					if (res.error.code === 'ERR_UNAUTHORIZED') {
						console.log(res.error.error);
					}
				}

				if (res.isOk()) {
					core.openCreateChannelModal.status = false;
					core.activateMapDragging();
				}
			}
		}
	});
</script>

<Dialog.Root
	onOpenChange={(s) => {
		if (!s) core.activateMapDragging();
		core.openCreateChannelModal.status = s;
	}}
	open={core.openCreateChannelModal.status}
>
	<Dialog.Portal>
		<Dialog.Overlay class="fixed inset-0 bg-black/20" />
		<CustomDialogContent
			class="bg-main-900 border-main-800 fixed top-1/2 left-1/2 w-[550px] -translate-1/2 rounded-2xl border"
		>
			<div class="border-b-main-800 relative mb-8 w-full border-b py-7">
				<Dialog.Close
					class="text-main-400 hocus:text-main-50 absolute top-1/2 right-5 -translate-y-1/2 transition-colors hover:cursor-pointer"
				>
					<Close width={18} height={18} />
				</Dialog.Close>
			</div>
			<div class="flex items-center justify-between px-8">
				<div>
					<Dialog.Title class="text-lg font-semibold">Create a new channel</Dialog.Title>
					<Dialog.Description class="text-main-400 max-w-[24rem]  text-sm">
						Create channels to interact with peoples based on topics. What you like, dislike, memes,
						you choose.
					</Dialog.Description>
				</div>
			</div>

			<form method="post" use:enhance>
				<div class="mt-4 flex flex-col px-8">
					<div class="flex items-center gap-x-1">
						<label
							for="channel-name"
							class={['text-sm', $errors.name ? 'text-red-400 ' : 'text-main-500']}
							>Channel name</label
						>
						{#if $errors.name}
							<p class="text-sm text-red-400">- {$errors.name}</p>
						{/if}
					</div>
					<input
						id="channel-name"
						type="text"
						bind:value={$form.name}
						placeholder="General"
						class={[
							'bg-main-800 border-main-600 placeholder:text-main-400 mt-1.5 rounded-xl border py-2.5 focus-visible:ring-0',
							$errors.name ? 'border-red-400' : 'border-main-600'
						]}
					/>
				</div>

				<div class="border-t-main-800 relative mt-8 w-full border-t py-9">
					<button
						type="submit"
						class="hocus:text-main-50 bg-accent-100/15 text-accent-50 hocus:bg-accent-100/75 absolute top-1/2 right-3 -translate-y-1/2 rounded-lg px-3.5 py-2 transition-colors hover:cursor-pointer"
					>
						Create channel
					</button>
				</div>
			</form>
		</CustomDialogContent>
	</Dialog.Portal>
</Dialog.Root>
