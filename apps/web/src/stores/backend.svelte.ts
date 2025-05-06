import ky from 'ky';
import { err, ok, type Result } from 'neverthrow';
import type { Server, Setup } from '../types/types';
import type { ServerErrors, SetupErrors, StandardError } from '../types/errors';
import type { CreateServerType } from '../types/schemas';

const client = ky.create({
	prefixUrl: import.meta.env.VITE_API_URL,
	credentials: 'include',
	retry: 2,
	timeout: 10000
});

class Backend {
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
				return err({ code: 'ERR_SETUP_UNAUTHORIZED', error: errBody.error });
			}
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async createServer(body: CreateServerType): Promise<Result<Server, ServerErrors>> {
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
}

export const backend = new Backend();
