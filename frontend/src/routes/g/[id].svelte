<script lang="ts">
	import Chat from '../../Components/Chat.svelte';
	import Navbar from '../../Components/Navbar.svelte';

	import { page } from '$app/stores';
	import CharacterGrid from '../../Components/CharacterGrid.svelte';

	import { Gobby, ö } from 'gobby-ts';
	import type { Message } from 'gobby-ts';
	import { onMount } from 'svelte';
	import {
		cards,
		currentLives,
		ingameName,
		maxLives,
		messages,
		packs,
		selectedCard,
		users,
		winnerName,
		winnerID,
		type Card,
		type ChatMessage,
		type Pack,
		selectedPackIndex
	} from '../../stores';
	import WinningScreen from '../../Components/WinningScreen.svelte';
	import LobbyPeek, { type PeekContent } from '../../Components/LobbyPeek.svelte';

	enum GameState {
		EnterName,
		ChooseCharacter,
		Lobby,
		Ingame,
		End
	}

	let gameStatus = GameState.EnterName;

	let _ingameName = '';

	let gobby: Gobby;

	// Peek lobby before joining
	let peek: undefined | null | PeekContent;
	const check = async () => {
		const resp = await fetch(`https://backend.wai.sap.lol/lobby/peek/${$page.params.id}`);
		if (resp.status === 404) {
			peek = null; // null = lobby not found
		} else {
			peek = (await resp.json()) as PeekContent;
		}
	};

	onMount(async () => {
		check();
		setInterval(check, 1000);

		try {
			gobby = new Gobby(`wss://backend.wai.sap.lol/lobby/socket/${$page.params.id}`);
			await gobby.connect();
		} catch (e) {
			console.error(e);
		}

		// SELECTED_PACK_CHANGED is called whenever the selected pack changes.
		// {selectedPackIndex: number}
		gobby.handle('SELECTED_PACK_CHANGED', (msg: Message) => {
			if (msg.args) {
				selectedPackIndex.set(msg.args[0] as number);
			}
		});

		// LIST returns a list of all users in the lobby.
		// {user1: string}, {user2: string}, ...
		gobby.handle('LIST', (msg: Message) => {
			if (msg.args) {
				users.set(msg.args as string[]);
			}
		});

		// CHAT is called whenever a message is sent in the chat.
		// {username: string}, {message: string}
		gobby.handle('CHAT', (msg: Message) => {
			if (msg.args) {
				const newMessage: ChatMessage = {
					message: msg.args[1] as string,
					user: msg.args[0] as string
				};
				messages.update((old) => {
					old.unshift(newMessage);
					return old;
				});
			}
		});

		// STATE_CHANGE is called whenever the game state changes.
		// {state: number}
		gobby.handle('STATE_CHANGE', (msg: Message) => {
			const states = [GameState.Lobby, GameState.ChooseCharacter, GameState.Ingame, GameState.End];
			if (msg.args) {
				gameStatus = states[Math.log2(msg.args[1] as number)];
			}
		});

		// AVAILABLE_CHARACTERS returns a list of all available characters you can choose from.
		// {character1: Card}, {character2: Card}, ...
		gobby.handle('AVAILABLE_CHARACTERS', (msg: Message) => {
			if (msg.args) {
				cards.set(msg.args as Card[]);
			}
		});

		// LIVES returns the current amount of lives you have.
		// {lives: number}, {maxLives: number}
		gobby.handle('LIVES', (msg: Message) => {
			if (msg.args) {
				currentLives.set(msg.args[0] as number);
				maxLives.set(msg.args[1] as number);
			}
		});

		// WINNER returns the name of the winner and the reason why they won.
		gobby.handle('WINNER', (msg: Message) => {
			if (msg.args) {
				winnerName.set(msg.args[0] as string);
				winnerID.set(msg.args[1] as number);
			}
		});

		gobby.handle('PACKS', (msg: Message) => {
			if (msg.args) {
				packs.set(msg.args[0] as Pack[]);
			}
		});
	});

	async function joinWithName(event: KeyboardEvent) {
		if (event.key === 'Enter') {
			try {
				await gobby.join(_ingameName);
				ingameName.set(_ingameName);
			} catch (eeee) {
				console.error(eeee);
				return;
			}
		}
	}

	function chooseCharacter(event: MouseEvent, _player: Card, i: number) {
		event.stopPropagation();
		gobby.send(ö('SELECT_CHARACTER', i.toString())).then((msg) => {
			if (msg.args && msg.args[0] === 'OK') {
				let card: Card = {
					name: msg.args[1] as string,
					avatar: msg.args[2] as string
				};
				selectedCard.set(card);
			}
		});
	}

	function startGame() {
		gobby.send(ö('START'));
	}

	function sendChatMessage(text: string) {
		gobby.send(ö('CHAT', text));
	}

	function changePack(index: number) {
		gobby.send(ö('SELECT_PACK', index.toString()));
	}

	function guess(index: number) {
		gobby.send(ö('GUESS', index.toString()));
	}
