import ky from 'ky';
import { err, ok, type Result } from 'neverthrow';
import type { Channel, Message, Server, Setup } from '../types/types';
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

const client = ky.create({
	prefixUrl: import.meta.env.VITE_API_URL,
	credentials: 'include',
	retry: 2,
	timeout: 10000
});

class Backend {
	wsConn = $state<WebSocket>();

	setupWebsocket(userId: number) {
		this.wsConn = new WebSocket(`ws://localhost:3000/v1/authenticated/connect/${userId}`);
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
	): Promise<Result<Message, CreateMessageErrors>> {
		try {
			const res = await client.post(`authenticated/messages/${serverId}/${channelId}`, {
				body: JSON.stringify(body)
			});

			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '' });
			}

			const data = (await res.json()) as Message;

			return ok(data);
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
