<script lang="ts">
	import Navbar from '../Components/Navbar.svelte';

	let lobbyCode: string;

	function enterLobby(e: KeyboardEvent) {
		if (e.key === 'Enter') {
			window.location.href = `/g/${lobbyCode}`;
		}
	}
	async function createGame() {
		const res = await fetch('https://backend.wai.d2a.io/lobby/create');
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
			<h1>WHO<span class="light">AM</span>I</h1>
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
			<div>
				<svg id="separator" viewBox="0 0 2 136" fill="none" xmlns="http://www.w3.org/2000/svg">
					<path d="M1 0L1.00001 136" stroke="#C0C0C0" stroke-width="2" />
				</svg>
			</div>
			<div id="howtoplay">
				<h3>HOW TO PLAY</h3>
				<p>
					<span class="green">whoami</span> is a two-player board game where players each guess the identity
					of the other chosen character
				</p>
				<ol>
					<li>Select character</li>
					<li>
						Ask <span class="green">yes</span> or <span class="green">no</span> questions about the
						other character. You can mark characters as "invalid" by
						<span class="blue">left-clicking</span>
					</li>
					<li>
						If you think you know the answer, <span class="blue">right-click</span> on the character
						and confirm the selection
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
	.green {
		color: #65d46e;
	}
	.blue {
		color: #49beff;
	}
	.light {
		font-weight: 200;
	}
	#separator {
		width: 2px;
		height: 100%;
	}
	ol {
		padding: 0;
		margin: 0;
	}
	li {
		background-color: #1a1a1a;
		padding: 0.6rem;
		border-radius: 0.4rem;
		margin: 0.6rem;
		margin-left: 0;
		margin-right: 0;
		list-style-type: none;

		counter-increment: inst;

		&::before {
			content: counter(inst) '.';
			background-color: white;
			color: black;
			width: 1rem;
			height: 1rem;
			padding: 0.3rem;
			margin-right: 0.5rem;
			border-radius: 0.2rem;
			display: inline-block;
			text-align: center;
			font-weight: bold;
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
