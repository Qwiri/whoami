<script lang="ts" context="module">
	export interface PeekContent {
		ID: string;
		State: number;
		Clients: string[];
	}
</script>

<script lang="ts">
	export let peek: PeekContent;

	let danger = false;
	$: {
        peek.Clients.sort();
		danger = peek.State != 1 << 0;
	}
</script>

<div id="peek">
	<span id="title">LOBBY AVAILABLE</span> <span id="state">STATE {peek.State}</span>

	{#if danger}
		<div id="warn">
			<span class="danger">⚠️ This game is already in progress</span>
		</div>
	{/if}

	<div id="connected">
		<span><span class="green">{peek.Clients.length}</span> player/s connected:</span>
		<ul>
			{#each peek.Clients as client}
				<li>{client}</li>
			{/each}
		</ul>
	</div>
</div>

<style lang="scss">
	#peek {
        #warn {
            margin: 1rem 0rem 1rem 0rem;
        }
		background-color: #1b1b1b;
		padding: 1rem;
		border-radius: 0.3rem;

		border-bottom: 2px solid #49beff;

		.danger {
			background-color: #ff6f6f;
            padding: .7rem;
            border-radius: 0.3rem;
		}

		span {
			font-size: 0.7rem;
		}

		.green {
			color: #65d46e;
		}

		#title {
			font-size: 0.7rem;
			font-weight: bold;
			color: #65d46e;
		}

		#state {
			font-size: 0.3rem;
			color: #49beff;
		}

        ul {
            margin: 0;
            font-size: 0.7rem;

            li {
                color: #787878;
            }
        }
	}
</style>
