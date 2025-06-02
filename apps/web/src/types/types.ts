export const ChannelTypes = {
	Textual: 'textual',
	Voice: 'voice'
} as const;
export type ChannelTypes = (typeof ChannelTypes)[keyof typeof ChannelTypes];

export const contextMenuTargets = [
	'serverButton',
	'channelButton',
	'message',
	'inServer',
	'mainMap'
] as const;
export type ContextMenuTarget = (typeof contextMenuTargets)[number];

export interface Window {
	id: string;
	channelId?: string;
	serverId?: string;
	friendId?: string;
}

export interface Channel {
	id: string;
	name: string;
	type: ChannelTypes;
	x: number;
	y: number;
	unread: boolean;
	messages?: Message[];
	users?: Partial<User>[];
}

export interface Server {
	id: string;
	owner_id: string;
	name: string;
	avatar: string;
	banner: string;
	description?: string;
	x: number;
	y: number;
	channels: Record<string, Channel>;
	active_count: string[];
	member_count: number;
	members: Partial<User>[];
	hidden: boolean;
}

export interface Fact {
	id: string;
	label: string;
	value: string;
}

export interface Link {
	id: string;
	label: string;
	url: string;
}

export interface User {
	id: string;
	email: string;
	username: string;
	display_name: string;
	avatar: string;
	banner: string;
	main_color?: string;
	about?: string;
	facts?: Fact[];
	links?: Link[];
}

export interface Friend extends Partial<User> {
	channel_id?: string;
	friendship_id: string;
	accepted: boolean;
	sender: boolean;
}

export interface Setup {
	user: User;
	friends: Friend[];
	servers: Record<string, Server>;
}

export interface DefaultResponse {
	message: string;
}

export interface Message {
	id: string;
	author: Partial<User>;
	server_id: string;
	channel_id: string;
	content: any;
	mentions_users: string[];
	mentions_channels: string[];
	updated_at: string;
	created_at: string;
}

export type ActorMessageTypes = 'channel:message';
