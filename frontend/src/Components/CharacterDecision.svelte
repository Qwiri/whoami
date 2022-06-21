<script lang="ts">
	import CharacterCard from './CharacterCard.svelte';
</script>

<div id="characterDecision">
	{#each Array(24) as _}
		<CharacterCard name="luca" />
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
