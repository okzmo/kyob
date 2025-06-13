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
import { print } from 'utils/print';

class RTC {
	currentVC = $state<Room>();

	async prepareConnection() {
		print('Preparing connection...');

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

		await room.prepareConnection(import.meta.env.VITE_LIVEKIT_URL);
	}

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
			.on(RoomEvent.LocalTrackPublished, this.#handleLocalTrackPublished)
			.on(RoomEvent.Reconnecting, this.#handleReconnecting)
			.on(RoomEvent.Connected, this.#handleConnected)
			.on(RoomEvent.Reconnected, this.#handleReconnected);

		const startTime = Date.now();
		try {
			await room.connect(import.meta.env.VITE_LIVEKIT_URL, token);
			print(`Connection established in ${Date.now() - startTime}ms`);
		} catch (error) {
			console.error(`Connection failed after ${Date.now() - startTime}ms:`, error);
		}

		try {
			await room.localParticipant.setMicrophoneEnabled(!userStore.mute);
		} catch (error) {
			console.error('Error enabling camera and microphone:', error);
		}

		if (userStore.deafen) {
			room.remoteParticipants.forEach((participant) => {
				const audioTrack = participant.getTrackPublication(Track.Source.Microphone);
				audioTrack?.setEnabled(false);
			});
		}

		this.currentVC = room;
	}

	async quitRoom() {
		if (!this.currentVC) return;
		await this.currentVC.disconnect();
	}

	#handleTrackSubscribed(
		track: RemoteTrack,
		pub: RemoteTrackPublication,
		participant: RemoteParticipant
	) {
		switch (track.kind) {
			case Track.Kind.Audio:
				{
					const audioContainer = document.getElementById('voice-audio-container');

					if (audioContainer) {
						const element = track.attach();
						element.id = `audio-${participant.identity}`;
						audioContainer.appendChild(element);
					}

					if (userStore.deafen) {
						pub.setEnabled(false);
					}
				}
				break;
			case Track.Kind.Video:
				{
					if (track.source === 'screen_share') {
						const vidElement = document.getElementById(`${participant.identity}-video-element`);

						if (vidElement) {
							setTimeout(() => {
								track.attach(vidElement as HTMLMediaElement);
							});
						}
					}
				}
				break;
		}
	}

	async toggleMute() {
		if (!this.currentVC) return;

		const audioTrack = this.currentVC.localParticipant.getTrackPublication(Track.Source.Microphone);
		if (!audioTrack) return;

		if (userStore.mute) await audioTrack.mute();
		else await audioTrack.unmute();
	}

	async toggleDeafen() {
		if (!this.currentVC) return;

		this.currentVC.remoteParticipants.forEach((participant) => {
			const audioTrack = participant.getTrackPublication(Track.Source.Microphone);
			audioTrack?.setEnabled(!userStore.deafen);
		});
	}

	#handleTrackUnsubscribed(track: RemoteTrack) {
		track.detach();
	}

	#handleActiveSpeakersChange() {}

	#handleDisconnect() {
		this.currentVC = undefined;
	}

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

	#handleReconnecting() {
		print('Reconnecting...');
	}

	#handleReconnected() {
		print('Reconnected!');
	}

	#handleConnected() {
		print('Connected!');
	}
}

export const rtc = new RTC();
