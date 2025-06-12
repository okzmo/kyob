import {
	LocalTrackPublication,
	RemoteParticipant,
	RemoteTrack,
	RemoteTrackPublication,
	Room,
	RoomEvent,
	Track,
	VideoPresets
} from 'livekit-client';
import { userStore } from './user.svelte';
class RTC {
	currentVC = $state<Room>();

	async connectToRoom(token: string) {
		//TODO: change to using user's settings
		const room = new Room({
			audioCaptureDefaults: {
				autoGainControl: false,
				echoCancellation: false,
				noiseSuppression: false
			},
			adaptiveStream: true,
			dynacast: true,
			videoCaptureDefaults: {
				resolution: VideoPresets.h720.resolution
			},
			publishDefaults: {
				screenShareEncoding: {
					maxBitrate: 1_500_000,
					maxFramerate: 60
				},
				dtx: true
			}
		});

		room
			.on(RoomEvent.TrackSubscribed, this.#handleTrackSubscribed)
			.on(RoomEvent.TrackUnsubscribed, this.#handleTrackUnsubscribed)
			.on(RoomEvent.ActiveSpeakersChanged, this.#handleActiveSpeakersChange)
			.on(RoomEvent.Disconnected, this.#handleDisconnect)
			.on(RoomEvent.LocalTrackUnpublished, this.#handleLocalTrackUnpublished)
			.on(RoomEvent.LocalTrackPublished, this.#handleLocalTrackPublished);

		await room.connect(import.meta.env.VITE_LIVEKIT_URL, token);

		this.currentVC = room;
	}

	async quitRoom() {
		if (!this.currentVC) return;
		await this.currentVC.disconnect();
	}

	#handleTrackSubscribed(
		track: RemoteTrack,
		_: RemoteTrackPublication,
		participant: RemoteParticipant
	) {
		switch (track.kind) {
			case Track.Kind.Audio:
				{
					const element = track.attach();
					document.body.appendChild(element);
				}
				break;
			case Track.Kind.Video:
				{
					if (track.source === 'screen_share') {
						setTimeout(() => {
							const vidElement = document.getElementById(
								`${participant.identity}-video-element`
							) as HTMLMediaElement;
							track.attach(vidElement);
						});
					}
				}
				break;
		}
	}

	#handleTrackUnsubscribed(track: RemoteTrack) {
		track.detach();
	}

	#handleActiveSpeakersChange() {}

	#handleDisconnect() {}

	#handleLocalTrackUnpublished(publication: LocalTrackPublication) {
		publication.track?.detach();
	}

	#handleLocalTrackPublished(publication: LocalTrackPublication) {
		if (publication.track?.kind === 'video' && publication.track?.source === 'screen_share') {
			const vidElement = document.getElementById(
				`${userStore.user!.id}-video-element`
			) as HTMLMediaElement;
			publication.track.attach(vidElement);
		}
	}
}

export const rtc = new RTC();
