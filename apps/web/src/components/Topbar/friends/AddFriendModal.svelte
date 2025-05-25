<script lang="ts">
	import { Dialog } from 'bits-ui';
	import CustomDialogContent from '../../ui/CustomDialogContent/CustomDialogContent.svelte';
	import { defaults, superForm, setError } from 'sveltekit-superforms';
	import { valibot } from 'sveltekit-superforms/adapters';
	import { core } from '../../../stores/core.svelte';
	import Close from '../../ui/icons/Close.svelte';
	import { AddFriendSchema } from '../../../types/schemas';
	import { backend } from '../../../stores/backend.svelte';

	const { form, errors, enhance } = superForm(defaults(valibot(AddFriendSchema)), {
		dataType: 'json',
		SPA: true,
		validators: valibot(AddFriendSchema),
		async onUpdate({ form }) {
			if (form.valid) {
				const res = await backend.addFriend({ username: form.data.username });
				if (res.isErr()) {
					if (res.error.code === 'ERR_USER_NOT_FOUND') {
						setError(form, 'username', 'User not found.');
					}
					if (res.error.code === 'ERR_UNKNOWN') {
						console.log(res.error.error);
					}
				}

				if (res.isOk()) {
					core.openAddFriendModal.status = false;
					core.activateMapDragging();
				}
			}
		}
	});
</script>

<Dialog.Root
	onOpenChange={(s) => {
		if (!s) core.activateMapDragging();
		core.openAddFriendModal.status = s;
	}}
	open={core.openAddFriendModal.status}
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
					<Dialog.Title class="text-lg font-semibold">Add a new friend</Dialog.Title>
					<Dialog.Description class="text-main-400 max-w-[24rem]  text-sm">
						Who's that friend mmh :)
					</Dialog.Description>
				</div>
			</div>

			<form method="post" use:enhance>
				<div class="mt-4 flex flex-col px-8">
					<div class="flex items-center gap-x-1">
						<label
							for="channel-name"
							class={['text-sm', $errors.username ? 'text-red-400 ' : 'text-main-500']}
							>Friend username</label
						>
						{#if $errors.username}
							<p class="text-sm text-red-400">- {$errors.username}</p>
						{/if}
					</div>
					<input
						id="channel-name"
						type="text"
						bind:value={$form.username}
						placeholder="bob"
						class={[
							'bg-main-800 border-main-600 placeholder:text-main-400 mt-1.5 rounded-xl border py-2.5 focus-visible:ring-0',
							$errors.username ? 'border-red-400' : 'border-main-600'
						]}
					/>
				</div>

				<div class="border-t-main-800 relative mt-8 w-full border-t py-9">
					<button
						type="submit"
						class="hocus:text-main-50 bg-accent-100/15 text-accent-50 hocus:bg-accent-100/75 absolute top-1/2 right-3 -translate-y-1/2 rounded-lg px-3.5 py-2 transition-colors hover:cursor-pointer"
					>
						Add friend
					</button>
				</div>
			</form>
		</CustomDialogContent>
	</Dialog.Portal>
</Dialog.Root>
