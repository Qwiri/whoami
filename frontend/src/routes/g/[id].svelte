<script lang="ts">
	import Chat from '../../Components/Chat.svelte';
	import Navbar from '../../Components/Navbar.svelte';

	import { page } from '$app/stores';
	import CharacterGrid from '../../Components/CharacterGrid.svelte';

	enum gameState {
		EnterName,
		ChooseCharacter,
		Lobby,
		Ingame
	}

	let gameStatus = gameState.Ingame;

	let ingameName = '';
	let characterName = '';

	function joinWithName(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			gameStatus = gameState.ChooseCharacter;
		}
	}
	function chooseCharacter(e: MouseEvent, p: { name: string; selected: boolean }) {
		characterName = p.name;
		gameStatus = gameState.Lobby;
		e.stopPropagation();
	}
</script>

<div id="container">
	<Navbar name={characterName} ingame={[gameState.Ingame, gameState.Lobby].includes(gameStatus)} />
	<div id="content">
		{#if gameStatus === gameState.Ingame}
			<CharacterGrid />
			<Chat />
		{:else if gameStatus === gameState.EnterName}
			<div id="lobbyContent">
				<h1>Lobby <span class="green">{$page.params.id}</span></h1>
				<input
					bind:value={ingameName}
					on:keyup={joinWithName}
					type="text"
					placeholder="Enter name"
				/>
			</div>
		{:else if gameStatus === gameState.ChooseCharacter}
			<div id="content">
				<h2>choose your character</h2>
				<CharacterGrid onLeftClick={chooseCharacter} onRightClick={(e) => e.stopPropagation()} />
			</div>
		{:else if gameStatus === gameState.Lobby}
			<div id="lobbyContent">
				<h1>In Lobby: <span class="green">{1} / 2</span></h1>
				<div>
					<p class="playerCard">{ingameName}</p>
				</div>
			</div>
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
		justify-content: center;
		gap: 4rem;
	}
	#lobbyContent {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}
	.green {
		color: #65d46e;
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
		color: black;
		border-radius: 0.5rem;
	}
</style>
