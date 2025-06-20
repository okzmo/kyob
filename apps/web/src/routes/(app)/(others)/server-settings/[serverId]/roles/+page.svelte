<script lang="ts">
	import { serversStore } from 'stores/servers.svelte';
	import { page } from '$app/state';
	import RoleSidebar from 'components/settings/roles/RoleSidebar.svelte';
	import RoleTabBar from 'components/settings/roles/RoleTabBar.svelte';
	import FormEditRole from 'components/settings/roles/FormEditRole.svelte';
	import FormNewRole from 'components/settings/roles/FormNewRole.svelte';

	const server = $derived(serversStore.getServer(page.params.serverId));

	let activeTab = $state('display');
	let activeRole = $state<number | undefined>();
	let creatingRole = $state(false);
</script>

<h1 class="text-2xl font-bold select-none">Roles</h1>
<p class="text-main-400 mt-3">
	Use roles to group your members together or/and assign them permissions.
</p>

<hr class="mt-5 w-full border-none" style="height: 1px; background-color: var(--color-main-800);" />

<div class="mt-5 flex h-[calc(100%-8rem)]">
	<RoleSidebar bind:activeRole bind:creatingRole />
	<div class="bg-main-800 mx-4 h-full w-[1px]"></div>
	<div class="flex flex-col">
		{#if creatingRole || activeRole}
			<RoleTabBar bind:activeTab />
		{/if}

		{#if activeRole}
			<FormEditRole {activeTab} />
		{:else if creatingRole}
			<FormNewRole {activeTab} />
		{:else}
			No role selected, either create one or select one.
		{/if}
	</div>
</div>
