import { fromBinary } from '@bufbuild/protobuf';
import { timestampDate } from '@bufbuild/protobuf/wkt';
import ky from 'ky';
import { err, ok, type Result } from 'neverthrow';
import { print } from 'utils/print';
import { WSMessageSchema } from '../gen/types_pb';
import type {
	AcceptFriendErrors,
	AddFriendErrors,
	CallErrors,
	CreateChannelErrors,
	CreateInviteErrors,
	CreateMessageErrors,
	CreateServerErrors,
	DeleteChannelErrors,
	DeleteEmojiErrors,
	DeleteFriendErrors,
	DeleteMessageErrors,
	DeleteServerErrors,
	EmojiErrors,
	GetUserErrors,
	JoinServerErrors,
	LeaveServerErrors,
	MessagesErrors,
	SetupErrors,
	StandardError,
	UpdateAccountErrors,
	UpdateAvatarErrors
} from '../types/errors';
import type {
	AcceptFriendType,
	AddEmojisType,
	AddFriendType,
	CreateChannelType,
	CreateMessageType,
	CreateServerType,
	DeleteFriendType,
	EditMessageType,
	JoinServerType,
	UpdateAccountType,
	UpdateAvatarType,
	UpdateProfileType
} from '../types/schemas';
import type {
	Channel,
	ChannelTypes,
	Emoji,
	Friend,
	LastState,
	Message,
	Server,
	Setup,
	User
} from '../types/types';
import { sounds } from './audio.svelte';
import { core } from './core.svelte';
import { serversStore } from './servers.svelte';
import { userStore } from './user.svelte';
import { windows } from './windows.svelte';

const client = ky.create({
	prefixUrl: `${import.meta.env.VITE_API_URL}/authenticated`,
	credentials: 'include',
	retry: 2,
	timeout: 10000
});

class Backend {
	wsConn = $state<WebSocket>();

