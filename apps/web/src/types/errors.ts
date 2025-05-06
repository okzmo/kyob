export interface StandardError {
	error: string;
	response: Response;
}

export interface SetupErrors {
	code: 'ERR_SETUP_UNAUTHORIZED' | 'ERR_UNKNOWN';
	error: string;
}

export interface ServerErrors {
	code: 'ERR_VALIDATION_FAILED' | 'ERR_UNKNOWN';
	error: string;
}
