<script lang="ts">
	import Check from 'components/ui/icons/Check.svelte';
	import { backend } from 'stores/backend.svelte';
	import type { Server } from 'types/types';

	interface Props {
		roleId: string;
		server: Server;
		members: string[];
	}

	let { roleId, members = $bindable(), server }: Props = $props();

	async function addMember(memberId: string) {
		const res = await backend.addRoleMember(server.id, roleId, memberId);

		if (res.isErr()) console.error(res.error);
	}

	async function removeMember(memberId: string) {
		const res = await backend.removeRoleMember(server.id, roleId, memberId);

		if (res.isErr()) console.error(res.error);
	}
</script>

<ul class="mt-4 flex flex-col">
	{#each server.members as member (member.id)}
		<li
			class="border-b-main-800 flex items-center justify-between border-b px-2 py-4 first:pt-2 last:border-b-transparent"
		>
			<div class="flex items-center gap-x-2">
				<img src={member.avatar} alt="" class="h-8 w-8" />
				{member.display_name}
			</div>

			<label
				for={`member-${member.id}`}
				class="group bg-main-800 border-main-700 has-checked:bg-accent-100/40 has-checked:border-accent-100 flex h-[1.25rem] w-[1.25rem] items-center justify-center border transition-colors hover:cursor-pointer"
			>
				<input
					id={`member-${member.id}`}
					type="checkbox"
					class="absolute h-0 w-0 opacity-0"
					checked={member.roles?.includes(roleId)}
					onchange={() => {
						if (member.roles?.includes(roleId)) {
							removeMember(member.id!);
						} else {
							addMember(member.id!);
						}
					}}
				/>
				<Check
					height={16}
					width={16}
					class="text-main-800 group-has-checked:text-accent-100 transition-colors"
				/>
			</label>
		</li>
	{/each}
</ul>
