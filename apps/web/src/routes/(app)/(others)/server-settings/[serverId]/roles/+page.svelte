<script lang="ts">
	import { serversStore } from 'stores/servers.svelte';
	import { page } from '$app/state';
	import RoleSidebar from 'components/settings/roles/RoleSidebar.svelte';
	import RoleTabBar from 'components/settings/roles/RoleTabBar.svelte';
	import FormNewRole from 'components/settings/roles/FormNewRole.svelte';
	import type { AbilitiesType, Role } from 'types/types';
	import RoleMembers from 'components/settings/roles/RoleMembers.svelte';

	const server = $derived(serversStore.getServer(page.params.serverId));

	let activeTab = $state('display');
	let activeRole = $state<Role | undefined>();
	let creatingRole = $state(false);

	const PERMISSIONS: { label: string; description: string; ability: AbilitiesType }[] = [
		{
			label: 'Manage Server',
			description: "Allow role to change the server's name and banner.",
			ability: 'MANAGE_SERVER'
		},
		{
			label: 'Manage Channels',
			description: 'Allow role to create, edit or delete channels.',
			ability: 'MANAGE_CHANNELS'
		},
		{
			label: 'Manage Roles',
			description: 'Allow role to create, edit or delete roles.',
			ability: 'MANAGE_ROLES'
		},
		{
			label: 'Manage Expressions',
			description: 'Allow role to create, edit or delete emojis.',
			ability: 'MANAGE_EXPRESSIONS'
		},
		{
			label: 'Manage messages',
			description: 'Allow role to delete other members messages.',
			ability: 'MANAGE_MESSAGES'
		},
		{
			label: 'Ban members',
			description: 'Allow role to ban other members.',
			ability: 'BAN'
		},
		{
			label: 'Kick members',
			description: 'Allow role to kick other members.',
			ability: 'KICK'
		},
		{
			label: 'Mute members',
			description: 'Allow role to mute other members.',
			ability: 'MUTE'
		},
		{
			label: 'Attach files',
			description: 'Allow role to attach files in text channels.',
			ability: 'ATTACH_FILES'
		},
		{
			label: 'Administrator',
			description: 'Allow role to do anything.',
			ability: 'ADMIN'
		}
	];

	$effect(() => {
		if (activeRole) activeTab = 'display';
	});
</script>

<h1 class="text-2xl font-bold select-none">Roles</h1>
<p class="text-main-400 mt-3">
	Use roles to group your members together or/and assign them permissions.
</p>

<hr class="mt-5 w-full border-none" style="height: 1px; background-color: var(--color-main-800);" />

<div class="mt-5 flex h-[calc(100%-8rem)] gap-x-4">
	<RoleSidebar serverId={server.id} bind:activeRole bind:creatingRole roles={server?.roles || []} />
	<div class="border-l-main-800 flex h-full w-full flex-col border-l pl-4">
		{#if creatingRole || activeRole}
			<RoleTabBar bind:activeTab {activeRole} />
		{/if}

		{#if activeRole || creatingRole}
			{#if activeTab === 'display' || activeTab === 'permissions'}
				<FormNewRole
					bind:roles={server.roles}
					bind:creatingRole
					serverId={server.id}
					{activeTab}
					{PERMISSIONS}
					{activeRole}
				/>
			{:else if activeTab === 'members' && activeRole}
				<RoleMembers roleId={activeRole.id} bind:members={activeRole.members} {server} />
			{/if}
		{:else}
			No role selected, either create one or select one.
		{/if}
	</div>
</div>
