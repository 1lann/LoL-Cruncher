
// LoL Cruncher - A Historical League of Legends Statistics Tracker
// Copyright (C) 2015  Jason Chu (1lann) 1lanncontact@gmail.com

// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.

// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.

// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

// jscs: disable


var statsSource = '\
<table class="stats-table">\
<tbody>\
	{{#each statsRow}}\
	<tr>\
		<td class="stats-data">{{data}}</td>\
		<td class="stats-label">{{label}}</td>\
	</tr>\
	{{/each}}\
</tbody>\
</table>'

var statsTemplate = Handlebars.compile(statsSource);

var timePlayingSource = '\
<span class="time-playing">{{time}}</span>\
<span class="spent-playing"> spent playing</span>\
'

var timePlayingTemplate = Handlebars.compile(timePlayingSource);

var filtersAreaSource = '\
<span>\
	<div class="ui inline dropdown" id="general-dropdown">\
		<input type="hidden" name="data-type">\
		<div class="text">All</div>\
			<i class="dropdown icon"></i>\
			<div class="menu">\
			<div class="item active selected" data-text="All">All</div>\
			<div class="item" data-text="Rates/average">Rates/Averages</div>\
		</div>\
	</div>\
	statistics\
	<div class="ui inline dropdown" id="date-dropdown">\
		<input type="hidden" name="timeframe">\
		<div class="text">since {{start}}</div>\
			<i class="dropdown icon"></i>\
			<div class="menu">\
			<div class="item active selected" data-text="since {{start}}">since {{start}}</div>\
			{{#each month}}\
			<div class="item" data-text="{{text}}">{{text}}</div>\
			{{/each}}\
		</div>\
	</div>\
	for\
	<div class="ui inline dropdown" id="queue-dropdown">\
		<input type="hidden" name="queue-type">\
		<div class="text">all queues</div>\
			<i class="dropdown icon"></i>\
			<div class="menu">\
			<div class="item active selected" data-text="all queues">all queues</div>\
			{{#each queueTypes}}\
			<div class="item" data-text="{{name}}">{{name}}</div>\
			{{/each}}\
		</div>\
	</div>\
	<span class="glyphicon glyphicon-info-sign" data-container="body" data-toggle="popover"\
	 data-trigger="hover" data-placement="top" data-content="Customs, Dominion and featured gamemodes\'\
	 statistics will not be recorded."></span>\
</span>\
<br>\
<span>\
	Display statistics as\
	<div class="ui inline dropdown" id="display-dropdown">\
		<input type="hidden" name="display-type">\
		<div class="text">cards</div>\
			<i class="dropdown icon"></i>\
			<div class="menu">\
			<div class="item active selected" data-text="cards">cards</div>\
			<div class="item" data-text="a table">a table</div>\
		</div>\
	</div>\
</span>\
<p class="delay">Statistics may be delayed by up to 24 hours</p>'

var filtersAreaTemplate = Handlebars.compile(filtersAreaSource);

var championCardSource = '\
<div class="stats-card">\
	<img class="champion-image" src="//ddragon.leagueoflegends.com/cdn/6.11.1/img/champion/{{imageName}}">\
	<div class="champion-label">\
		<p class="title">Stats for {{displayName}}</p>\
		<p class="subheading">{{dateFilter}}</p>\
	</div>\
	<div class="stats-area">\
	{{{stats}}}\
	</div>\
</div>\
'

var championCardTemplate = Handlebars.compile(championCardSource);

var generalCardSource = '\
<div class="stats-card">\
	<div class="title">All champions stats</div>\
	<p class="subheading">{{dateFilter}}</p>\
	<div class="stats-area">\
	{{{stats}}}\
	</div>\
</div>'

var generalCardTemplate = Handlebars.compile(generalCardSource);

var profileSource = '\
<img src="//avatar.leagueoflegends.com/{{regionCode}}/{{{imageName}}}.png">\
<div class="profile-info">\
	<p class="username">{{username}}</p>\
	<p class="region">{{region}}</p>\
</div>'

var profileTemplate = Handlebars.compile(profileSource);

var tableAllHeader = '\
<th class="no-pointer"></th> <!-- Image -->\
<th>Champion</th>\
<th>Played</th>\
<th class="won">Won</th>\
<th class="lost">Lost</th>\
<th class="minions">Minions</th>\
<th class="jungle">Jungle</th>\
<th class="gold">Gold</th>\
<th class="wards">Wards</th>\
<th class="kills">Kills</th>\
<th class="deaths">Deaths</th>\
<th class="assists">Assists</th>'

var tableRatesHeader = '\
<th class="no-pointer"></th> <!-- Image -->\
<th>Champion</th>\
<th>Played</th>\
<th class="won">Winrate</th>\
<th class="minions">Minions/10m</th>\
<th class="jungle">Jungle/10m</th>\
<th class="gold">Gold/10m</th>\
<th class="wards">Wards</th>\
<th class="kills">Kills</th>\
<th class="deaths">Deaths</th>\
<th class="assists">Assists</th>'

var tableStatsSource = '\
<table class="table" id="stats-table">\
<thead>\
	<tr>\
		{{{tableHeader}}}\
	</tr>\
</thead>\
<tbody>\
	{{#each championStats}}\
	<tr>\
	{{{stats}}}\
	</tr>\
	{{/each}}\
</tbody>\
{{{tableFooter}}}\
</table>'


var tableStatsTemplate = Handlebars.compile(tableStatsSource);

var tableRowSource = '\
<td><img src="//ddragon.leagueoflegends.com/cdn/6.11.1/img/champion/{{imageName}}"></td>\
<td>{{championName}}</td>\
<td>{{games}}</td>\
<td class="won">{{wins}}</td>\
{{#if losses}}\
<td class="lost">{{losses}}</td>\
{{/if}}\
<td class="minions">{{minions}}</td>\
<td class="jungle">{{jungle}}</td>\
<td class="gold">{{gold}}</td>\
<td class="wards">{{wards}}</td>\
<td class="kills">{{kills}}</td>\
<td class="deaths">{{deaths}}</td>\
<td class="assists">{{assists}}</td>'

var tableRowTemplate = Handlebars.compile(tableRowSource);

tableFooterSource = '\
<tfoot>\
<tr>\
<td></td>\
<td>Total</td>\
<td>{{games}}</td>\
<td class="won">{{wins}}</td>\
{{#if losses}}\
<td class="lost">{{losses}}</td>\
{{/if}}\
<td class="minions">{{minions}}</td>\
<td class="jungle">{{jungle}}</td>\
<td class="gold">{{gold}}</td>\
<td class="wards">{{wards}}</td>\
<td class="kills">{{kills}}</td>\
<td class="deaths">{{deaths}}</td>\
<td class="assists">{{assists}}</td>\
</tr>\
</tfoot>\
'

var tableFooterTemplate = Handlebars.compile(tableFooterSource);
