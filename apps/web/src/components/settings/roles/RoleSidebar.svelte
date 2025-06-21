<script lang="ts">
	import DragHandle from 'components/ui/icons/DragHandle.svelte';
	import PlusSimple from 'components/ui/icons/PlusSimple.svelte';
	import { backend } from 'stores/backend.svelte';
	import type { Role } from 'types/types';

	interface Props {
		activeRole: Role | undefined;
		creatingRole: boolean;
		roles: Role[];
		serverId: string;
	}

	let { serverId, activeRole = $bindable(), creatingRole = $bindable(), roles }: Props = $props();

	let draggedIndex = $state<number | null>(null);
	let draggedOverIndex = $state<number | null>(null);

	function handleDragStart(e: DragEvent, idx: number) {
		draggedIndex = idx;
		if (e.dataTransfer) e.dataTransfer.effectAllowed = 'move';
	}

	function handleDragOver(e: DragEvent, idx: number) {
		e.preventDefault();
		draggedOverIndex = idx;
	}

	function handleDrop(e: DragEvent, idx: number) {
		e.preventDefault();

		if (draggedIndex === null || draggedOverIndex === null) return;

		const draggedItem = roles[draggedIndex];
		const newRoles = [...roles];

		newRoles.splice(draggedIndex, 1);
		newRoles.splice(idx, 0, draggedItem);

		roles = newRoles;

		backend.moveRole(serverId, draggedItem.id, draggedIndex, idx);

		draggedIndex = null;
		draggedOverIndex = null;
	}

	function handleDragEnd() {
		draggedIndex = null;
		draggedOverIndex = null;
	}

	function hoverRole(target: HTMLElement, color: string) {
		target.style.backgroundColor = `rgba(${color}, 0.25)`;
		target.style.borderColor = `rgba(${color}, 0.85)`;
		target.style.color = `rgba(${color}, 1)`;
	}

	function unhoverRole(target: HTMLElement, color: string, roleId: string) {
		if (activeRole?.id === roleId) return;
		target.style.backgroundColor = `rgba(${color}, 0.15)`;
		target.style.borderColor = `rgba(${color}, 0.4)`;
		target.style.color = `rgba(${color}, 0.75)`;
	}
</script>

<ul class="flex w-[15rem] flex-col gap-y-2">
	{#each roles as role, idx (role.id)}
		<li>
			<button
				draggable="true"
				class="justify-left flex w-full items-center gap-x-2 border py-2 pr-4 pl-2 font-semibold transition-all duration-100 hover:cursor-pointer {draggedIndex ===
				idx
					? 'opacity-50'
					: ''} {draggedOverIndex === idx ? 'border-t-4 border-t-blue-500' : ''}"
				style={activeRole?.id === role.id
					? `color: rgba(${role.color}, 1); background-color: rgba(${role.color}, 0.25); border-color: rgba(${role.color}, 0.85);`
					: `color: rgba(${role.color}, 0.75); background-color: rgba(${role.color}, 0.15); border-color: rgba(${role.color}, 0.4);`}
				onclick={() => (activeRole = role)}
				onmouseover={(e) => hoverRole(e.currentTarget, role.color)}
				onfocus={(e) => hoverRole(e.currentTarget, role.color)}
				onmouseleave={(e) => unhoverRole(e.currentTarget, role.color, role.id)}
				onblur={(e) => unhoverRole(e.currentTarget, role.color, role.id)}
				ondragstart={(e) => handleDragStart(e, idx)}
				ondragover={(e) => handleDragOver(e, idx)}
				ondrop={(e) => handleDrop(e, idx)}
				ondragend={handleDragEnd}
			>
				<span class="drag-handle hover:cursor-grab active:cursor-grabbing">
					<DragHandle height={20} width={20} />
				</span>
				{role.name}
			</button>
		</li>
	{/each}
	<li>
		<button
			class="border-main-400 bg-main-400/15 text-main-400 hocus:bg-main-400/25 flex w-full items-center justify-center gap-x-2 border px-2 py-2 font-semibold transition-colors duration-100 hover:cursor-pointer"
			onclick={() => {
				creatingRole = true;
				activeRole = undefined;
			}}
		>
			<PlusSimple height={20} width={20} />
		</button>
	</li>
</ul>
