<script lang="ts">
	import CharacterCard from './CharacterCard.svelte';

	import { cards, type Card } from '../stores';

	export let onLeftClick: undefined | ((e: MouseEvent, player: Card, i: number) => void);
	export let onRightClick: undefined | ((e: MouseEvent, player: Card, i: number) => void);

	// function removeSelection() {
	// 	cards.filter((p) => p.selected).forEach((p) => (p.selected = false));
	// }
</script>

<div id="characterDecision">
	{#each $cards as card, i}
		<div
			on:contextmenu|capture={(e) => {
				if (onRightClick) {
					onRightClick(e, card, i);
				}
			}}
			on:click|capture={(e) => {
				if (onLeftClick) {
					onLeftClick(e, card, i);
				}
			}}
		>
			<CharacterCard {card} />
		</div>
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
