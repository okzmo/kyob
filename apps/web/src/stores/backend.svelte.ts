import ky from 'ky';
import { err, ok, type Result } from 'neverthrow';
import type { Setup } from '../types/types';
import type { SetupErrors } from '../types/errors';

const client = ky.create({
	prefixUrl: import.meta.env.VITE_API_URL,
	credentials: 'include',
	headers: {
		'Content-Type': 'application/json'
	},
	retry: 2,
	timeout: 10000
});

class Backend {
	async getSetup(): Promise<Result<Setup, SetupErrors>> {
		try {
			const res = await client.get('authenticated/setup');

			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '' });
			}
			const data = (await res.json()) as Setup;

			return ok(data);
		} catch (error) {
			const errBody = await error.response.json();
			if (errBody.status === 401) {
				return err({ code: 'ERR_SETUP_UNAUTHORIZED', error: errBody.error });
			}
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}
}

export const backend = new Backend();
