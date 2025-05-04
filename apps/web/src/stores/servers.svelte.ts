import type { Server } from '../types/types';

class Servers {
	servers = $state<Server[]>([]);

	getServer(id: number) {
		return this.servers.find((s) => s.id === id);
	}

	getChannels(id: number) {
		return this.servers.find((s) => s.id === id)?.channels;
	}

	isOwner(userId: number, serverId: number) {
		return Boolean(this.servers.find((s) => s.id === serverId && s.owner_id === userId));
	}

	isMember(serverId: number) {
		return Boolean(this.servers.find((s) => s.id === serverId && s.is_member));
	}
}

export const serversStore = new Servers();
