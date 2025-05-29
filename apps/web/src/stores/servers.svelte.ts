import type { Channel, Message, Server } from '../types/types';
import { backend } from './backend.svelte';

class Servers {
	servers = $state<Record<string, Server>>({});

	setupServers(servers: Record<string, Server>) {
		this.servers = servers;
	}

	isOwner(userId: string, serverId: string) {
		return this.servers[serverId].owner_id === userId;
	}

	getServers() {
		return Object.values(this.servers);
	}

	getServer(id: string) {
		return this.servers?.[id];
	}

	getChannels(id: string) {
		return Object.values(this.servers[id]?.channels || {});
	}

	getChannel(serverId: string, channelId: string) {
		return this.servers[serverId].channels[channelId];
	}

	getActiveMembers(serverId: string) {
		return this.servers[serverId]?.active_count?.length || 0;
	}

	getMemberById(serverId: string, userId: string) {
		return this.servers[serverId].members.find((m) => m.id === userId);
	}

	async getMessages(serverId: string, channelId: string) {
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

	addChannel(serverId: string, channel: Channel) {
		if (!this.servers[serverId]) return;
		this.servers[serverId].channels[channel.id] = channel;
	}

	hideServer(serverId: string) {
		this.servers[serverId].hidden = true;
	}

	showServer(serverId: string) {
		this.servers[serverId].hidden = false;
	}

	removeServer(serverId: string) {
		delete this.servers[serverId];
	}

	removeChannel(serverId: string, channelId: string) {
		delete this.servers[serverId].channels[channelId];
	}

	addMessage(serverId: string, message: Message) {
		const messages = this.servers[serverId]?.channels[message.channel_id]?.messages;
		if (Array.isArray(messages)) {
			messages.push(message);
		} else {
			this.servers[serverId].channels[message.channel_id].messages = [message];
		}
	}

	editMessage(
		serverId: string,
		channelId: string,
		messageId: string,
		content: any,
		mentions_users: string[],
		mentions_channels: string[],
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

	getMessage(serverId: string, channelId: string, messageId: string) {
		const messages = this.servers[serverId]?.channels[channelId]?.messages;

		if (Array.isArray(messages)) {
			const idx = messages.findIndex((m) => m.id === messageId);
			return messages[idx];
		}
	}

	deleteMessage(serverId: string, channelId: string, messageId: string) {
		const messages = this.servers[serverId]?.channels[channelId]?.messages;
		if (Array.isArray(messages)) {
			const idx = messages.findIndex((m) => m.id === messageId);
			messages.splice(idx, 1);
		}
	}

	connectUser(serverId: string, userId: string, connectedUsers: string[], type: string) {
		if (serverId === 'global') return;

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

	disconnectUser(serverId: string, userId: string, type: string) {
		if (serverId === 'global') return;

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
