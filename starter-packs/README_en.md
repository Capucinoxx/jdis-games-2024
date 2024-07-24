# Magellan - JDIS Games 2024

## Scenario

Imagine yourself in the early 16th century, a time when world maps were still filled with unknown and mysterious areas. During this period of discoveries and explorations, Ferdinand Magellan, a Portuguese explorer, was preparing. His goal: to find a sea route to the spice islands, whose riches were coveted by all of Europe.

On board his ship, Magellan and his crew ventured into uncharted waters, armed with patience and their knowledge of navigation. Their main ally was the astrolabe, an instrument capable of measuring the altitude of stars and planets. This tool was essential for charting their course through vast oceans, allowing them to find their way to the unknown.

## Your Objective

Like Magellan, you will become explorers navigating uncharted waters in search of treasure. You know its coordinates, but the location of obstacles remains unclear. There will be two phases in your adventure: discovery and acquisition.

## Map Elements

During previous explorations, some information was gathered.

The following elements will be found on the map:

|        |                                             |
| ------ | ------------------------------------------- |
| Coin   | ![image](/starter-packs/docs/coin.png)      |
| Treasure| ![image](/starter-packs/docs/astrolab2.png)|
| Wall   | ...                                         |

### Discrete Grid

The discrete grid contains the number of walls per 4 squares. The outer walls delineating the map are also counted.

<div align="center">
  <img width="1000" alt="logo" src="./docs/grille_murs.png">
</div>

### Scoring

| Item       | Points |
| ---------- | ------ |
| Pistol     | 15     |
| Sword      | 40     |
| Coin       | 40     |
| Treasure   | 1200   |

## Game Flow

In a game, there are two phases.

### Phase 1

Duration: x time
Agents and coins are placed randomly on the map 🗺️

When an agent picks up a coin, it reappears randomly on the map.

At the end of the phase, agents are removed from the map.

### Phase 2

Duration: x time
A treasure is placed in the center of the map, and agents will appear at an equidistant distance from the treasure. In this phase, there will be no coins on the map.

### Death

When the agent loses all their life, they disappear from the map and receive no data from the server for a defined period.

## Actions

Several actions can be sent to the server in the same refresh cycle with certain usage constraints.

### Movement

The position sent is the one the agent will move towards. The agent cannot pass through walls. No pathfinding algorithm is implemented; you will need to implement it to navigate the map. Note that there are no collisions between players.

This action has no usage constraints.

### Attack

#### Equipment

The agent must be equipped with a weapon to use it.

This action cannot be accompanied by the use of a weapon in the same refresh cycle.

#### Pistol

To use the pistol, send the desired destination position of a projectile. 
This action cannot be accompanied by equipping a weapon in the same refresh cycle. 

#### Sword

When the sword is equipped, it appears at 0 radians as illustrated in the following image.  

<div align="center">
  <img width="200" alt="logo" src="./docs/blade2.png">
</div>

To change the rotation of the sword, send the new rotation in radians.

### Save

A limited amount of bytes can be sent to the server. These bytes will be saved for the duration of a game. These data will be received by the player each time they connect to the server.

## Received Data

### Map State

The map state is received when connecting to the server. 

