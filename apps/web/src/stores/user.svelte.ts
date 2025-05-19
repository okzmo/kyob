import type { User } from '../types/types';

class UserStore {
	user = $state<User>();
	mention = $state(false);
}

export const userStore = new UserStore();
