<script lang="ts">
	import { currentLives, maxLives, selectedCard } from '../stores';

	export let ingame = false;
	let startTime = Date.now();
	let ingameTime = '00:00';

	setInterval(() => {
		let seconds = Math.floor((Date.now() - startTime) / 1000);
		let minutes = Math.floor(seconds / 60);
		seconds = Math.floor(seconds % 60);

		ingameTime = `${minutes.toLocaleString('en-US', {
			minimumIntegerDigits: 2
		})}:${seconds.toLocaleString('en-US', { minimumIntegerDigits: 2 })}`;
	}, 1000);
</script>

<div id="nav">
	<img id="wai" src="/WAI.svg" alt="wai icon" />
	{#if ingame}
		<div class="hearts">
			{#each Array($maxLives - $currentLives) as item}
				<img src="/broken_heart.png" alt="lost life" />
			{/each}
			{#each Array($currentLives) as item}
				<img src="/heart.png" alt="one life" />
			{/each}
		</div>
		<img alt="avatar" src={$selectedCard?.avatar} />
		<hr id="hr" />
		<p id="time">{ingameTime}</p>
		<img on:click={() => (window.location.href = '/')} id="leave" alt="leave" src="/leave.svg" />
	{/if}
</div>

<style lang="scss">
	#nav {
		width: 3rem;
		padding: 1rem;
		background-color: black;
		align-self: stretch;

		display: flex;
		flex-direction: column;
		justify-content: flex-end;
	}
	.hearts {
		img {
			width: 3rem;
		}
	}
	#hr {
		width: 90%;
		height: 2px;
		background-color: #202020;
		border: none;
	}
	#wai {
		margin-bottom: auto;
	}
	#time {
		font-weight: bold;
	}
	#leave {
		background-color: #333333;
		padding: 0.5rem;
		border-radius: 0.5rem;

		&:hover {
			cursor: pointer;
		}
	}
</style>
