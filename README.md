# LoL Cruncher

## What is LoL Cruncher?

LoL Cruncher _will_ be a web app that is an upgrade of [LoL Reader](http://1lann.github.io/lolreader), except instead of being powered by your logs, it will be powered by Riot Games API. It will monitor your stats and will update regularily without you having to take action, apart from having to check your account on the service once to activate it. However, if you haven't played a game for more than a month, monitoring will be stopped and will require you to visit the service to reactivate it. LoL Cruncher is designed to process lots of data over long periods of time. It will have the same features of LoL Reader and even more (except for time wasted on loading, lulz).

## Running

LoL Cruncher is written with [Revel](http://revel.github.io). Setup revel and clone this repository in your `$GOPATH/src` folder. This application does not come with its own Riot API key to make API requests. You will need to place the key in the conf/app.conf file (You can copy the one from app.conf.backup as the base) as the field of "riotapikey".

For example you would add this line somewhere near the top of the file (Obviously replace the Xs with your key):

`riotapikey = xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxx`

Then, run `revel run cruncher` and navigate your browser to http://127.0.0.1:9000

## The Plan/Todo List

- ~~Database structure~~
- ~~Core functionality without cron jobs~~
- ~~Design pages~~
- ~~Add cron jobs~~
- Apply for deployment key
- ???
- Profit


## License

LoL Cruncher is licensed under GNU Affero GPL v3.0. It can be found under /LICENSE and a summary and full license is available at http://choosealicense.com/licenses/agpl-3.0/. In a summary, it allows for modifications, commercial use and distribution as long as the source code is published, including server code if it is distributed on a network, and that I am not held liable for the code.

Note that this license does not apply to material used in LoL Cruncher I do not own. You can find their respective licenses at the top of their respective files.

This includes
- [Handlebars](https://github.com/wycats/handlebars.js/blob/master/LICENSE)
- [Bootstrap](https://github.com/twbs/bootstrap/blob/master/LICENSE)
- [jQuery](https://github.com/jquery/jquery/blob/master/LICENSE.txt)
- [Semantic UI](https://github.com/Semantic-Org/Semantic-UI/blob/master/LICENSE.md)
- [Revel](https://github.com/revel/revel/blob/master/LICENSE)

### Why did I choose GNU Affero?

Because I want people who make forks, improvements or distribute my work, also publish their source code for others to take and learn from, to gain inspiration and knowledge.

