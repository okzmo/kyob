export interface StandardError {
	error: string;
	response: Response;
}

export interface SetupErrors {
	code: 'ERR_UNAUTHORIZED' | 'ERR_UNKNOWN';
	error: string;
}

export interface CreateServerErrors {
	code: 'ERR_VALIDATION_FAILED' | 'ERR_UNKNOWN';
	error: string;
}

export interface DeleteServerErrors {
	code: 'ERR_UNAUTHORIZED' | 'ERR_UNKNOWN';
	error: string;
}

export interface CreateChannelErrors {
	code: 'ERR_VALIDATION_FAILED' | 'ERR_UNAUTHORIZED' | 'ERR_UNKNOWN';
	error: string;
}

export interface DeleteChannelErrors {
	code: 'ERR_UNAUTHORIZED' | 'ERR_UNKNOWN';
	error: string;
}
