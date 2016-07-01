# Codenames board generator

This week I made a little app for playing the board game [Codenames](https://boardgamegeek.com/boardgame/178900/codenames). It solves a handful of small pain points:

* I need the physical game on hand to play.
* Using the spymaster's key to look at your words gets annoying.
* Word cards are double-sided, so there are pairs of words that you never get to play together.
* Word cards are small and with a large group, it's hard to make sure everyone can see the board.

To play using this app, go to [www.horsepaste.com](http://www.horsepaste.com), enter a name and click 'New.' The app will create a random board for you. If you copy & paste the URL, other devices can connect to the same game.

![codenames-spymaster-view](/static/images/codenames-spymaster-view.png)
<p class="caption">An in-progress game from the spymaster view. Words with colored backgrounds have been revealed. The darker gray tile is the Assassin.</p>

Each device can toggle between spymaster and normal player views by the buttons in the bottom right corner. Players can click on words to reveal them. The app will keep track of whose turn it is and announce when a team wins.

If you're interested or want to contribute, the source is on github at [github.com/jbowens/codenames](https://github.com/jbowens/codenames). The server is in Go, and the frontend uses React.js.
