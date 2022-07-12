<script lang="ts">
	import ChatMessage from './ChatMessage.svelte';

	import { messages, ingameName } from '../stores';

	// let messages = [{ me: false, text: 'hey, nice to meet you! :)' }];
	export let sendMessageCallback: (text: string) => void;

	let text: string;

	function sendMessage(e: KeyboardEvent) {
		if (e.key === 'Enter' && text) {
			sendMessageCallback(text);
			text = '';
		}
	}
</script>

<div id="chatContainer">
	<div id="messageContainer">
		{#each $messages as message}
			<ChatMessage me={message.user === $ingameName} text={message.message} />
		{/each}
	</div>
	<hr id="hr" />
	<input
		bind:value={text}
		id="textInput"
		type="text"
		on:keyup={sendMessage}
		placeholder="type message"
	/>
</div>

<style lang="scss">
	#messageContainer {
		display: flex;
		flex-direction: column-reverse;
		max-height: 75vh;
		overflow-y: scroll;

		&::-webkit-scrollbar {
			width: 12px;
			overflow-y: hidden;
		}
		&::-webkit-scrollbar-track {
			background-color: transparent;
		}
		&::-webkit-scrollbar-thumb {
			background-color: #d6dee1;
			border-radius: 20px;
			border: 4px solid transparent;
			background-clip: content-box;
		}
		&::-webkit-scrollbar-thumb:hover {
			background-color: #a8bbbf;
			border: 2px solid transparent;
			background-clip: content-box;
		}
	}
	#chatContainer {
		background-color: #272727;
		border-radius: 0.5rem;
		padding: 1rem;
		width: 20vw;
		display: flex;
		flex-direction: column;
		justify-content: flex-end;
	}
	#hr {
		border: none;
		height: 1px;
		width: 90%;
		background-color: #3c3c3c;
	}
	#textInput {
		background-color: #1e1e1e;
		border-radius: 0.2rem;
		border: none;
		padding: 0.5rem;
	}
</style>
