import type { Friend, User } from '../types/types';

class UserStore {
	user = $state<User>();
	friends = $state<Friend[]>([]);
	mention = $state(false);

	addFriend(friend: Friend) {
		if (Array.isArray(this.friends)) {
			this.friends.push(friend);
		} else {
			this.friends = [friend];
		}
	}

	getFriend(id: string) {
		return this.friends?.find((f) => f.id === id);
	}

	acceptFriend({
		friendshipId,
		friend,
		sender = false
	}: {
		friendshipId: string;
		friend?: Friend;
		sender?: boolean;
	}) {
		if (!sender) {
			const friendIdx = this.friends.findIndex((f) => f.friendship_id === friendshipId);
			this.friends[friendIdx].accepted = true;
		} else if (friend) {
			if (Array.isArray(this.friends)) {
				this.friends.push(friend);
			} else {
				this.friends = [friend];
			}
		}
	}

	setFriendChannelId(friendshipId: string, channelId: string) {
		const friendIdx = this.friends.findIndex((f) => f.friendship_id === friendshipId);
		this.friends[friendIdx].channel_id = channelId;
	}

	deleteFriend(friendshipId: string) {
		const friendIdx = this.friends.findIndex((f) => f.friendship_id === friendshipId);
		if (friendIdx > -1) this.friends.splice(friendIdx, 1);
	}
}

export const userStore = new UserStore();
