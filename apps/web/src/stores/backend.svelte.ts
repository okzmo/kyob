import ky from 'ky';
import { err, ok, type Result } from 'neverthrow';
import type { ActorMessageTypes, Channel, Message, Server, Setup } from '../types/types';
import { WSMessageSchema } from '../gen/types_pb';
import type {
	CreateChannelErrors,
	CreateMessageErrors,
	CreateServerErrors,
	DeleteChannelErrors,
	DeleteServerErrors,
	SetupErrors,
	StandardError
} from '../types/errors';
import type { CreateChannelType, CreateMessageType, CreateServerType } from '../types/schemas';
import { fromBinary } from '@bufbuild/protobuf';
import { serversStore } from './servers.svelte';
import { timestampDate } from '@bufbuild/protobuf/wkt';

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
			switch (wsMess.type as ActorMessageTypes) {
				case 'channel:message':
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

			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '' });
			}
			const data = (await res.json()) as Setup;

			return ok(data);
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			if (errBody.status === 401) {
				return err({ code: 'ERR_UNAUTHORIZED', error: errBody.error });
			}
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
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

			const res = await client.post('authenticated/server', {
				body: formData
			});

			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '' });
			}
			const data = (await res.json()) as Server;

			return ok(data);
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			if (errBody.status === 400) {
				return err({ code: 'ERR_VALIDATION_FAILED', error: errBody.error });
			}
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async deleteServer(serverId: number): Promise<Result<void, DeleteServerErrors>> {
		try {
			const res = await client.delete(`authenticated/servers/${serverId}`);

			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '' });
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

	async createChannel(
		serverId: number,
		body: CreateChannelType
	): Promise<Result<Channel, CreateChannelErrors>> {
		try {
			const res = await client.post(`authenticated/channels/${serverId}`, {
				body: JSON.stringify(body)
			});

			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '' });
			}

			const data = (await res.json()) as Channel;

			return ok(data);
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

			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '' });
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

			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '' });
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

			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '' });
			}

			const data = (await res.json()) as Message[];

			return ok(data);
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}
}

export const backend = new Backend();