	setupWebsocket(userId: string) {
		const ws = new WebSocket(`ws://localhost:3000/v1/authenticated/connect/${userId}`);
		if (!ws) return;

		this.wsConn = ws;
		ws.onopen = () => {
			print('Connection established');
			window.setInterval(() => {
				ws.send('heartbeat');
			}, 10 * 1000);
		};

		ws.onmessage = async (event) => {
			if (event.data === 'heartbeat') return;

			const arrayBuffer = await event.data.arrayBuffer();
			const uint8Array = new Uint8Array(arrayBuffer);
			const wsMess = fromBinary(WSMessageSchema, uint8Array, {
				readUnknownFields: false
			});
			switch (wsMess.content.case) {
				case 'chatMessage':
					{
						if (!wsMess.content.value) return;
						const value = wsMess.content.value;
						const contentStr = new TextDecoder().decode(value?.content);
						const attachments = new TextDecoder().decode(value?.attachments);

						const message: Message = {
							id: value.id,
							author_id: value.authorId,
							server_id: value.serverId,
							channel_id: value.channelId,
							content: JSON.parse(contentStr),
							everyone: value.everyone,
							mentions_users: value.mentionsUsers,
							mentions_channels: value.mentionsChannels,
							attachments: attachments.length > 0 ? JSON.parse(attachments) : [],
							updated_at: timestampDate(value.createdAt!).toISOString(),
							created_at: timestampDate(value.createdAt!).toISOString()
						};
						serversStore.addMessage(value?.serverId, message);

						const isAMention =
							message.mentions_users.includes(userStore.user!.id) || message.everyone;
						const isADM = message.server_id === 'global';
						const windowIsActive = windows.getActiveWindow()?.channelId === message.channel_id;

						if (
							(isAMention || isADM) &&
							message.author_id !== userStore.user!.id &&
							!windowIsActive
						) {
							sounds.playSound('notification');
						}
					}
					break;
				case 'channelCreation':
					{
						if (!wsMess.content.value) return;
						const value = wsMess.content.value;
						const channel: Channel = {
							id: value.id,
							name: value.name,
							type: value.type as ChannelTypes,
							unread: false,
							last_message_sent: '',
							last_message_read: '',
							last_mentions: [],
							x: value.x,
							y: value.y,
							users: value.users.map((u) => ({
								id: u.id,
								avatar: u.avatar,
								username: u.username,
								display_name: u.displayName
							})),
							voice_users: []
						};
						serversStore.addChannel(value.serverId, channel);
					}
					break;
				case 'channelRemoved':
					{
						if (!wsMess.content.value) return;
						const value = wsMess.content.value;
						serversStore.removeChannel(value.serverId, value.channelId);
						const channelWindow = windows.getWindow({ channelId: value.channelId });
						if (channelWindow) windows.closeDeadWindow(channelWindow.id);
					}
					break;
				case 'newUser':
					{
						if (!wsMess.content.value) return;
						const value = wsMess.content.value;
						const about = new TextDecoder().decode(value.user?.about);

						const newUser: Partial<User> = {
							id: value.user?.id,
							username: value.user?.username,
							display_name: value.user?.displayName,
							avatar: value.user?.avatar
						};

						if (about.length > 0) newUser.about = JSON.parse(about);
						serversStore.addMember(value.serverId, newUser);
					}
					break;
				case 'userConnect':
					{
						if (!wsMess.content.value) return;
						const value = wsMess.content.value;
						setTimeout(() => {
							serversStore.connectUser(value.serverId, value.userId, value.users, value.type);
						}, 500);
					}
					break;
				case 'userDisconnect':
					{
						if (!wsMess.content.value) return;
						const value = wsMess.content.value;
						if (value.userId === userStore.user!.id) return;
						serversStore.disconnectUser(value.serverId, value.userId, value.type);
					}
					break;
				case 'deleteMessage':
					{
						if (!wsMess.content.value) return;
						const value = wsMess.content.value;
						serversStore.deleteMessage(value.serverId, value.channelId, value.messageId);
					}
					break;
				case 'editMessage':
					{
						if (!wsMess.content.value) return;
						const value = wsMess.content.value;
						const contentStr = new TextDecoder().decode(value.content);
						const message = serversStore.editMessage(
							value.serverId,
							value.channelId,
							value.messageId,
							JSON.parse(contentStr),
							value.mentionsUsers,
							value.mentionsChannels,
							timestampDate(value.updatedAt!).toISOString()
						);

						if (
							(message?.mentions_users?.includes(userStore.user!.id) || message?.everyone) &&
							message.author_id !== userStore.user!.id
						) {
							sounds.playSound('notification');
						}
					}
					break;
				case 'friendInvite':
					{
						if (!wsMess.content.value) return;
						const value = wsMess.content.value;
						const about = new TextDecoder().decode(value.user?.about);

						const newFriend: Friend = {
							id: value.user?.id,
							display_name: value.user?.displayName,
							avatar: value.user?.avatar,
							friendship_id: value.inviteId,
							accepted: false,
							sender: false
						};
						if (about.length > 0) newFriend.about = JSON.parse(about);

						userStore.addFriend(newFriend);
					}
					break;
				case 'acceptFriend':
					{
						if (!wsMess.content.value) return;
						const value = wsMess.content.value;
						const about = new TextDecoder().decode(value.user?.about);

						if (value.sender) {
							const newFriend: Friend = {
								id: value.user?.id,
								display_name: value.user?.displayName,
								avatar: value.user?.avatar,
								friendship_id: value.inviteId,
								channel_id: value.channelId,
								accepted: true,
								sender: true
							};
							if (about.length > 0) newFriend.about = JSON.parse(about);

							userStore.acceptFriend({
								friendshipId: value.inviteId,
								friend: newFriend,
								sender: true
							});
						} else {
							userStore.setFriendChannelId(value.inviteId, value.channelId);
						}
					}
					break;
				case 'deleteFriend':
					{
						if (!wsMess.content.value) return;
						const value = wsMess.content.value;
						userStore.deleteFriend(value.inviteId);
						const friendChatWindow = windows.getWindow({ friendId: value.userId });
						if (friendChatWindow) windows.closeDeadWindow(friendChatWindow.id);
					}
					break;
				case 'userChanged':
					{
						if (!wsMess.content.value) return;
						const value = wsMess.content.value;
						const about = new TextDecoder().decode(value.userInformations?.about);
						const facts = new TextDecoder().decode(value.userInformations?.facts);
						const links = new TextDecoder().decode(value.userInformations?.links);

						const user: Partial<User> = {
							avatar: value.userInformations?.avatar,
							banner: value.userInformations?.banner,
							username: value.userInformations?.username,
							display_name: value.userInformations?.displayName,
							main_color: value.userInformations?.mainColor
						};

						if (about.length > 0) user.about = JSON.parse(about);
						if (facts.length > 0) user.facts = JSON.parse(facts);
						if (links.length > 0) user.links = JSON.parse(links);

						if (value.serverId) {
							serversStore.modifyMember(value.serverId, value.userId, user);
						} else {
							userStore.modifyFriend(value.userId, user);
						}

						core.modifyProfile(value.userId, user);
					}
					break;
				case 'connectToCall':
					{
						if (!wsMess.content.value) return;
						const value = wsMess.content.value;
						serversStore.connectUserToCall(value.serverId, value.channelId, value.userId);

						if (serversStore.isInCall(value.serverId, value.channelId, userStore.user!.id)) {
							sounds.playSound('call-on');
						}
					}
					break;
				case 'disconnectFromCall':
					{
						if (!wsMess.content.value) return;
						const value = wsMess.content.value;
						serversStore.disconnectUserFromCall(value.serverId, value.channelId, value.userId);

						if (serversStore.isInCall(value.serverId, value.channelId, userStore.user!.id)) {
							sounds.playSound('call-off');
						}
					}
					break;
				case 'callUsers':
					{
						if (!wsMess.content.value) return;
						const value = wsMess.content.value;
						for (const call of value.callUsers) {
							serversStore.connectUserToCall(call.serverId, call.channelId, call.userId);
						}
					}
					break;
			}
		};

		ws.onclose = (event) => {
			print('Connection closed:', event);
		};

		ws.onerror = (error) => {
			console.error('WebSocket error:', error);
		};
	}

