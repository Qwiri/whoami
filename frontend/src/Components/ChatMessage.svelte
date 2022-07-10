<script lang="ts">
	export let me: Boolean;
	export let text: string;
</script>

<p id="message" class={me ? 'me' : 'notMe'}>{text}</p>

<style lang="scss">
	#message {
		border-radius: 2rem;
		padding: 0.5rem 1rem;
		overflow-wrap: anywhere;
		max-width: 80%;

		margin: 0.2rem;
	}
	.me {
		background-color: #181818;
		color: #6f6f6f;
		align-self: flex-end;
	}
	.notMe {
		background-color: #65d46e;
		align-self: flex-start;
	}
</style>
