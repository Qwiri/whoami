<script lang="ts">
	import CharacterCard from './CharacterCard.svelte';

	let players = [
		{ name: 'Luca', selected: false },
		{ name: 'Tobus', selected: false },
		{ name: 'Darmiel', selected: false },
		{ name: 'Jimmy', selected: false },
		{ name: 'JuJu', selected: false },
		{ name: 'Bevin', selected: false },
		{ name: 'Bernd', selected: false },
		{ name: 'Janni boi', selected: false },
		{ name: 'JJ', selected: false },
		{ name: 'Anni', selected: false },
		{ name: 'Alexa', selected: false },
		{ name: 'Pandora', selected: false },
		{ name: 'Yo Mamer', selected: false },
		{ name: 'Dora', selected: false },
		{ name: 'Chad', selected: false },
		{ name: 'Buchard', selected: false },
		{ name: 'Sebastian', selected: false },
		{ name: 'Karlheinz', selected: false },
		{ name: 'Thea', selected: false },
		{ name: 'Linda', selected: false },
		{ name: 'Josephine', selected: false },
		{ name: 'Aaron', selected: false },
		{ name: 'Theo', selected: false },
		{ name: 'Markus', selected: false }
	];

	function removeSelection() {
		players.filter((p) => p.selected).forEach((p) => (p.selected = false));
	}
</script>

<div id="characterDecision" on:contextmenu|capture={removeSelection}>
	{#each players as player}
		<CharacterCard bind:selected={player.selected} name={player.name} />
	{/each}
</div>

<style lang="scss">
	#characterDecision {
		display: grid;
		grid-template-columns: repeat(6, 1fr);
		grid-template-rows: repeat(4, 1fr);
		justify-items: center;
		align-items: center;

		gap: 1rem;
		direction: rtl;
	}
</style>