	async getSetup(): Promise<Result<Setup, SetupErrors>> {
		try {
			const res = await client.get('setup', {
				headers: {
					'Content-Type': 'application/json'
				}
			});

			const data = (await res.json()) as Setup;
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok(data);
		} catch (error) {
			const errBody = await (error as StandardError).response?.json();

			switch (true) {
				case errBody.status === 401:
					return err({ code: 'ERR_UNAUTHORIZED', error: errBody.error });
				default:
					return err({ code: 'ERR_UNKNOWN', error: errBody.error });
			}
		}
	}

	async createServer(body: CreateServerType): Promise<Result<Server, CreateServerErrors>> {
		try {
			const formData = new FormData();
			formData.append('name', body.name);
			formData.append('avatar', body.avatar);
			formData.append('crop', JSON.stringify(body.crop));
			formData.append('private', String(body.private));
			formData.append('x', String(body.x));
			formData.append('y', String(body.y));

			if (body.description) formData.append('description', JSON.stringify(body.description));

			const res = await client.post('server', {
				body: formData
			});

			const data = (await res.json()) as Server;
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok(data);
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			switch (true) {
				case errBody.status === 400:
					return err({ code: 'ERR_VALIDATION_FAILED', error: errBody.error });
				case errBody.status === 403 && errBody.code === 'ERR_TOO_MANY_SERVERS':
					return err({ code: 'ERR_TOO_MANY_SERVERS', error: errBody.error });
				default:
					return err({ code: 'ERR_UNKNOWN', error: errBody.error });
			}
		}
	}

	async joinServer(body: JoinServerType): Promise<Result<Server, JoinServerErrors>> {
		try {
			const res = await client.post('server/join', {
				body: JSON.stringify(body)
			});

			const data = (await res.json()) as { server: Server };
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok(data.server);
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			switch (true) {
				case errBody.status === 400:
					return err({ code: errBody.code, error: errBody.error });
				case errBody.status === 404:
					return err({ code: 'ERR_INVITE_SERVER_NOT_FOUND', error: errBody.error });
				default:
					return err({ code: 'ERR_UNKNOWN', error: errBody.error });
			}
		}
	}

	async deleteServer(serverId: string): Promise<Result<void, DeleteServerErrors>> {
		try {
			const res = await client.delete(`servers/${serverId}`);

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();

			switch (true) {
				case errBody.status === 401:
					return err({ code: 'ERR_UNAUTHORIZED', error: errBody.error });
				default:
					return err({ code: 'ERR_UNKNOWN', error: errBody.error });
			}
		}
	}

