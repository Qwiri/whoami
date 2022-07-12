import { writable, type Writable } from 'svelte/store';

export interface Card {
	name: string;
	avatar: string;
}

export interface Pack {
	name: string;
	icon: string;
	description: string;
	avatars: Card[];
}

export interface ChatMessage {
	message: string;
	user: string;
}
export const packs: Writable<Pack[]> = writable([]);
export const selectedPack: Writable<Pack> = writable();
export const cards: Writable<Card[]> = writable([]);
export const users: Writable<string[]> = writable([]);
export const messages: Writable<ChatMessage[]> = writable([]);
export const selectedCard: Writable<Card> = writable();
export const tentativeCard: Writable<Card> = writable();
export const ingameName: Writable<string> = writable();
export const currentLives: Writable<number> = writable();
export const maxLives: Writable<number> = writable();
export const winnerName: Writable<string> = writable('');
export const winnerID: Writable<number> = writable();

// character
// packs (leader)
// leader?
// selected
// ws
