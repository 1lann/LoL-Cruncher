
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

$(".ui.dropdown").dropdown();

var selectedDate = "since " + playerData.RecordStart;
var selectedFilter = "All"

var templateMonths = [];
var dataMonths = {};
dataMonths[selectedDate] = playerData.All;

var selectedGameType = "All stats"

var championsSearchIndex = {};
var championFilter = ""

var selectedCollection = dataMonths[selectedDate];

var monthsKeys = {
	"1": "January",
	"2": "February",
	"3": "March",
	"4": "April",
	"5": "May",
	"6": "June",
	"7": "July",
	"8": "August",
	"9": "September",
	"10": "October",
	"11": "November",
	"12": "December",
}

var dateRegex = /(\d+)\s(\d+)/;

var getStringedDate = function(date) {
	var resp = date.match(dateRegex)
	return monthsKeys[resp[2]] + " " + resp[1]
}

var loadMonths = function() {
	var months = Object.keys(playerData.MonthlyStats)
	months.sort();
	for (var i = 0; i < months.length; i++) {
		var monthKey = "for " + getStringedDate(months[i]);
		dataMonths[monthKey] = playerData.MonthlyStats[months[i]];
		templateMonths.push({text: monthKey});
	}
}

var getHumanTime = function(seconds) {
	var numdays = Math.floor(seconds / 86400);
	var numhours = Math.floor((seconds % 86400) / 3600);
	var numminutes = Math.floor(((seconds % 86400) % 3600) / 60);
	var construct = "";

	if (numdays == 1) {
		construct = "1 day";
	} else if (numdays > 1) {
		construct = numdays + " days";
	}

	if (numhours == 1) {
		if (construct != "") {
			construct = construct + ", ";
		}
		construct = construct + "1 hour";
	} else if (numhours > 1) {
		if (construct != "") {
			construct = construct + ", ";
		}
		construct = construct + numhours + " hours";
	}

	if (construct != "") {
		construct = construct + ", and ";
	}
	if (numminutes == 1) {
		construct = construct + "1 minute";
	} else {
		construct = construct + numminutes + " minutes";
	}
	return construct;
}

var getGoldAmount = function(gold) {
	if (gold > 999999) {
		var prefix = Math.round(gold/10000) / 100;
		return prefix + "m"
	} else {
		var prefix = Math.round(gold/100) / 10;
		return prefix + "k"
	}
}

var oneDecRound = function(num) {
	return (Math.round(num * 10) / 10).toString();
}