	async leaveServer(serverId: string): Promise<Result<void, LeaveServerErrors>> {
		try {
			const res = await client.post(`server/${serverId}/leave`);

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async createChannel(
		serverId: string,
		body: CreateChannelType
	): Promise<Result<void, CreateChannelErrors>> {
		try {
			const res = await client.post(`channels/${serverId}`, {
				body: JSON.stringify(body)
			});

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();

			switch (true) {
				case errBody.status === 401:
					return err({ code: 'ERR_UNAUTHORIZED', error: errBody.error });
				default:
					return err({ code: 'ERR_UNKNOWN', error: errBody.error });
			}
		}
	}

	async deleteChannel(
		serverId: string,
		channelId: string
	): Promise<Result<void, DeleteChannelErrors>> {
		try {
			const res = await client.delete(`channels/${serverId}/${channelId}`);

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();

			switch (true) {
				case errBody.status === 401:
					return err({ code: 'ERR_UNAUTHORIZED', error: errBody.error });
				default:
					return err({ code: 'ERR_UNKNOWN', error: errBody.error });
			}
		}
	}

	async sendMessage(
		serverId: string,
		channelId: string,
		body: CreateMessageType
	): Promise<Result<void, CreateMessageErrors>> {
		try {
			const formData = new FormData();
			formData.append('type', 'SEND');
			formData.append('content', JSON.stringify(body.content));
			formData.append('everyone', body.everyone ? 'true' : 'false');
			body.mentions_users?.forEach((user) => formData.append('mentions_users[]', user));
			body.attachments?.forEach((attachment) => formData.append('attachments[]', attachment));

			const res = await client.post(`messages/${serverId}/${channelId}`, {
				body: formData
			});

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();

			switch (true) {
				case errBody.status === 400:
					return err({ code: 'ERR_VALIDATION_FAILED', error: errBody.error });
				case errBody.status === 413:
					return err({ code: 'ERR_MESSAGE_TOO_BIG', error: errBody.error });
				default:
					return err({ code: 'ERR_UNKNOWN', error: errBody.error });
			}
		}
	}

	async editMessage(
		serverId: string,
		channelId: string,
		messageId: string,
		body: EditMessageType
	): Promise<Result<void, CreateMessageErrors>> {
		const formData = new FormData();
		formData.append('type', 'EDIT');
		formData.append('content', JSON.stringify(body.content));
		body.mentions_users?.forEach((user) => formData.append('mentions_users[]', user));

		try {
			const res = await client.patch(`messages/${serverId}/${channelId}/${messageId}`, {
				body: formData
			});

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();

			switch (true) {
				case errBody.status === 400:
					return err({ code: 'ERR_VALIDATION_FAILED', error: errBody.error });
				default:
					return err({ code: 'ERR_UNKNOWN', error: errBody.error });
			}
		}
	}

	async getMessages(channelId: string): Promise<Result<Message[], MessagesErrors>> {
		try {
			const res = await client.get(`messages/${channelId}`);

			const data = (await res.json()) as Message[];
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok(data);
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async createInvite(serverId: string): Promise<Result<string, CreateInviteErrors>> {
		try {
			const res = await client.get(`server/create_invite/${serverId}`);

			const data = (await res.json()) as { invite_link: string };
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok(data.invite_link);
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async getUserProfile(userId: string): Promise<Result<User, GetUserErrors>> {
		try {
			const res = await client.get(`user/${userId}`);

			const data = (await res.json()) as User;
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok(data);
		} catch (error) {
			const errBody = await (error as StandardError).response.json();

			switch (true) {
				case errBody.status === 404:
					return err({ code: 'ERR_USER_NOT_FOUND', error: errBody.error });
				default:
					return err({ code: 'ERR_UNKNOWN', error: errBody.error });
			}
		}
	}

	async deleteMessage(
		serverId: string,
		channelId: string,
		messageId: string
	): Promise<Result<void, DeleteMessageErrors>> {
		try {
			const res = await client.delete(`messages/${serverId}/${channelId}/${messageId}`);

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();

			switch (true) {
				case errBody.status === 400:
					return err({ code: 'ERR_VALIDATION_FAILED', error: errBody.error });
				default:
					return err({ code: 'ERR_UNKNOWN', error: errBody.error });
			}
		}
	}

	async addFriend(body: AddFriendType): Promise<Result<void, AddFriendErrors>> {
		try {
			const res = await client.post('friends/add', {
				body: JSON.stringify(body)
			});

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();

			switch (true) {
				case errBody.status === 404:
					return err({ code: 'ERR_USER_NOT_FOUND', error: errBody.error });
				case errBody.status === 403:
					return err({ code: 'ERR_ADDING_ITSELF', error: errBody.error });
				default:
					return err({ code: 'ERR_UNKNOWN', error: errBody.error });
			}
		}
	}

	async acceptFriend(body: AcceptFriendType): Promise<Result<void, AcceptFriendErrors>> {
		try {
			const res = await client.post('friends/accept', {
				body: JSON.stringify(body)
			});

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();

			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async deleteFriend(body: DeleteFriendType): Promise<Result<void, DeleteFriendErrors>> {
		try {
			const res = await client.post('friends/delete', {
				body: JSON.stringify(body)
			});

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();

			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async updateAccount(body: UpdateAccountType): Promise<Result<void, UpdateAccountErrors>> {
		try {
			const res = await client.post('user/update_account', {
				body: JSON.stringify(body)
			});

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();

			switch (true) {
				case errBody.status === 403 && errBody.code === 'ERR_USERNAME_IN_USE':
					return err({ code: 'ERR_USERNAME_IN_USE', error: errBody.error });
				case errBody.status === 403 && errBody.code === 'ERR_EMAIL_IN_USE':
					return err({ code: 'ERR_EMAIL_IN_USE', error: errBody.error });
				default:
					return err({ code: 'ERR_UNKNOWN', error: errBody.error });
			}
		}
	}

	async logout() {
		try {
			await client.post('logout');
		} catch (error) {
			console.error(error);
		}
	}

	async updateAvatar(
		body: UpdateAvatarType
	): Promise<Result<{ banner: string; avatar: string; main_color: string }, UpdateAvatarErrors>> {
		try {
			const formData = new FormData();
			formData.append('avatar', body.avatar);
			formData.append('crop_banner', JSON.stringify(body.crop_banner));
			formData.append('crop_avatar', JSON.stringify(body.crop_avatar));
			if (body.main_color) formData.append('main_color', body.main_color);

			const res = await client.post('user/update_avatar', {
				body: formData
			});

			const data = (await res.json()) as { banner: string; avatar: string; main_color: string };
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok(data);
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async updateProfile(body: UpdateProfileType): Promise<Result<void, UpdateAvatarErrors>> {
		try {
			const res = await client.post('user/update_profile', {
				body: JSON.stringify(body)
			});

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async connectToCall(
		serverId: string,
		channelId: string
	): Promise<Result<{ token: string }, CallErrors>> {
		try {
			const res = await client.post(`channels/${serverId}/${channelId}/join_call`);

			const data = (await res.json()) as { token: string };
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok(data);
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async disconnectFromCall(serverId: string, channelId: string): Promise<Result<void, CallErrors>> {
		try {
			const res = await client.post(`channels/${serverId}/${channelId}/quit_call`);

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();
			return err({ code: 'ERR_UNKNOWN', error: errBody.error });
		}
	}

	async uploadEmojis(body: AddEmojisType): Promise<Result<Emoji[], EmojiErrors>> {
		try {
			const formData = new FormData();
			for (let i = 0; i < body.emojis.length; ++i) {
				formData.append('emojis[]', body.emojis[i]);
				formData.append('shortcodes[]', body.shortcodes[i]);
			}

			const res = await client.post(`user/upload_emojis`, {
				body: formData
			});

			const data = (await res.json()) as Emoji[];
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok(data);
		} catch (error) {
			const errBody = await (error as StandardError).response.json();

			switch (true) {
				case errBody.status === 400 && errBody.code === 'ERR_MISSING_EMOJIS':
					return err({ code: 'ERR_MISSING_EMOJIS', error: errBody.error });
				case errBody.status === 400 && errBody.code === 'ERR_MISSING_SHORTCODES':
					return err({
						code: 'ERR_MISSING_SHORTCODES',
						error: 'Some shortcodes are missing or empty'
					});
				case errBody.status === 400 && errBody.code === 'ERR_EMOJIS_INVALID':
					return err({
						code: 'ERR_EMOJIS_INVALID',
						error: 'Your emojis are invalid. Please refer to the above requirements.'
					});
				case errBody.status === 400 && errBody.code === 'ERR_SHORTCODES_INVALID':
					return err({
						code: 'ERR_SHORTCODES_INVALID',
						error:
							'The given shortcodes are invalid. The shortcodes must be lowercase, without spaces. e.g. my_emoji_1'
					});
				default:
					return err({ code: 'ERR_UNKNOWN', error: errBody.error });
			}
		}
	}

	async deleteEmoji(emojiId: string): Promise<Result<void, DeleteEmojiErrors>> {
		try {
			const res = await client.delete(`user/delete_emoji/${emojiId}`);

			const data = await res.json();
			if (!res.ok) {
				return err({ code: 'ERR_UNKNOWN', error: '', cause: data });
			}

			return ok();
		} catch (error) {
			const errBody = await (error as StandardError).response.json();

			switch (true) {
				case errBody.status === 403:
					return err({ code: 'ERR_FORBIDDEN', error: errBody.error });
				default:
					return err({ code: 'ERR_UNKNOWN', error: errBody.error });
			}
		}
	}

	async saveState(body: LastState) {
		navigator.sendBeacon(
			`${import.meta.env.VITE_API_URL}/authenticated/save_state`,
			JSON.stringify(body)
		);
	}
}

export const backend = new Backend();
