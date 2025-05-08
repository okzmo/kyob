import type { Channel, Server } from '../types/types';

class Servers {
	servers = $state<Server[]>([]);

	getServer(id: number) {
		return this.servers.find((s) => s.id === id);
	}

	getChannels(id: number) {
		return this.servers.find((s) => s.id === id)?.channels;
	}

	getChannel(serverId: number, channelId: number) {
		const channels = this.getChannels(serverId);
		return channels?.find((c) => c.id === channelId);
	}

	isOwner(userId: number, serverId: number) {
		return Boolean(this.servers.find((s) => s.id === serverId && s.owner_id === userId));
	}

	isMember(serverId: number) {
		return Boolean(this.servers.find((s) => s.id === serverId && s.is_member));
	}

	addChannel(serverId: number, channel: Channel) {
		const serverIdx = this.servers.findIndex((s) => s.id === serverId);
		if (Array.isArray(this.servers[serverIdx].channels)) {
			this.servers[serverIdx].channels.push(channel);
		} else {
			this.servers[serverIdx].channels = [channel];
		}
	}

	removeServer(serverId: number) {
		for (let i = 0; i < this.servers.length; ++i) {
			if (this.servers[i].id === serverId) {
				this.servers.splice(i, 1);
			}
		}
	}

	removeChannel(serverId: number, channelId: number) {
		const serverIdx = this.servers.findIndex((s) => s.id === serverId);
		if (!this.servers[serverIdx].channels) return;

		for (let i = 0; i < this.servers[serverIdx].channels.length; ++i) {
			if (this.servers[serverIdx].channels[i].id === channelId) {
				this.servers[serverIdx].channels.splice(i, 1);
			}
		}
	}
}

export const serversStore = new Servers();
