<script lang="ts">
	import { assets } from '$app/paths';
	import { tentativeCard, type Card } from '../stores';
	export let card: Card;
	$: selected = $tentativeCard?.name === card.name;

	let grayed = false;

	function handleClick(e: MouseEvent) {
		grayed = !grayed;
		selected = false;
	}

	function handleRightClick(e: MouseEvent) {
		selected = !selected;
		grayed = false;
		tentativeCard.set(selected ? card : ({} as Card));
	}
</script>

<div
	class:selected
	class:grayed
	id="characterCard"
	on:click={handleClick}
	on:contextmenu|preventDefault={handleRightClick}
>
	{#if selected}
		<div id="selectionCircle">
			<img alt="confirm" src="{assets}/checkmark.svg" />
		</div>
	{/if}
	<img alt="avatar" src={card.avatar} />
	<h2>{card.name}</h2>
</div>

<style lang="scss">
	#characterCard {
		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;

		background-color: #181818;
		border-radius: 0.5rem;

		position: relative;

		&:hover {
			background-color: #3b3b3b;
			cursor: pointer;
		}

		img {
			width: 7rem;
			height: 8rem;
		}

		h2 {
			margin: 0;
			font-family: 'Inter', sans-serif;
		}
	}

	.grayed {
		filter: grayscale(1);
		background-color: #3b3b3b !important;
	}

	.selected {
		position: relative;
		&::after {
			content: '';

			position: absolute;
			top: -0.2rem;
			bottom: -0.2rem;
			left: -0.2rem;
			right: -0.2rem;

			background-image: linear-gradient(to bottom left, #65d46e, #d46565);
			border-radius: 0.5rem;
			z-index: -1;
		}
	}

	#selectionCircle {
		background-color: #242424;
		border-radius: 100vw;

		width: 3rem;
		height: 3rem;

		position: absolute;
		top: -1.5rem;
		right: -1.5rem;
		display: flex;
		justify-content: center;
		align-items: center;

		img {
			width: 1rem;
			height: 1rem;
		}
	}
</style>
