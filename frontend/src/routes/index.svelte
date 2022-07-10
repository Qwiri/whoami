<script lang="ts">
	import Chat from '../Components/Chat.svelte';
	import CharacterDecision from '../Components/CharacterDecision.svelte';
</script>

<div id="container">
	<div id="nav">
		<img src="/WAI.svg" alt="wai icon" />
	</div>
	<div id="content">
		<CharacterDecision />
		<Chat />
	</div>
</div>

<style lang="scss">
	#container {
		width: 100%;
		height: 100vh;

		display: flex;
		// justify-content: space-between;
		align-items: center;
	}
	#content {
		display: flex;
		width: 100%;
		justify-content: center;
		gap: 4rem;
	}
	#nav {
		width: 3rem;
		padding: 1rem;
		background-color: black;
		height: 100%;
	}
</style>
