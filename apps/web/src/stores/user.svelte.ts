import type { Emoji, Friend, User } from '../types/types';
import { sounds } from './audio.svelte';
import { rtc } from './rtc.svelte';
import { serversStore } from './servers.svelte';

interface DM {
  friendId: string;
  channelId: string;
  avatar: string;
  username: string;
}

class UserStore {
  user = $state<User>();
  friends = $state<Friend[]>([]);
  emojis = $state<Emoji[]>([]);
  callTokens = $state<Record<string, string>>({});
  mention = $state(false);
  mute = $state(false);
  deafen = $state(false);

  getDms() {
    const global = serversStore.getServer('global');
    if (!global) return [];
    const dms: DM[] = [];

    for (const channel of Object.values(global.channels)) {
      if (Number(channel.last_message_read) < Number(channel.last_message_sent)) {
        const friend = channel.users?.find((u) => u.id !== this.user!.id);
        if (!friend) continue;

        dms.push({
          friendId: friend.id!,
          channelId: channel.id,
          avatar: friend.avatar!,
          username: friend.username!
        });
      }
    }

    return dms;
  }

  async toggleMute() {
    this.mute = !this.mute;
    sounds.playSound(this.mute ? 'mute-on' : 'mute-off');
    if (rtc.currentVC) await rtc.toggleMute();
  }

  async toggleDeafen() {
    this.deafen = !this.deafen;
    sounds.playSound(this.deafen ? 'mute-on' : 'mute-off');
    if (rtc.currentVC) await rtc.toggleDeafen();
  }

  addEmojis(emojis: Emoji[]) {
    if (Array.isArray(this.emojis)) {
      this.emojis.push(...emojis);
    } else {
      this.emojis = [...emojis];
    }
  }

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

  modifyFriend(friendId: string, informations: Partial<User>) {
    const friend = this.friends.find((f) => f.id === friendId);
    if (!friend) return;

    if (informations.display_name) friend.display_name = informations.display_name;
    if (informations.avatar) friend.avatar = informations.avatar;
    if (informations.about) friend.about = informations.about;
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
