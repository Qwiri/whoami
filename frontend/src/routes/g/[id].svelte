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
		selectedPack,
		users,
		winnerName,
		winnerID,
		type Card,
		type ChatMessage,
		type Pack
	} from '../../stores';
	import WinningScreen from '../../Components/WinningScreen.svelte';

	enum GameState {
		EnterName,
		ChooseCharacter,
		Lobby,
		Ingame,
		End
	}

	let gameStatus = GameState.EnterName;

	let _ingameName = '';
	let characterName = '';

	let gobby: Gobby;

	onMount(async () => {
		try {
			gobby = new Gobby(`wss://backend.wai.sap.lol/lobby/socket/${$page.params.id}`);
			const ws = await gobby.connect();

			gobby.handle('SELECTED_PACK_CHANGED', (msg: Message) => {
				if (msg.args) {
					selectedPack.set($packs[msg.args[0] as number]);
				}
			});

			gobby.handle('PACK', (msg: Message) => {
				if (msg.args) {
					selectedPack.set(msg.args[0] as Pack);
				}
			});

			gobby.handle('LIST', (msg: Message) => {
				if (msg.args) {
					users.set(msg.args as string[]);
				}
			});

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

			gobby.handle('STATE_CHANGE', (msg: Message) => {
				if (msg.args) {
					switch (msg.args[1] as number) {
						case 1 << 0:
							gameStatus = GameState.Lobby;
							break;
						case 1 << 1:
							// select character
							gameStatus = GameState.ChooseCharacter;
							break;
						case 1 << 2:
							// game halt, ingame
							gameStatus = GameState.Ingame;
							break;
						case 1 << 3:
							// gewonnen, dieser winning screen
							gameStatus = GameState.End;
							break;

						default:
							break;
					}
				}
			});

			gobby.handle('AVAILABLE_CHARACTERS', (msg: Message) => {
				if (msg.args) {
					cards.set(msg.args as Card[]);
				}
			});

			gobby.handle('LIVES', (msg: Message) => {
				if (msg.args) {
					currentLives.set(msg.args[0] as number);
					maxLives.set(msg.args[1] as number);
				}
			});

			gobby.handle('WINNER', (msg: Message) => {
				if (msg.args) {
					$winnerName = msg.args[0] as string;
					$winnerID = msg.args[1] as number;
					// winnningReason = msg.args[1] as string;
				}
			});
		} catch (e) {
			console.error(e);
		}
	});

	async function joinWithName(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			try {
				await gobby.join(_ingameName);
				ingameName.set(_ingameName);
			} catch (eeee) {
				console.error(eeee);
				return;
			}
			gobby.send(ö('PACKS')).then((msg) => {
				console.log('PACKS');
				if (msg.args) {
					packs.set(msg.args[0] as Pack[]);
				}
			});
		}
	}
	function chooseCharacter(e: MouseEvent, card: Card, i: number) {
		console.log('CHOOSE CHARACTER');
		gobby.send(ö('SELECT_CHARACTER', i.toString())).then((msg) => {
			if (msg.args && msg.args[0] === 'OK') {
				let card: Card = {
					name: msg.args[1] as string,
					avatar: msg.args[2] as string
				};
				selectedCard.set(card);
			}
		});
		e.stopPropagation();
	}

	function startGame() {
		if ($users.length >= 2) {
			gobby.send(ö('START'));
		}
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
			</div>
		{:else if gameStatus === GameState.ChooseCharacter}
			<div id="content">
				<h2>choose your character</h2>
				<CharacterGrid onLeftClick={chooseCharacter} onRightClick={(e) => e.stopPropagation()} />
			</div>
		{:else if gameStatus === GameState.Lobby}
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
							class:selectedPack={pack.name === $selectedPack?.name}
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
		padding: 0.5rem;
		border-radius: 0.2rem;

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
</style>