</script>

<div id="container">
	<Navbar ingame={[GameState.Ingame, GameState.End].includes(gameStatus)} />
	<div id="content">
		{#if gameStatus === GameState.Ingame}
			<CharacterGrid onGuess={guess} />
			<Chat sendMessageCallback={sendChatMessage} />
		{:else if gameStatus === GameState.EnterName}
			<div id="lobbyContent">
				<h1>Lobby <span class="green">{$page.params.id}</span></h1>
				<input
					bind:value={_ingameName}
					on:keyup={joinWithName}
					type="text"
					placeholder="Enter name"
				/>
				<!-- Peek -->
				{#if peek !== undefined}
					<div id="peekContent">
						{#if peek === null}
							<h2 class="status unavailable">LOBBY UNAVAILABLE</h2>
						{:else}
							<LobbyPeek {peek} />
						{/if}
					</div>
				{/if}
			</div>
		{:else if gameStatus === GameState.ChooseCharacter}
			<div id="content">
				<h2>choose your character</h2>
				<CharacterGrid onLeftClick={chooseCharacter} onRightClick={(e) => e.stopPropagation()} />
			</div>
		{:else if gameStatus === GameState.Lobby}
			<div class="ewe">
				<div id="playerVS">
					<div>
						<h3>In Lobby: <span class="green">{$users.length} / 2</span></h3>
						<div id="playerVS">
							<p class="playerCard pc1">{$users[0] || ' '}</p>
							<p><b>VS</b></p>
							<p class="playerCard pc2">{$users[1] || ' '}</p>
						</div>
						<button
							id="startButton"
							class={$users.length >= 2 ? 'startPlayable' : 'gray'}
							on:click={startGame}>Start</button
						>
					</div>
					<div>
						<h1>Select a pack</h1>
						{#each $packs as pack, i}
							<div
								class="pack"
								class:selectedPack={i === $selectedPackIndex}
								on:click={() => changePack(i)}
							>
								<img src={pack.icon} alt="cover for {pack.name}" />
								<div>
									<h3>{pack.name}</h3>
									<p>{pack.description}</p>
								</div>
							</div>
						{/each}
					</div>
				</div>
				<Chat sendMessageCallback={sendChatMessage} />
			</div>
		{:else if gameStatus === GameState.End}
			<WinningScreen />
		{/if}
	</div>
</div>

<style lang="scss">
	#container {
		width: 100%;
		height: 100vh;

		display: flex;
		align-items: center;
	}
	.ewe {
		display: flex;
		width: 100%;
		justify-content: space-around;
	}
	#content {
		display: flex;
		width: 100%;
		height: 100%;
		justify-content: center;
		gap: 4rem;
	}
	#lobbyContent {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}
	#playerVS {
		display: flex;
		justify-content: center;
		align-items: center;
		gap: 1rem;
	}
	.selectedPack {
		background-image: linear-gradient(to right, #65d46e, #00bde7) !important;
	}
	.green {
		color: #65d46e;
	}
	.gray {
		color: black;
		background-color: gray;
	}
	.startPlayable {
		color: white;
		background-color: #65d46e;

		&:hover {
			cursor: pointer;
		}
	}
	#startButton {
		padding: 0.5rem 1rem;
		font-size: 1.2rem;
		border-radius: 0.5rem;
		border: none;

		p {
			margin: 0;
		}
	}
	input {
		background-color: #1e1e1e;
		border-radius: 2rem;
		border: none;
		padding: 1rem 2rem;
	}
	.playerCard {
		padding: 1rem;
		background-color: #d9d9d9;
		width: 10rem;
		color: #2e2e2e;
		border-radius: 0.5rem;
		font-weight: bold;
	}
	.pc1 {
		background-image: linear-gradient(to right, #c4d465, #ff736a);
	}
	.pc2 {
		background-image: linear-gradient(to right, #65d46e, #6aa6ff);
	}
	.pack {
		display: flex;
		align-items: center;
		background-image: linear-gradient(to left, #6f6f6f, #4d4d4d);
		border-radius: 1rem;
		margin: 0.2rem;

		&:hover {
			cursor: pointer;
		}
		p {
			margin-top: 0;
		}
		h1 {
			margin-bottom: 0;
		}
		img {
			height: 4rem;
			border-radius: 0.2rem;
			margin: 0.5rem;
		}
	}

	#peekContent {
		margin-top: 2rem;
	}
</style>
