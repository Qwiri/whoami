<script lang="ts">
	import { winnerName, winnerID, ingameName, cards } from '../stores';
	import Particles from 'svelte-particles';
	import { loadFull } from 'tsparticles';

	let particlesConfig = {
		emitters: [
			{
				position: {
					x: 0,
					y: 0
				},
				rate: {
					quantity: 50,
					delay: 0.15
				},
				particles: {
					move: {
						direction: 'top-right',
						outModes: {
							top: 'none',
							left: 'none',
							default: 'destroy'
						}
					}
				}
			},
			{
				position: {
					x: 100,
					y: 0
				},
				rate: {
					quantity: 5,
					delay: 0.15
				},
				particles: {
					move: {
						direction: 'top-left',
						outModes: {
							top: 'none',
							right: 'none',
							default: 'destroy'
						}
					}
				}
			}
		],
		particles: {
			color: {
				value: ['#ffffff', '#FF0000']
			},
			move: {
				decay: 0.05,
				direction: 'top',
				enable: true,
				gravity: {
					enable: true
				},
				outModes: {
					top: 'none',
					default: 'destroy'
				},
				speed: {
					min: 1,
					max: 10
				}
			},
			number: {
				value: 0
			},
			opacity: {
				value: 1
			},
			rotate: {
				value: {
					min: 0,
					max: 360
				},
				direction: 'random',
				animation: {
					enable: true,
					speed: 30
				}
			},
			tilt: {
				direction: 'random',
				enable: true,
				value: {
					min: 0,
					max: 360
				},
				animation: {
					enable: true,
					speed: 30
				}
			},
			size: {
				value: {
					min: 0,
					max: 2
				},
				animation: {
					enable: true,
					startValue: 'min',
					count: 1,
					speed: 16,
					sync: true
				}
			},
			roll: {
				darken: {
					enable: true,
					value: 25
				},
				enable: true,
				speed: {
					min: 5,
					max: 15
				}
			},
			wobble: {
				distance: 30,
				enable: true,
				speed: {
					min: -7,
					max: 7
				}
			},
			shape: {
				type: ['circle', 'square', 'triangle', 'polygon', 'character', 'image'],
				options: {
					polygon: [
						{
							sides: 5
						},
						{
							sides: 6
						}
					],
					character: {
						fill: true,
						font: 'Verdana',
						style: '',
						weight: 400,
						particles: {
							size: {
								value: 8
							}
						},
						value: ['💩', '🤡', '🍀', '🍙', '🦄', '⭐️']
					},
					image: [
						{
							src: 'https://particles.js.org/images/fruits/apple.png',
							width: 32,
							height: 32,
							particles: {
								size: {
									value: 16
								}
							}
						},
						{
							src: 'https://particles.js.org/images/fruits/avocado.png',
							width: 32,
							height: 32,
							particles: {
								size: {
									value: 16
								}
							}
						},
						{
							src: 'https://particles.js.org/images/fruits/banana.png',
							width: 32,
							height: 32,
							particles: {
								size: {
									value: 16
								}
							}
						},
						{
							src: 'https://particles.js.org/images/fruits/berries.png',
							width: 32,
							height: 32,
							particles: {
								size: {
									value: 16
								}
							}
						},
						{
							src: 'https://particles.js.org/images/fruits/cherry.png',
							width: 32,
							height: 32,
							particles: {
								size: {
									value: 16
								}
							}
						},
						{
							src: 'https://particles.js.org/images/fruits/grapes.png',
							width: 32,
							height: 32,
							particles: {
								size: {
									value: 16
								}
							}
						},
						{
							src: 'https://particles.js.org/images/fruits/lemon.png',
							width: 32,
							height: 32,
							particles: {
								size: {
									value: 16
								}
							}
						},
						{
							src: 'https://particles.js.org/images/fruits/orange.png',
							width: 32,
							height: 32,
							particles: {
								size: {
									value: 16
								}
							}
						},
						{
							src: 'https://particles.js.org/images/fruits/peach.png',
							width: 32,
							height: 32,
							particles: {
								size: {
									value: 16
								}
							}
						},
						{
							src: 'https://particles.js.org/images/fruits/pear.png',
							width: 32,
							height: 32,
							particles: {
								size: {
									value: 16
								}
							}
						},
						{
							src: 'https://particles.js.org/images/fruits/pepper.png',
							width: 32,
							height: 32,
							particles: {
								size: {
									value: 16
								}
							}
						},
						{
							src: 'https://particles.js.org/images/fruits/plum.png',
							width: 32,
							height: 32,
							particles: {
								size: {
									value: 16
								}
							}
						},
						{
							src: 'https://particles.js.org/images/fruits/star.png',
							width: 32,
							height: 32,
							particles: {
								size: {
									value: 16
								}
							}
						},
						{
							src: 'https://particles.js.org/images/fruits/strawberry.png',
							width: 32,
							height: 32,
							particles: {
								size: {
									value: 16
								}
							}
						},
						{
							src: 'https://particles.js.org/images/fruits/watermelon.png',
							width: 32,
							height: 32,
							particles: {
								size: {
									value: 16
								}
							}
						},
						{
							src: 'https://particles.js.org/images/fruits/watermelon_slice.png',
							width: 32,
							height: 32,
							particles: {
								size: {
									value: 16
								}
							}
						}
					]
				}
			}
		}
	};

	let onParticlesLoaded = (event) => {
		const particlesContainer = event.detail.particles;

		// you can use particlesContainer to call all the Container class
		// (from the core library) methods like play, pause, refresh, start, stop
	};

	let particlesInit = async (engine) => {
		// you can use main to customize the tsParticles instance adding presets or custom shapes
		// this loads the tsparticles package bundle, it's the easiest method for getting everything ready
		// starting from v2 you can add only the features you need reducing the bundle size
		await loadFull(engine);
	};

	let winnerCard = $cards[$winnerID];
</script>

<div class="content">
	<div id="tsparticles">
		<Particles
			id="tsparticles"
			options={particlesConfig}
			on:particlesLoaded={onParticlesLoaded}
			{particlesInit}
		/>
	</div>
	<div>
		<img src={winnerCard.avatar} alt="character of {$winnerName}" />
	</div>
	<div class="text">
		<h1><span class="green">{$winnerName === $ingameName ? 'You' : $winnerName}</span> won</h1>
		<p>as {winnerCard.name}</p>
	</div>
</div>

<style lang="scss">
	#tsparticles {
		position: absolute;
		z-index: -1;
	}
	.content {
		position: relative;
		height: 100%;
		width: 100%;

		display: flex;
		// justify-content: space-around;

		background-image: linear-gradient(to top, #6CD475, #d4ae6500);
		background-size: 100% 50%;
		background-position: 0 100%;
		background-repeat: no-repeat;

		.green {
			color: #65d46e;
		}
		.text {
			align-items: flex-start;
		}
		div {
			display: flex;
			width: 50%;
			flex-direction: column;
			justify-content: center;
			align-items: center;

			h1 {
				margin: 0;
				font-size: 5rem;
			}
			p {
				margin: 0;
				color: white;
				font-weight: 100;
			}
		}
	}
</style>
