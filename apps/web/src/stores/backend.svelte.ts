import ky from 'ky';
import { err, ok, type Result } from 'neverthrow';
import type { Channel, ChannelTypes, Message, Server, Setup } from '../types/types';
import { WSMessageSchema } from '../gen/types_pb';
import type {
	CreateChannelErrors,
	CreateInviteErrors,
	CreateMessageErrors,
	CreateServerErrors,
	DeleteChannelErrors,
	DeleteServerErrors,
	JoinServerErrors,
	LeaveServerErrors,
	SetupErrors,
	StandardError
} from '../types/errors';
import type {
	CreateChannelType,
	CreateMessageType,
	CreateServerType,
	JoinServerType
} from '../types/schemas';
import { fromBinary } from '@bufbuild/protobuf';
import { serversStore } from './servers.svelte';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import { windows } from './windows.svelte';

const client = ky.create({
	prefixUrl: import.meta.env.VITE_API_URL,
	credentials: 'include',
	retry: 2,
	timeout: 10000
});

class Backend {
	wsConn = $state<WebSocket>();

	setupWebsocket(userId: number) {
		const ws = new WebSocket(`ws://localhost:3000/v1/authenticated/connect/${userId}`);
		if (!ws) return;

		this.wsConn = ws;
		ws.onopen = () => {
			console.log('Connection established');
			window.setInterval(() => {
				ws.send('heartbeat');
			}, 10 * 1000);
		};

		ws.onmessage = async (event) => {
			if (event.data === 'heartbeat') return;

			const arrayBuffer = await event.data.arrayBuffer();
			const uint8Array = new Uint8Array(arrayBuffer);
			const wsMess = fromBinary(WSMessageSchema, uint8Array, {
				readUnknownFields: false
			});
			switch (wsMess.content.case) {
				case 'chatMessage':
					{
						if (!wsMess.content.value) return;
						const contentStr = new TextDecoder().decode(wsMess.content.value?.content);
						const message: Message = {
							id: wsMess.content.value.id,
							author: {
								id: wsMess.content.value.author!.id,
								username: wsMess.content.value.author!.username,
								display_name: wsMess.content.value.author!.displayName,
								avatar: wsMess.content.value.author!.avatar,
								about: wsMess.content.value.author!.about,
								banner: wsMess.content.value.author!.banner,
								email: wsMess.content.value.author!.email,
								gradient_top: wsMess.content.value.author!.gradientTop,
								gradient_bottom: wsMess.content.value.author!.gradientBottom,
								facts: wsMess.content.value.author!.facts,
								links: wsMess.content.value.author!.links
							},
							server_id: wsMess.content.value.serverId,
							channel_id: wsMess.content.value.channelId,
							content: JSON.parse(contentStr),
							mentions_users: wsMess.content.value.mentionsUsers,
							mentions_channels: wsMess.content.value.mentionsChannels,
							created_at: timestampDate(wsMess.content.value.createdAt!).toISOString()
						};
						serversStore.addMessage(Number(wsMess.content.value?.serverId), message);
					}
					break;
				case 'channelCreation':
					{
						if (!wsMess.content.value) return;
						const value = wsMess.content.value;
						const channel: Channel = {
							id: value.id,
							name: value.name,
							type: value.type as ChannelTypes,
							unread: false,
							x: value.x,
							y: value.y
						};
						serversStore.addChannel(value.serverId, channel);
					}
					break;
				case 'channelRemoved':
					{
						if (!wsMess.content.value) return;
						const value = wsMess.content.value;
						serversStore.removeChannel(value.serverId, value.channelId);
						windows.closeDeadWindow(value.channelId);
					}
					break;
				case 'userConnect':
					{
						if (!wsMess.content.value) return;
						const value = wsMess.content.value;
						serversStore.connectUser(value.serverId, value.userId, value.users, value.type);
					}
					break;
				case 'userDisconnect': {
					if (!wsMess.content.value) return;
					const value = wsMess.content.value;
					serversStore.disconnectUser(value.serverId, value.userId, value.type);
				}
			}
		};

		ws.onclose = (event) => {
			console.log('Connection closed:', event);
		};

		ws.onerror = (error) => {
			console.error('WebSocket error:', error);
		};
	}

