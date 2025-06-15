import type { Channel, LastState, Message, Server, User } from '../types/types';
import { backend } from './backend.svelte';
import { userStore } from './user.svelte';
import { windows } from './windows.svelte';

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
		return this.servers?.[serverId]?.members?.find((m) => m.id === userId);
	}

	async getMessages(serverId: string, channelId: string) {
		const messages = this.servers[serverId]?.channels[channelId]?.messages;

		if (!messages) {
			const res = await backend.getMessages(channelId);
			if (res.isOk()) {
				this.servers[serverId].channels[channelId].messages = res.value || [];
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

	removeServer(serverId: string) {
		delete this.servers[serverId];
	}

	removeChannel(serverId: string, channelId: string) {
		delete this.servers[serverId].channels[channelId];
	}

	markChannelAsRead(serverId: string, channelId: string) {
		const channel = this.getChannel(serverId, channelId);

		if (Array.isArray(channel.messages) && channel.messages.length > 0) {
			channel.last_message_read = channel.messages[0].id;
			channel.last_mentions = [];
		}
	}

	addMessage(serverId: string, message: Message) {
		const messages = this.servers[serverId]?.channels[message.channel_id]?.messages;
		const channel = this.getChannel(serverId, message.channel_id);

		if (Array.isArray(messages)) {
			messages.unshift(message);
		}

		if (!windows.getWindow({ channelId: message.channel_id })) {
			channel.last_message_sent = message.id;

			if (message.mentions_users.find((id) => id === userStore.user!.id) || message.everyone) {
				if (Array.isArray(channel.last_mentions))
					channel.last_mentions = [...channel.last_mentions, message.id];
				else channel.last_mentions = [message.id];
			}
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
			if (idx > -1) {
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

	addMember(serverId: string, user: Partial<User>) {
		if (serverId === 'global') return;

		const server = this.getServer(serverId);
		server.members.push(user);
	}

	modifyMember(serverId: string, userId: string, user: Partial<User>) {
		const member = this.servers[serverId].members.find((m) => m.id === userId);
		if (!member) return;

		if (user.avatar) member.avatar = user.avatar;
		if (user.display_name) member.display_name = user.display_name;
		if (user.username) member.username = user.username;
	}

	disconnectUser(serverId: string, userId: string, type: string) {
		if (serverId === 'global') return;

		const server = this.getServer(serverId);
		if (!server.active_count) return;
		for (let i = 0; i < server.active_count.length; ++i) {
			if (server.active_count[i] === userId) {
				server.active_count.splice(i, 1);
			}
		}

		if (type === 'LEAVE_SERVER') {
			const memberIdx = server.members.findIndex((m) => m.id === userId);

			server.members.splice(memberIdx, 1);
			server.member_count -= 1;
		}
	}

	isInCall(serverId: string, channelId: string, userId: string) {
		const channel = this.getChannel(serverId, channelId);
		return channel.voice_users.find((u) => u.user_id === userId);
	}

	initiateCall(serverId: string, channelId: string, users: any) {
		const channel = this.getChannel(serverId, channelId);
		channel.voice_users = users;
	}

	connectUserToCall(serverId: string, channelId: string, userId: string) {
		const channel = this.getChannel(serverId, channelId);
		channel.voice_users.push({ user_id: userId, deafen: false, mute: false });
	}

	disconnectUserFromCall(serverId: string, channelId: string, userId: string) {
		const channel = this.getChannel(serverId, channelId);
		const userIdx = channel.voice_users.findIndex((u) => u.user_id === userId);
		if (userIdx > -1) {
			channel.voice_users.splice(userIdx, 1);
		}
	}

	hasUnreadChannels(serverId: string) {
		const channels = this.getChannels(serverId);

		for (const channel of channels) {
			if (
				channel.last_message_read &&
				channel.last_message_sent &&
				channel.last_message_read < channel.last_message_sent
			) {
				return true;
			}
		}

		return false;
	}

	hasMentionsInChannels(serverId: string) {
		const channels = this.getChannels(serverId);

		for (const channel of channels) {
			if (channel.last_mentions && channel.last_mentions.length > 0) {
				return channel.last_mentions.length;
			}
		}

		return false;
	}

	getLastState(): LastState {
		const lastState: LastState = {
			channel_ids: [],
			last_message_ids: [],
			mentions_ids: []
		};

		for (const server of Object.values(this.servers)) {
			for (const channel of Object.values(server.channels)) {
				lastState.channel_ids.push(channel.id);
				lastState.last_message_ids.push(channel.last_message_read || '');
				lastState.mentions_ids.push(channel.last_mentions || []);
			}
		}

		return lastState;
	}
}

export const serversStore = new Servers();
