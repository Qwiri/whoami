<script>
	import CharacterDecision from '../Components/CharacterDecision.svelte';
</script>

<div id="container">
	<h1>WhoAmI</h1>
	<CharacterDecision />
</div>

<style lang="scss">
	#container {
		width: 100%;
		height: 100vh;

		display: flex;
		flex-direction: column;
		justify-content: center;
		align-items: center;
	}
</style>
