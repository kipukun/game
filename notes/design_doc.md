# design document

## high level

### working title
J Walking

### concept statement
J Walking is a party game that mixes traditional JRPG combat with
the dynamics of turn-based party games.

### genre
party, JRPG

### target audience
fans of games like mario marty, and fans of games like the final fantasy series.

### unique selling points
J Walking will have:
- familiar party game turn based combat
- but with interesting combat that changes with every game, including chance and strategy

## product design
### player experience
the player will be one of N classes (rogue, mage, warrior, ...), inside of a typical 
fantasy world that would include different vistas. the player should feel satsified 
when their strategy comes together, but should still feel worried that everything
might not go as expected. random chance and the discovery of new items will keep players
on their toes and wanting to discover more.

### visual and audio style
the look and feel is a pixel 2D top-down with themes borrowed from games like classic JRPGs.
this will keep the focus on gameplay in a familiar setting.
[insert images of classic JRPGs here]

### game world fiction
you are a champion for a kingdom in a world split into different factions. you are sent to
the Happenings against other champions
from neighboring factions in order to fight for your kingdom and bring back glory.
this is similar to the olympics, with each country hosting their own version
of the Happenings.


### monetization
probably free

### tech
PC, 2D, with my own engine. this is with a one man team, so god knows how long 
until its playable.


## game system design
### core loops
the core loop is a turn, where the player rolls a random number and proceeds
the amount of spaces. depending on the space object, the player is thrust into a 
different situation. they have to use the tools they acquired to deal with it, or be 
punished. emergent results i hope to see:
- different builds depending on the players preference
- planning

### game systems
these are the systems needed:
- a progression tracker (keep track of player progress throughout a game)
- 2D object system (does not require collision/physics)
- multiplayer system
- item system
- battle system
- mini game system