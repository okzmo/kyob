import type { Channel, Message, Server } from '../types/types';
import { backend } from './backend.svelte';

class Servers {
	servers = $state<Record<number, Server>>({});

	setupServers(servers: Record<number, Server>) {
		this.servers = servers;
	}

	isOwner(userId: number, serverId: number) {
		return this.servers[serverId].owner_id === userId;
	}

	isMember(serverId: number) {
		return this.servers[serverId]?.is_member;
	}

	getServers() {
		return Object.values(this.servers);
	}

	getServer(id: number) {
		return this.servers[id];
	}

	getChannels(id: number) {
		return Object.values(this.servers[id]?.channels || {});
	}

	getChannel(serverId: number, channelId: number) {
		return this.servers[serverId].channels[channelId];
	}

	async getMessages(serverId: number, channelId: number) {
		const messages = this.servers[serverId]?.channels[channelId]?.messages;

		if (!messages || messages.length <= 0) {
			const res = await backend.getMessages(channelId);
			if (res.isOk()) {
				this.servers[serverId].channels[channelId].messages = res.value;
				return this.servers[serverId].channels[channelId].messages;
			}
		}

		return messages;
	}

	addServer(server: Server) {
		this.servers[server.id] = server;
	}

	addChannel(serverId: number, channel: Channel) {
		if (!this.servers[serverId]) return;
		this.servers[serverId].channels[channel.id] = channel;
	}

	removeServer(serverId: number) {
		delete this.servers[serverId];
	}

	removeChannel(serverId: number, channelId: number) {
		delete this.servers[serverId].channels[channelId];
	}

	addMessage(serverId: number, message: Message) {
		const messages = this.servers[serverId]?.channels[message.channel_id]?.messages;
		if (Array.isArray(messages)) {
			this.servers[serverId].channels[message.channel_id].messages!.push(message);
		} else {
			this.servers[serverId].channels[message.channel_id].messages = [message];
		}
	}
}

export const serversStore = new Servers();
