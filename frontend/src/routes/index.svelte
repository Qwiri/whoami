<script lang="ts">
	import Navbar from '../Components/Navbar.svelte';

	let lobbyCode: string;

	function enterLobby(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			window.location.href = `/g/${lobbyCode}`;
		}
	}
	async function createGame() {
		const res = await fetch('https://backend.wai.sap.lol/lobby/create');
		console.log('res');
		console.log(res);
		const json = await res.json();
		console.log('JSON');
		console.log(json);
		window.location.href = `/g/${json.ID}`;
	}
</script>

<div id="root" class="horizontalContent">
	<Navbar />
	<div id="main" class="verticalContent">
		<div id="header">
			<img src="/banner.svg" alt="header image" />
			<h1>WHOAMI</h1>
		</div>
		<div id="mainContent" class="horizontalContent justifyAround alignCenter">
			<div id="createLobby">
				<button on:click={createGame}>CREATE GAME</button>
				<p>or</p>
				<input
					bind:value={lobbyCode}
					on:keyup={enterLobby}
					type="text"
					placeholder="Enter lobby code"
				/>
			</div>
			<div id="howtoplay">
				<h1>How to play</h1>
				<p>
					whoami is a two-player board game where players each guess the identity of the other
					chosen character
				</p>
				<ol>
					<li>Select character</li>
					<li>
						Ask yes or no questions about the other character. You can mark characters as "invalid"
						by left-clicking
					</li>
					<li>
						If you think you know the answer, right-click on the character and confirm the selection
					</li>
					<li>The first player to select the other's character wins</li>
				</ol>
			</div>
		</div>
	</div>
</div>

<style lang="scss">
	#root {
		width: 100%;
		height: 100vh;
	}
	#main {
		width: 100%;
		height: 100%;
	}
	#header {
		display: flex;
		width: 52%;
		justify-content: space-between;
		padding-top: 3rem;
		img {
			width: 30vw;
		}
	}
	#mainContent {
		height: 100%;
		align-items: center;
	}
	.justifyCenter {
		justify-content: center;
	}
	.justifyAround {
		justify-content: space-around;
	}
	.alignCenter {
		align-items: center;
	}
	.horizontalContent {
		display: flex;
	}
	.verticalContent {
		display: flex;
		flex-direction: column;
	}
	li {
		background-color: #1a1a1a;
		padding: 0.5rem;
		border-radius: 0.2rem;
		margin: 0.2rem;
		list-style-type: none;

		counter-increment: inst;

		&::before {
			content: counter(inst) '.';
			background-color: white;
			color: black;
			width: 1rem;
			height: 1rem;
			padding: 0.2rem;
			margin-right: 0.5rem;
			border-radius: 0.2rem;
			display: inline-block;
			text-align: center;
		}
	}
	#howtoplay {
		max-width: 50%;
	}
	#createLobby {
		background-image: linear-gradient(to left, #00bde7, #65d46e);
		border: none;
		border-radius: 0 2rem 0 2rem;
		padding: 3rem 1rem;
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;

		width: 40%;

		button {
			background-image: linear-gradient(to left, #be65d4, #d47965);
			border: none;
			border-radius: 2rem;
			padding: 0.5rem 1rem;
			font-weight: bold;

			&:hover {
				cursor: pointer;
			}
		}
		input {
			background-color: #1e1e1e;
			padding: 0.5rem 1rem;
			border: none;
			border-radius: 1rem;
		}
	}
</style>
