import type { Friend, User } from '../types/types';

class UserStore {
	user = $state<User>();
	friends = $state<Friend[]>([]);
	mention = $state(false);

	acceptFriend(id: string) {
		for (const friend of this.friends) {
			if (friend.id === id) {
				friend.accepted = true;
			}
		}
	}

	deleteFriend(id: string) {
		const friendIdx = this.friends.findIndex((f) => f.id === id);
		this.friends.splice(friendIdx, 1);
	}
}

export const userStore = new UserStore();
