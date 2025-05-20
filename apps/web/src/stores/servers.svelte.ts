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

	getActiveMembers(serverId: number) {
		return this.servers[serverId]?.active_count?.length || 0;
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
		let messages = this.servers[serverId]?.channels[message.channel_id]?.messages;
		if (Array.isArray(messages)) {
			messages.push(message);
		} else {
			messages = [message];
		}
	}

	editMessage(
		serverId: number,
		channelId: number,
		messageId: number,
		content: any,
		mentions_users: number[],
		mentions_channels: number[],
		updated_at: string
	) {
		const messages = this.servers[serverId]?.channels[channelId]?.messages;

		if (Array.isArray(messages)) {
			const idx = messages.findIndex((m) => m.id === messageId);
			if (idx) {
				messages[idx].content = content;
				messages[idx].mentions_users = mentions_users;
				messages[idx].mentions_channels = mentions_channels;
				messages[idx].updated_at = updated_at;
			}
			return messages[idx];
		}
	}

	deleteMessage(serverId: number, channelId: number, messageId: number) {
		const messages = this.servers[serverId]?.channels[channelId]?.messages;
		if (Array.isArray(messages)) {
			const idx = messages.findIndex((m) => m.id === messageId);
			if (idx) messages.splice(idx, 1);
		}
	}

	connectUser(serverId: number, userId: number, connectedUsers: number[], type: string) {
		const server = this.getServer(serverId);
		if (!server.active_count || server.active_count.length <= 0) {
			this.servers[server.id].active_count = [];
		}

		if (connectedUsers?.length > 0) {
			this.servers[server.id].active_count = connectedUsers;
		}

		if (!this.servers[server.id].active_count.includes(userId)) {
			this.servers[server.id].active_count.push(userId);
		}

		if (type === 'JOIN_SERVER') {
			this.servers[server.id].member_count += 1;
		}
	}

	disconnectUser(serverId: number, userId: number, type: string) {
		const server = this.getServer(serverId);
		if (!server.active_count) return;
		for (let i = 0; i < server.active_count.length; ++i) {
			if (this.servers[server.id].active_count[i] === userId) {
				this.servers[server.id].active_count.splice(i, 1);
			}
		}

		if (type === 'LEAVE_SERVER') {
			this.servers[server.id].member_count -= 1;
		}
	}
}

export const serversStore = new Servers();
