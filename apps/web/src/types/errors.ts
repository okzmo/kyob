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
  code: 'ERR_VALIDATION_FAILED' | 'ERR_TOO_MANY_SERVERS' | 'ERR_UNKNOWN';
  error: string;
  cause?: any;
}

export interface DeleteServerErrors {
  code: 'ERR_UNAUTHORIZED' | 'ERR_UNKNOWN';
  error: string;
  cause?: any;
}

export interface LeaveServerErrors {
  code: 'ERR_UNKNOWN';
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
  code: 'ERR_VALIDATION_FAILED' | 'ERR_MESSAGE_TOO_BIG' | 'ERR_UNKNOWN';
  error: string;
  cause?: any;
}

export interface DeleteMessageErrors {
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

export interface GetUserErrors {
  code: 'ERR_USER_NOT_FOUND' | 'ERR_UNKNOWN';
  error: string;
  cause?: any;
}

export interface AddFriendErrors {
  code: 'ERR_ADDING_ITSELF' | 'ERR_USER_NOT_FOUND' | 'ERR_UNKNOWN';
  error: string;
  cause?: any;
}

export interface AcceptFriendErrors {
  code: 'ERR_UNKNOWN';
  error: string;
  cause?: any;
}

export interface DeleteFriendErrors {
  code: 'ERR_UNKNOWN';
  error: string;
  cause?: any;
}

export interface UpdateAccountErrors {
  code: 'ERR_USERNAME_IN_USE' | 'ERR_EMAIL_IN_USE' | 'ERR_UNKNOWN';
  error: string;
  cause?: any;
}

export interface UpdateProfileErrors {
  code: 'ERR_UNKNOWN';
  error: string;
  cause?: any;
}

export interface UpdateAvatarErrors {
  code: 'ERR_UNKNOWN';
  error: string;
  cause?: any;
}

export interface CallErrors {
  code: 'ERR_UNKNOWN';
  error: string;
  cause?: any;
}

export interface EmojiErrors {
  code:
  | 'ERR_MISSING_EMOJIS'
  | 'ERR_SHORTCODES_INVALID'
  | 'ERR_EMOJIS_INVALID'
  | 'ERR_MISSING_SHORTCODES'
  | 'ERR_UNKNOWN';
  error: string;
  cause?: any;
}

export interface DeleteEmojiErrors {
  code: 'ERR_FORBIDDEN' | 'ERR_UNKNOWN';
  error: string;
  cause?: any;
}

export interface RoleErrors {
  code: 'ERR_FORBIDDEN' | 'ERR_UNKNOWN';
  error: string;
  cause?: any;
}