	async getSetup(): Promise<Result<Setup, SetupErrors>> {
		try {
			const res = await client.get('authenticated/setup', {
				headers: {
					'Content-Type': 'application/json'
				}
			});

			const data = (await res.json()) as Setup;
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok(data);
		} catch (error) {
			const errBody = await (error as StandardError).response?.json();
			if (errBody?.status === 401) {
				return err({ code: 'ERR_UNAUTHORIZED', error: errBody.error });
			}
			return err({ code: 'ERR_UNKNOWN', error: errBody?.error || '' });
		}
	}

	async createServer(body: CreateServerType): Promise<Result<Server, CreateServerErrors>> {
		try {
			const formData = new FormData();
			formData.append('name', body.name);
			formData.append('description', body.description);
			formData.append('avatar', body.avatar);
			formData.append('crop', JSON.stringify(body.crop));
			formData.append('private', String(body.private));
			formData.append('x', String(body.x));
			formData.append('y', String(body.y));

			const res = await client.post('authenticated/server', {
				body: formData
			});

			const data = (await res.json()) as Server;
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok(data);
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			if (errBody.status === 400) {
				return err({ code: 'ERR_VALIDATION_FAILED', error: errBody.error });
			}
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async joinServer(body: JoinServerType): Promise<Result<Server, JoinServerErrors>> {
		try {
			const res = await client.post('authenticated/server/join', {
				body: JSON.stringify(body)
			});

			const data = (await res.json()) as { server: Server };
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok(data.server);
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			if (errBody.status === 400) {
				return err({ code: errBody.code, error: errBody.error });
			}
			if (errBody.status === 404) {
				return err({ code: 'ERR_INVITE_SERVER_NOT_FOUND', error: errBody.error });
			}
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async deleteServer(serverId: number): Promise<Result<void, DeleteServerErrors>> {
		try {
			const res = await client.delete(`authenticated/servers/${serverId}`);

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			if (errBody.status === 401) {
				return err({ code: 'ERR_UNAUTHORIZED', error: errBody.error });
			}
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async leaveServer(serverId: number): Promise<Result<void, LeaveServerErrors>> {
		try {
			const res = await client.post(`authenticated/server/${serverId}/leave`);

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async createChannel(
		serverId: number,
		body: CreateChannelType
	): Promise<Result<void, CreateChannelErrors>> {
		try {
			const res = await client.post(`authenticated/channels/${serverId}`, {
				body: JSON.stringify(body)
			});

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			if (errBody.status === 401) {
				return err({ code: 'ERR_UNAUTHORIZED', error: errBody.error });
			}
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async deleteChannel(
		serverId: number,
		channelId: number
	): Promise<Result<void, DeleteChannelErrors>> {
		try {
			const res = await client.delete(`authenticated/channels/${serverId}/${channelId}`);

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			if (errBody.status === 401) {
				return err({ code: 'ERR_UNAUTHORIZED', error: errBody.error });
			}
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async sendMessage(
		serverId: number,
		channelId: number,
		body: CreateMessageType
	): Promise<Result<void, CreateMessageErrors>> {
		try {
			const res = await client.post(`authenticated/messages/${serverId}/${channelId}`, {
				body: JSON.stringify(body)
			});

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			if (errBody.status === 400) {
				return err({ code: 'ERR_VALIDATION_FAILED', error: errBody.error });
			}
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async getMessages(channelId: number): Promise<Result<Message[], SetupErrors>> {
		try {
			const res = await client.get(`authenticated/messages/${channelId}`);

			const data = (await res.json()) as Message[];
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok(data);
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async createInvite(serverId: number): Promise<Result<string, CreateInviteErrors>> {
		try {
			const res = await client.get(`authenticated/server/create_invite/${serverId}`);

			const data = (await res.json()) as { invite_link: string };
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok(data.invite_link);
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}
}

export const backend = new Backend();
