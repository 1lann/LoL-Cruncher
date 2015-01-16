
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
		<div class="text">All your</div>\
			<i class="dropdown icon"></i>\
			<div class="menu">\
			<div class="item active selected" data-text="All your">All</div>\
			<div class="item" data-text="Your rates/average">Rates/Averages</div>\
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
</span>\
<p class="delay">Statistics may be delayed by up to 24 hours</p>'

var filtersAreaTemplate = Handlebars.compile(filtersAreaSource)

var leftDropdownSource ='\
<div class="ui inline dropdown">\
	<input type="hidden" name="data-type">\
	<div class="text title">All stats</div>\
		<i class="dropdown icon"></i>\
		<div class="menu">\
		<div class="item active selected" data-text="All stats">All stats</div>\
		{{#each gametype}}\
		<div class="item" data-text="{{text}}">{{text}}</div>\
		{{/each}}\
	</div>\
</div>'

var leftDropdownTemplate = Handlebars.compile(leftDropdownSource)

var championCardSource = '\
{{#each championCard}}\
<div class="stats-card">\
	<img class="champion-image" src="http://ddragon.leagueoflegends.com/cdn/5.1.1/img/champion/{{imageName}}.png">\
	<div class="champion-label">\
		<p class="title">Stats for {{displayName}}</p>\
		<p class="subheading">{{dateFilter}}</p>\
	</div>\
	<div class="stats-area">\
	{{{stats}}}\
	</div>\
</div>\
{{/each}}\
'

var championCardTemplate = Handlebars.compile(championCardSource)

var generalCardSource = '\
<div class="stats-card">\
	{{{dropdown}}}\
	<p class="subheading">{{dateFilter}}</p>\
	<div class="stats-area">\
	{{{stats}}}\
	</div>\
</div>'

var generalCardTemplate = Handlebars.compile(generalCardSource)
