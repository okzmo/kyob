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
  width: number;
  height: number;
  x: number;
  y: number;
  tab: 'chat' | 'call';
}

export interface Channel {
  id: string;
  name: string;
  type: ChannelTypes;
  x: number;
  y: number;
  unread: boolean;
  last_message_sent?: string;
  last_message_read?: string;
  last_mentions?: string[];
  messages?: Message[];
  users?: Partial<User>[];
  voice_users: {
    user_id: string;
    deafen: boolean;
    mute: boolean;
  }[];
}

export interface Server {
  id: string;
  owner_id: string;
  name: string;
  avatar: string;
  banner: string;
  description?: any;
  main_color?: string;
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
  about?: any;
  facts: Fact[];
  links: Link[];
  rpm_avatar_id: string;
  rpm_token: string;
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
  emojis: Emoji[];
  servers: Record<string, Server>;
}

export interface DefaultResponse {
  message: string;
}

export interface Attachment {
  id: string;
  url: string;
  file_name: string;
  file_size: string;
  type: string;
}

export interface Message {
  id: string;
  author_id: string;
  server_id: string;
  channel_id: string;
  content: any;
  everyone: boolean;
  mentions_users: string[];
  mentions_channels: string[];
  attachments: Attachment[];
  updated_at: string;
  created_at: string;
}

export interface LastState {
  channel_ids: string[];
  last_message_ids: string[];
  mentions_ids: string[][];
}

export interface Emoji {
  id: string;
  url: string;
  shortcode: string;
}