// You may want to collapse these functions in your IDE/Text Editor
var generateGeneralStats = function() {
	var outputStats = [];
	var statsSource = selectedCollection.All

	if (selectedGameType != "All stats") {
		statsSource = selectedCollection.GameTypeStats[selectedGameType];
	}

	var spentPlaying = timePlayingTemplate({
		time: getHumanTime(statsSource.TimePlayed)
	});

	if (selectedFilter == "All") {
		outputStats.push({
			label: "Games played",
			data: (statsSource.Wins + statsSource.Losses).toString(),
		});
		outputStats.push({
			label: "Games won",
			data: statsSource.Wins.toString(),
		});
		outputStats.push({
			label: "Games lost",
			data: statsSource.Losses.toString(),
		});
		outputStats.push({
			label: "Games played on red",
			data: (statsSource.Red.Wins + statsSource.Red.Losses).toString(),
		});
		outputStats.push({
			label: "Games played on blue",
			data: (statsSource.Blue.Wins + statsSource.Blue.Losses).toString(),
		});
		outputStats.push({
			label: "Minions killed",
			data: statsSource.MinionsKilled.toString(),
		});
		outputStats.push({
			label: "Jungle monsters killed",
			data: statsSource.MonstersKilled.toString(),
		});
		outputStats.push({
			label: "Gold earned",
			data: getGoldAmount(statsSource.GoldEarned),
		});
		outputStats.push({
			label: "Wards placed",
			data: statsSource.WardsPlaced.toString(),
		});
		outputStats.push({
			label: "Wards killed",
			data: statsSource.WardsKilled.toString(),
		});
		outputStats.push({
			label: "Kills",
			data: statsSource.Kills.toString(),
		});
		outputStats.push({
			label: "Deaths",
			data: statsSource.Deaths.toString(),
		});
		outputStats.push({
			label: "Assists",
			data: statsSource.Assists.toString(),
		});
		outputStats.push({
			label: "Double kills",
			data: statsSource.DoubleKills.toString(),
		});
		outputStats.push({
			label: "Triple kills",
			data: statsSource.TripleKills.toString(),
		});
		outputStats.push({
			label: "Quadra kills",
			data: statsSource.QuadraKills.toString(),
		});
		outputStats.push({
			label: "Pentakills",
			data: statsSource.PentaKills.toString(),
		});
	} else if (selectedFilter == "Rates/average") {
		var numGames = statsSource.Wins + statsSource.Losses;
		var timePlayed = statsSource.TimePlayed;
		outputStats.push({
			label: "Games played",
			data: numGames.toString(),
		});
		outputStats.push({
			label: "Average game time in minutes",
			data: Math.round(timePlayed/numGames/60),
		});
		outputStats.push({
			label: "Winrate",
			data: Math.round((statsSource.Wins/numGames) * 100) + "%",
		});
		var redGames = statsSource.Red.Wins + statsSource.Red.Losses
		outputStats.push({
			label: "Red team winrate",
			data: Math.round((statsSource.Red.Wins/redGames) * 100) + "%",
		});
		var blueGames = statsSource.Blue.Wins + statsSource.Blue.Losses
		outputStats.push({
			label: "Blue team winrate",
			data: Math.round((statsSource.Blue.Wins/blueGames) * 100) + "%",
		});
		outputStats.push({
			label: "Minions killed per 10 minutes",
			data: oneDecRound(statsSource.MinionsKilled/(timePlayed/600)),
		});
		outputStats.push({
			label: "Jungle monsters killed per 10 minutes",
			data: oneDecRound(statsSource.MonstersKilled/(timePlayed/600)),
		});
		outputStats.push({
			label: "Gold earned per 10 minutes",
			data: getGoldAmount(statsSource.GoldEarned/(timePlayed/600)),
		});
		outputStats.push({
			label: "Wards placed per game",
			data: oneDecRound(statsSource.WardsPlaced/numGames),
		});
		outputStats.push({
			label: "Wards killed per game",
			data: oneDecRound(statsSource.WardsKilled/numGames),
		});
		outputStats.push({
			label: "Kills per game",
			data: oneDecRound(statsSource.Kills/numGames),
		});
		outputStats.push({
			label: "Deaths per game",
			data: oneDecRound(statsSource.Deaths/numGames),
		});
		outputStats.push({
			label: "Assists per game",
			data: oneDecRound(statsSource.Assists/numGames),
		});
	}

	return spentPlaying + statsTemplate({statsRow: outputStats});
}

var generateChampionStats = function(championId) {
	var outputStats = [];
	var statsSource = selectedCollection.Champions[championId];

	var spentPlaying = timePlayingTemplate({
		time: getHumanTime(statsSource.TimePlayed)
	});

	if (selectedFilter == "All") {
		outputStats.push({
			label: "Games played",
			data: (statsSource.Wins + statsSource.Losses).toString(),
		});
		outputStats.push({
			label: "Games won",
			data: statsSource.Wins.toString(),
		});
		outputStats.push({
			label: "Games lost",
			data: statsSource.Losses.toString(),
		});
		outputStats.push({
			label: "Minions killed",
			data: statsSource.MinionsKilled.toString(),
		});
		outputStats.push({
			label: "Jungle monsters killed",
			data: statsSource.MonstersKilled.toString(),
		});
		outputStats.push({
			label: "Wards placed",
			data: statsSource.WardsPlaced.toString(),
		});
		outputStats.push({
			label: "Kills",
			data: statsSource.Kills.toString(),
		});
		outputStats.push({
			label: "Deaths",
			data: statsSource.Deaths.toString(),
		});
		outputStats.push({
			label: "Assists",
			data: statsSource.Assists.toString(),
		});
	} else if (selectedFilter == "Rates/average") {
		var numGames = statsSource.Wins + statsSource.Losses;
		var timePlayed = statsSource.TimePlayed;
		outputStats.push({
			label: "Games played",
			data: numGames.toString(),
		});
		outputStats.push({
			label: "Average game time in minutes",
			data: Math.round(timePlayed/numGames/60),
		});
		outputStats.push({
			label: "Winrate",
			data: Math.round((statsSource.Wins/numGames) * 100) + "%",
		});
		outputStats.push({
			label: "Minions killed per 10 minutes",
			data: oneDecRound(statsSource.MinionsKilled/(timePlayed/600)),
		});
		outputStats.push({
			label: "Jungle monsters killed per 10 minutes",
			data: oneDecRound(statsSource.MonstersKilled/(timePlayed/600)),
		});
		outputStats.push({
			label: "Wards placed per game",
			data: oneDecRound(statsSource.WardsPlaced/numGames),
		});
		outputStats.push({
			label: "Kills per game",
			data: oneDecRound(statsSource.Kills/numGames),
		});
		outputStats.push({
			label: "Deaths per game",
			data: oneDecRound(statsSource.Deaths/numGames),
		});
		outputStats.push({
			label: "Assists per game",
			data: oneDecRound(statsSource.Assists/numGames),
		});
	}

	return spentPlaying + statsTemplate({statsRow: outputStats});
}

