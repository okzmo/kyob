export interface StandardError {
	error: string;
}

export interface SetupErrors {
	code: 'ERR_SETUP_UNAUTHORIZED' | 'ERR_UNKNOWN';
	error: string;
}
