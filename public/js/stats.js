
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

var startDate = "since " + playerData.RecordStart;
var selectedDate = startDate;
var selectedFilter = "All";
var selectedQueue = "all queues";
var selectedDisplay = "cards";

var templateMonths = [];
var templateQueues = [];
var monthResolver = {};

var championsSearchIndex = {};
var championFilter = ""

var selectedCollection;

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
	var months = Object.keys(playerData.AllQueues.MonthlyStats)
	months.sort();
	for (var i = 0; i < months.length; i++) {
		var monthKey = "for " + getStringedDate(months[i]);
		monthResolver[monthKey] = months[i];
		templateMonths.push({text: monthKey});
	}
}

var loadQueues = function() {
	var queues = Object.keys(playerData.QueueStats)
	for (var i = 0; i < queues.length; i++) {
		templateQueues.push({name: queues[i]})
	}
}

var selectCollection = function() {
	var queueSelection;

	if (selectedQueue == "all queues") {
		queueSelection = playerData.AllQueues;
	} else {
		queueSelection = playerData.QueueStats[selectedQueue];
	}

	if (selectedDate == startDate) {
		selectedCollection = queueSelection.AllMonths
	} else {
		selectedCollection =
			queueSelection.MonthlyStats[monthResolver[selectedDate]]
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
		if (numGames <= 0) {
			outputStats.push({
				label: "Winrate",
				data: "0%",
			});
		} else {
			outputStats.push({
				label: "Winrate",
				data: Math.round((statsSource.Wins/numGames) * 100) + "%",
			});
		}

		var redGames = statsSource.Red.Wins + statsSource.Red.Losses
		if (redGames <= 0) {
			outputStats.push({
				label: "Red team winrate",
				data: "0%",
			});
		} else {
			outputStats.push({
				label: "Red team winrate",
				data: Math.round((statsSource.Red.Wins/redGames) * 100) + "%",
			});
		}

		var blueGames = statsSource.Blue.Wins + statsSource.Blue.Losses
		if (blueGames <= 0) {
			outputStats.push({
				label: "Blue team winrate",
				data: "0%",
			});
		} else {
			outputStats.push({
				label: "Blue team winrate",
				data: Math.round((statsSource.Blue.Wins/blueGames) * 100) + "%",
			});
		}

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
	// TODO: Add gold stats

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
			label: "Gold earned",
			data: getGoldAmount(statsSource.GoldEarned),
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
		if (numGames <= 0) {
			outputStats.push({
				label: "Winrate",
				data: "0%",
			});
		} else {
			outputStats.push({
				label: "Winrate",
				data: Math.round((statsSource.Wins/numGames) * 100) + "%",
			});
		}
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

var generate

var generateGeneralArea = function() {
	$(".general-card").empty();

	var statsArea = generateGeneralStats();

	var generalArea = $(generalCardTemplate({
		dateFilter: selectedDate + " for " + selectedQueue,
		stats: statsArea,
	}))

	$(".general-card").append(generalArea);
}

var updateGeneralArea = function() {
	// Only call if stats-area exists
	$(".general-card .stats-area").empty()
	$(".general-card .stats-area").append(generateGeneralStats());
}

var generateFiltersArea = function() {
	loadQueues();
	loadMonths();
	var filtersArea = $(filtersAreaTemplate({
		month: templateMonths,
		start: playerData.RecordStart,
		queueTypes: templateQueues,
	}));

	filtersArea.find("#general-dropdown").dropdown({
		onChange: function(value, text) {
			selectedFilter = text;
			regenerate();
		},
		on: "hover"
	});

	filtersArea.find("#date-dropdown").dropdown({
		onChange: function(value, text) {
			selectedDate = text;
			regenerate();
		},
		on: "hover"
	});

	filtersArea.find("#queue-dropdown").dropdown({
		onChange: function(value, text) {
			selectedQueue = text;
			regenerate();
		},
		on: "hover"
	});

	filtersArea.find("#display-dropdown").dropdown({
		onChange: function(value, text) {
			selectedDisplay = text;
			regenerate();
		},
		on: "hover"
	});

	filtersArea.find(".glyphicon.glyphicon-info-sign").popover();

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
			dateFilter: selectedDate + " for " + selectedQueue,
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

var generateTable = function() {
	$(".stats-table-container").empty();

	var stats = [];
	var tableHeader = "";

	if (selectedFilter == "All") {
		tableHeader = tableAllHeader;

		for (var championId in selectedCollection.Champions) {
			var statsSource = selectedCollection.Champions[championId];

			var renderArgs = {
				imageName: championsDatabase[championId].image,
				championName: championsDatabase[championId].name,
				games: (statsSource.Wins + statsSource.Losses).toString(),
				wins: statsSource.Wins.toString(),
				losses: statsSource.Losses.toString(),
				minions: statsSource.MinionsKilled.toString(),
				jungle: statsSource.MonstersKilled.toString(),
				gold: getGoldAmount(statsSource.GoldEarned),
				wards: statsSource.WardsPlaced.toString(),
				kills: statsSource.Kills.toString(),
				deaths: statsSource.Deaths.toString(),
				assists: statsSource.Assists.toString()
			}

			stats.push({stats: tableRowTemplate(renderArgs)});
		}
	} else if (selectedFilter == "Rates/average") {
		tableHeader = tableRatesHeader;

		for (var championId in selectedCollection.Champions) {
			var statsSource = selectedCollection.Champions[championId];

			var numGames = statsSource.Wins + statsSource.Losses;
			var timePlayed = statsSource.TimePlayed;

			var winrate = "0%";
			if (numGames > 0) {
				winrate = Math.round((statsSource.Wins/numGames) * 100) + "%";
			}

			var renderArgs = {
				imageName: championsDatabase[championId].image,
				championName: championsDatabase[championId].name,
				games: numGames.toString(),
				wins: winrate,
				minions: oneDecRound(statsSource.MinionsKilled/(timePlayed/600)),
				jungle: oneDecRound(statsSource.MonstersKilled/(timePlayed/600)),
				gold: getGoldAmount(statsSource.GoldEarned/(timePlayed/600)),
				wards: oneDecRound(statsSource.WardsPlaced/numGames),
				kills: oneDecRound(statsSource.Kills/numGames),
				deaths: oneDecRound(statsSource.Deaths/numGames),
				assists: oneDecRound(statsSource.Assists/numGames)
			}

			stats.push({stats: tableRowTemplate(renderArgs)});
		}
	}

	var renderArgs = {
		tableHeader: tableHeader,
		championStats: stats
	}

	var tableElement = $(tableStatsTemplate(renderArgs));

	$(".stats-table-container").append(tableElement);
	$(".stats-table-container").scroll(function() {
		$(window).trigger("resize.stickyTableHeaders");
	})

	tableElement.stickyTableHeaders();

	if (selectedFilter == "All") {
		tableElement.tablesorter({
			sortList: [[2, 1]],
			headers: {
				0: {
					sorter: false
				},
				7: {
					sorter: "gold"
				}
			}
		});
	} else if (selectedFilter == "Rates/average") {
		tableElement.tablesorter({
			sortList: [[2, 1]],
			headers: {
				0: {
					sorter: false
				},
				6: {
					sorter: "gold"
				}
			}
		});
	}
}

var regenerate = function() {
	selectCollection();

	if (selectedDisplay == "cards") {
		$("#table-area").hide();
		$("#cards-area").show();
		indexChampions();
		generateGeneralArea();
		generateChampionCards();
	} else {
		$("#cards-area").hide();
		$("#table-area").show();

		generateTable();
	}
}

var onDatabaseLoaded = function() {
	console.log("Database loaded, loading player data...")

	generateFiltersArea();
	selectCollection();

	regenerate();

	$("#champion-input").on("input", function() {
		generateChampionCards();
	})

	$(".profile").append(profileTemplate({
		username: summonerName,
		imageName: summonerName.replace(" ", ""),
		region: regionCodes[playerData.Region],
		regionCode: playerData.Region.toUpperCase(),
	}));
}

$.tablesorter.addParser({
	// set a unique id
	id: "gold",
	is: function(s) {
		// return false so this parser is not auto detected
		return false;
	},
	format: function(s) {
		// format your data for normalization
		if (s.indexOf("k") > 0) {
			return Math.round(parseFloat(s) * 1000);
		} else {
			return Math.round(parseFloat(s) * 1000000)
		}
	},
	// set type, either numeric or text
	type: "numeric"
});

$.tablesorter.defaults.sortInitialOrder = "desc";