var generateGeneralArea = function() {
	$(".general-card").empty();

	var gameTypes = [];
	for (var gameType in selectedCollection.GameTypeStats) {
		gameTypes.push({text: gameType});
	}

	var leftDropdown = leftDropdownTemplate({
		gametype: gameTypes,
	});

	var statsArea = generateGeneralStats();

	var generalArea = $(generalCardTemplate({
		dateFilter: selectedDate,
		dropdown: leftDropdown,
		stats: statsArea,
	}))

	generalArea.find(".dropdown").dropdown({
		onChange: function(value, text) {
			selectedGameType = text;
			updateGeneralArea();
		}
	})

	$(".general-card").append(generalArea);
}

var updateGeneralArea = function() {
	// Only call if stats-area exists
	$(".general-card .stats-area").empty()
	$(".general-card .stats-area").append(generateGeneralStats());
}

var generateFiltersArea = function() {
	loadMonths();
	var filtersArea = $(filtersAreaTemplate({
		month: templateMonths,
		start: playerData.RecordStart
	}));

	filtersArea.find("#general-dropdown").dropdown({
		onChange: function(value, text) {
			selectedFilter = text;
			regenerate();
		}
	});

	filtersArea.find("#date-dropdown").dropdown({
		onChange: function(value, text) {
			selectedDate = text;
			selectedCollection = dataMonths[selectedDate];
			regenerate();
		}
	});

	$(".filters-area").append(filtersArea);
}

var indexChampions = function() {
	var championIds = Object.keys(selectedCollection.Champions);

	for (var i = 0; i < championIds.length; i++) {
		var championId = championIds[i];
		championsSearchIndex[championsDatabase[championId].name] = championId;
	}
	return
}

var championSearch = function(query) {
	var results = []; // In champion ID form plz.

	if (query.trim() == "") {
		$("#champion-input").val("");
		// Sort by number of games
		var championIds = Object.keys(selectedCollection.Champions);

		championIds.sort(function(a, b) {
			var championA = selectedCollection.Champions[a];
			var championB = selectedCollection.Champions[b]
			return (championB.Wins + championB.Losses)
				- (championA.Wins + championA.Losses);
		})

		results = championIds;
	} else {
		// First get indexOf == 0
		// Followed by everything else in natrual order.
		var exact = false;
		var startsWith = [];
		var contains = [];
		for (var championName in championsSearchIndex) {
			var lowerChampionName = championName.toLowerCase();
			if (lowerChampionName == query) {
				exact = championsSearchIndex[championName];
			} else if (lowerChampionName.indexOf(query) >= 0) {
				if (lowerChampionName.indexOf(query) == 0) {
					startsWith.push(championsSearchIndex[championName]);
				} else {
					contains.push(championsSearchIndex[championName]);
				}
			}
		}

		results = startsWith.concat(contains);
		if (exact) {
			results.splice(0, 0, exact);
		}
	}
	if (results.length > 10) {
		return results.splice(0, 10);
	}
	return results
}

var generateChampionCards = function() {
	$(".champion-cards").empty();
	$(".more-champions").empty();

	var results = championSearch($("#champion-input").val().toLowerCase())
	var renderInput = [];
	for (var i = 0; i < results.length; i++) {
		var renderArgs = {
			imageName: championsDatabase[results[i]].image,
			displayName: championsDatabase[results[i]].name,
			dateFilter: selectedDate,
			stats: generateChampionStats(results[i]),
		}
		var championCard = championCardTemplate(renderArgs);

		if (window.innerWidth <= 991) {
			$(".champion-cards").append(championCard)
		} else if (i == 0) {
			$(".champion-cards").append(championCard)
		} else if (i % 2 == 0) {
			$(".more-champions").append(championCard)
		} else {
			$(".champion-cards").append(championCard)
		}
	}

	if (results.length <= 0) {
		$(".champion-cards").append('<p class="no-champions">No results</p>');
	}
}

var regenerate = function() {
	generateGeneralArea();
	generateChampionCards();
}

generateFiltersArea();
indexChampions();

$("#champion-input").on("input", function() {
	generateChampionCards();
})

regenerate();

$(".profile").append(profileTemplate({
	username: summonerName,
	region: regionCodes[playerData.Region],
	regionCode: playerData.Region.toUpperCase(),
}));

