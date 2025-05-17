export interface StandardError {
	error: string;
	cause?: any;
	response: Response;
}

export interface SetupErrors {
	code: 'ERR_UNAUTHORIZED' | 'ERR_UNKNOWN';
	error: string;
	cause?: any;
}

export interface MessagesErrors {
	code: 'ERR_UNKNOWN';
	error: string;
	cause?: any;
}

export interface CreateServerErrors {
	code: 'ERR_VALIDATION_FAILED' | 'ERR_UNKNOWN';
	error: string;
	cause?: any;
}

export interface DeleteServerErrors {
	code: 'ERR_UNAUTHORIZED' | 'ERR_UNKNOWN';
	error: string;
	cause?: any;
}

export interface CreateChannelErrors {
	code: 'ERR_VALIDATION_FAILED' | 'ERR_UNAUTHORIZED' | 'ERR_UNKNOWN';
	error: string;
	cause?: any;
}

export interface DeleteChannelErrors {
	code: 'ERR_UNAUTHORIZED' | 'ERR_UNKNOWN';
	error: string;
	cause?: any;
}

export interface CreateMessageErrors {
	code: 'ERR_VALIDATION_FAILED' | 'ERR_UNKNOWN';
	error: string;
	cause?: any;
}

export interface CreateInviteErrors {
	code: 'ERR_UNKNOWN';
	error: string;
	cause?: any;
}

export interface JoinServerErrors {
	code:
		| 'ERR_INVITE_MISSING_ID'
		| 'ERR_INVITE_SERVER_NOT_FOUND'
		| 'ERR_VALIDATION_FAILED'
		| 'ERR_UNKNOWN';
	error: string;
	cause?: any;
}
