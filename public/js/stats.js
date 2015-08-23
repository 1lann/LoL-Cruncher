
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

var startDate = "since " + playerData.rs
var selectedDate = startDate
var selectedFilter = "All"
var selectedQueue = "all queues"
var selectedDisplay = "cards"

var templateMonths = []
var templateQueues = []
var monthResolver = {}

var championsSearchIndex = {}
var championFilter = ""

var queueSelection = "all"
var dateSelection = "all"

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

var dateRegex = /(\d+)\s(\d+)/

var getStringedDate = function(date) {
	var resp = date.match(dateRegex)
	return monthsKeys[resp[2]] + " " + resp[1]
}

var loadMonths = function() {
	var monthMap = {}
	for (var i = 0; i < playerData.detailed.length; i++) {
		var period = playerData.detailed[i].p
		if (!monthMap[period] && period != "all") {
			monthMap[period] = true
		}
	}

	var months = Object.keys(monthMap)
	months.sort()

	for (var i = 0; i < months.length; i++) {
		var monthKey = "for " + getStringedDate(months[i])
		monthResolver[monthKey] = months[i]
		templateMonths.push({text: monthKey})
	}
}

var loadQueues = function() {
	var queueMap = {}
	for (var i = 0; i < playerData.detailed.length; i++) {
		var queue = playerData.detailed[i].q
		if (!queueMap[queue] && queue != "all") {
			queueMap[queue] = true
			templateQueues.push({name: queue})
		}
	}
}

var selectCollection = function() {
	if (selectedQueue == "all queues") {
		queueSelection = "all"
	} else {
		queueSelection = selectedQueue
	}

	if (selectedDate == startDate) {
		dateSelection = "all"
	} else {
		dateSelection = monthResolver[selectedDate]
	}
}

var getHumanTime = function(seconds) {
	var numdays = Math.floor(seconds / 86400)
	var numhours = Math.floor((seconds % 86400) / 3600)
	var numminutes = Math.floor(((seconds % 86400) % 3600) / 60)
	var construct = ""

	if (numdays == 1) {
		construct = "1 day"
	} else if (numdays > 1) {
		construct = numdays + " days"
	}

	if (numhours == 1) {
		if (construct != "") {
			construct = construct + ", "
		}
		construct = construct + "1 hour"
	} else if (numhours > 1) {
		if (construct != "") {
			construct = construct + ", "
		}
		construct = construct + numhours + " hours"
	}

	if (construct != "") {
		construct = construct + ", and "
	}
	if (numminutes == 1) {
		construct = construct + "1 minute"
	} else {
		construct = construct + numminutes + " minutes"
	}
	return construct
}

var getGoldAmount = function(gold) {
	if (gold > 999999) {
		var prefix = Math.round(gold/10000) / 100
		return prefix + "M"
	} else {
		var prefix = Math.round(gold/100) / 10
		return prefix + "k"
	}
}

var oneDecRound = function(num) {
	return (Math.round(num * 10) / 10).toString()
}

var getDetailedCollection = function(date, queue) {
	for (var i = 0; i < playerData.detailed.length; i++) {
		var currentDetailed = playerData.detailed[i]
		if (currentDetailed.p == date && currentDetailed.q == queue) {
			return currentDetailed
		}
	}
}

var getBasicCollection = function(date, queue, champion) {
	for (var i = 0; i < playerData.basic.length; i++) {
		var currentChampion = playerData.basic[i]
		if (currentChampion.p == date && currentChampion.q == queue &&
			currentChampion.c == champion) {
			return currentChampion
		}
	}
}

var getListChampions = function(date, queue) {
	var championMap = {}
	for (var i = 0; i < playerData.basic.length; i++) {
		var basicData = playerData.basic[i]
		if (basicData.q == queue && basicData.p == date) {
			championMap[basicData.c] = true
		}
	}

	return Object.keys(championMap)
}

// You may want to collapse these functions in your IDE/Text Editor
var generateGeneralStats = function() {
	var outputStats = []
	var statsSource = getDetailedCollection(dateSelection, queueSelection)

	var spentPlaying = timePlayingTemplate({
		time: getHumanTime(statsSource.t)
	})

	if (selectedFilter == "All") {
		outputStats.push({
			label: "Games played",
			data: (statsSource.w + statsSource.l).toString(),
		})
		outputStats.push({
			label: "Games won",
			data: statsSource.w.toString(),
		})
		outputStats.push({
			label: "Games lost",
			data: statsSource.l.toString(),
		})
		outputStats.push({
			label: "Games played on red",
			data: (statsSource.r.w + statsSource.r.l).toString(),
		})
		outputStats.push({
			label: "Games played on blue",
			data: (statsSource.b.w + statsSource.b.l).toString(),
		})
		outputStats.push({
			label: "Minions killed",
			data: statsSource.m.toString(),
		})
		outputStats.push({
			label: "Jungle monsters killed",
			data: statsSource.n.toString(),
		})
		outputStats.push({
			label: "Gold earned",
			data: getGoldAmount(statsSource.g),
		})
		outputStats.push({
			label: "Wards placed",
			data: statsSource.wp.toString(),
		})
		outputStats.push({
			label: "Wards killed",
			data: statsSource.wk.toString(),
		})
		outputStats.push({
			label: "Kills",
			data: statsSource.k.toString(),
		})
		outputStats.push({
			label: "Deaths",
			data: statsSource.d.toString(),
		})
		outputStats.push({
			label: "Assists",
			data: statsSource.a.toString(),
		})
		outputStats.push({
			label: "Double kills",
			data: statsSource.dk.toString(),
		})
		outputStats.push({
			label: "Triple kills",
			data: statsSource.tk.toString(),
		})
		outputStats.push({
			label: "Quadra kills",
			data: statsSource.qk.toString(),
		})
		outputStats.push({
			label: "Pentakills",
			data: statsSource.pk.toString(),
		})
	} else if (selectedFilter == "Rates/average") {
		var numGames = statsSource.w + statsSource.l
		var timePlayed = statsSource.t
		outputStats.push({
			label: "Games played",
			data: numGames.toString(),
		})
		outputStats.push({
			label: "Average game time in minutes",
			data: Math.round(timePlayed/numGames/60),
		})
		if (numGames <= 0) {
			outputStats.push({
				label: "Winrate",
				data: "0%",
			})
		} else {
			outputStats.push({
				label: "Winrate",
				data: Math.round((statsSource.w/numGames) * 100) + "%",
			})
		}

		var redGames = statsSource.r.w + statsSource.r.l
		if (redGames <= 0) {
			outputStats.push({
				label: "Red team winrate",
				data: "0%",
			})
		} else {
			outputStats.push({
				label: "Red team winrate",
				data: Math.round((statsSource.r.w/redGames) * 100) + "%",
			})
		}

		var blueGames = statsSource.b.w + statsSource.b.l
		if (blueGames <= 0) {
			outputStats.push({
				label: "Blue team winrate",
				data: "0%",
			})
		} else {
			outputStats.push({
				label: "Blue team winrate",
				data: Math.round((statsSource.b.w/blueGames) * 100) + "%",
			})
		}

		outputStats.push({
			label: "Minions killed per 10 minutes",
			data: oneDecRound(statsSource.m/(timePlayed/600)),
		})
		outputStats.push({
			label: "Jungle monsters killed per 10 minutes",
			data: oneDecRound(statsSource.n/(timePlayed/600)),
		})
		outputStats.push({
			label: "Gold earned per 10 minutes",
			data: getGoldAmount(statsSource.g/(timePlayed/600)),
		})
		outputStats.push({
			label: "Wards placed per game",
			data: oneDecRound(statsSource.wp/numGames),
		})
		outputStats.push({
			label: "Wards killed per game",
			data: oneDecRound(statsSource.wk/numGames),
		})
		outputStats.push({
			label: "Kills per game",
			data: oneDecRound(statsSource.k/numGames),
		})
		outputStats.push({
			label: "Deaths per game",
			data: oneDecRound(statsSource.d/numGames),
		})
		outputStats.push({
			label: "Assists per game",
			data: oneDecRound(statsSource.a/numGames),
		})
	}

	return spentPlaying + statsTemplate({statsRow: outputStats})
}

var generateChampionStats = function(championId) {
	// TODO: Add gold stats

	var outputStats = []
	var statsSource = getBasicCollection(dateSelection, queueSelection,
		championId)

	var spentPlaying = timePlayingTemplate({
		time: getHumanTime(statsSource.t)
	})

	if (selectedFilter == "All") {
		outputStats.push({
			label: "Games played",
			data: (statsSource.w + statsSource.l).toString(),
		})
		outputStats.push({
			label: "Games won",
			data: statsSource.w.toString(),
		})
		outputStats.push({
			label: "Games lost",
			data: statsSource.l.toString(),
		})
		outputStats.push({
			label: "Minions killed",
			data: statsSource.m.toString(),
		})
		outputStats.push({
			label: "Jungle monsters killed",
			data: statsSource.n.toString(),
		})
		outputStats.push({
			label: "Gold earned",
			data: getGoldAmount(statsSource.g),
		})
		outputStats.push({
			label: "Wards placed",
			data: statsSource.wp.toString(),
		})
		outputStats.push({
			label: "Kills",
			data: statsSource.k.toString(),
		})
		outputStats.push({
			label: "Deaths",
			data: statsSource.d.toString(),
		})
		outputStats.push({
			label: "Assists",
			data: statsSource.a.toString(),
		})
	} else if (selectedFilter == "Rates/average") {
		var numGames = statsSource.w + statsSource.l
		var timePlayed = statsSource.t
		outputStats.push({
			label: "Games played",
			data: numGames.toString(),
		})
		outputStats.push({
			label: "Average game time in minutes",
			data: Math.round(timePlayed/numGames/60),
		})
		if (numGames <= 0) {
			outputStats.push({
				label: "Winrate",
				data: "0%",
			})
		} else {
			outputStats.push({
				label: "Winrate",
				data: Math.round((statsSource.w/numGames) * 100) + "%",
			})
		}
		outputStats.push({
			label: "Minions killed per 10 minutes",
			data: oneDecRound(statsSource.m/(timePlayed/600)),
		})
		outputStats.push({
			label: "Jungle monsters killed per 10 minutes",
			data: oneDecRound(statsSource.n/(timePlayed/600)),
		})
		outputStats.push({
			label: "Gold earned per 10 minutes",
			data: getGoldAmount(statsSource.g/(timePlayed/600)),
		})
		outputStats.push({
			label: "Wards placed per game",
			data: oneDecRound(statsSource.wp/numGames),
		})
		outputStats.push({
			label: "Kills per game",
			data: oneDecRound(statsSource.k/numGames),
		})
		outputStats.push({
			label: "Deaths per game",
			data: oneDecRound(statsSource.d/numGames),
		})
		outputStats.push({
			label: "Assists per game",
			data: oneDecRound(statsSource.a/numGames),
		})
	}

	return spentPlaying + statsTemplate({statsRow: outputStats})
}

var generate

var generateGeneralArea = function() {
	$(".general-card").empty()

	var statsArea = generateGeneralStats()

	var generalArea = $(generalCardTemplate({
		dateFilter: selectedDate + " for " + selectedQueue,
		stats: statsArea,
	}))

	$(".general-card").append(generalArea)
}

var updateGeneralArea = function() {
	// Only call if stats-area exists
	$(".general-card .stats-area").empty()
	$(".general-card .stats-area").append(generateGeneralStats())
}

var generateFiltersArea = function() {
	loadQueues()
	loadMonths()
	var filtersArea = $(filtersAreaTemplate({
		month: templateMonths,
		start: playerData.rs,
		queueTypes: templateQueues,
	}))

	filtersArea.find("#general-dropdown").dropdown({
		onChange: function(value, text) {
			selectedFilter = text
			regenerate()
		},
		on: "hover"
	})

	filtersArea.find("#date-dropdown").dropdown({
		onChange: function(value, text) {
			selectedDate = text
			regenerate()
		},
		on: "hover"
	})

	filtersArea.find("#queue-dropdown").dropdown({
		onChange: function(value, text) {
			selectedQueue = text
			regenerate()
		},
		on: "hover"
	})

	filtersArea.find("#display-dropdown").dropdown({
		onChange: function(value, text) {
			selectedDisplay = text
			regenerate()
		},
		on: "hover"
	})

	filtersArea.find(".glyphicon.glyphicon-info-sign").popover()

	$(".filters-area").append(filtersArea)
}

var indexChampions = function() {
	var championIds = getListChampions(dateSelection, queueSelection)

	for (var i = 0; i < championIds.length; i++) {
		var championId = championIds[i]
		championsSearchIndex[championsDatabase[championId].name] = championId
	}
	return
}

var championSearch = function(query) {
	var results = [] // In champion ID form plz.

	if (query.trim() == "") {
		$("#champion-input").val("")
		// Sort by number of games
		championIds = getListChampions(dateSelection, queueSelection)

		championIds.sort(function(a, b) {
			var championA = getBasicCollection(dateSelection,
			queueSelection, a)
			var championB = getBasicCollection(dateSelection,
			queueSelection, b)
			return (championB.w + championB.l)
				- (championA.w + championA.l)
		})

		results = championIds
	} else {
		// First get indexOf == 0
		// Followed by everything else in natrual order.
		var exact = false
		var startsWith = []
		var contains = []
		for (var championName in championsSearchIndex) {
			var lowerChampionName = championName.toLowerCase()
			if (lowerChampionName == query) {
				exact = championsSearchIndex[championName]
			} else if (lowerChampionName.indexOf(query) >= 0) {
				if (lowerChampionName.indexOf(query) == 0) {
					startsWith.push(championsSearchIndex[championName])
				} else {
					contains.push(championsSearchIndex[championName])
				}
			}
		}

		results = startsWith.concat(contains)
		if (exact) {
			results.splice(0, 0, exact)
		}
	}
	if (results.length > 10) {
		return results.splice(0, 10)
	}
	return results
}

var generateChampionCards = function() {
	$(".champion-cards").empty()
	$(".more-champions").empty()

	var results = championSearch($("#champion-input").val().toLowerCase())
	var renderInput = []
	for (var i = 0; i < results.length; i++) {
		var renderArgs = {
			imageName: championsDatabase[results[i]].image,
			displayName: championsDatabase[results[i]].name,
			dateFilter: selectedDate + " for " + selectedQueue,
			stats: generateChampionStats(results[i]),
		}
		var championCard = championCardTemplate(renderArgs)

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
		$(".champion-cards").append('<p class="no-champions">No results</p>')
	}
}

var generateTable = function() {
	$(".stats-table-container").empty()

	var stats = []
	var tableHeader = ""
	var tableFooter = ""

	if (selectedFilter == "All") {
		tableHeader = tableAllHeader

		var championsList = getListChampions(dateSelection, queueSelection)

		for (var i = 0; i < championsList.length; i++) {
			var championId = championsList[i];

			var statsSource = getBasicCollection(dateSelection, queueSelection,
				championId)

			var renderArgs = {
				imageName: championsDatabase[championId].image,
				championName: championsDatabase[championId].name,
				games: (statsSource.w + statsSource.l).toString(),
				wins: statsSource.w.toString(),
				losses: statsSource.l.toString(),
				minions: statsSource.m.toString(),
				jungle: statsSource.n.toString(),
				gold: getGoldAmount(statsSource.g),
				wards: statsSource.wp.toString(),
				kills: statsSource.k.toString(),
				deaths: statsSource.d.toString(),
				assists: statsSource.a.toString()
			}

			stats.push({stats: tableRowTemplate(renderArgs)})
		}

		var footerStats = getDetailedCollection(dateSelection, queueSelection)

		var footerArgs = {
			games: (footerStats.w + footerStats.l).toString(),
			wins: footerStats.w.toString(),
			losses: footerStats.l.toString(),
			minions: footerStats.m.toString(),
			jungle: footerStats.n.toString(),
			gold: getGoldAmount(footerStats.g),
			wards: footerStats.wp.toString(),
			kills: footerStats.k.toString(),
			deaths: footerStats.d.toString(),
			assists: footerStats.a.toString()
		}

		tableFooter = tableFooterTemplate(footerArgs)
	} else if (selectedFilter == "Rates/average") {
		tableHeader = tableRatesHeader
		var championsList = getListChampions(dateSelection, queueSelection)

		for (var i = 0; i < championsList.length; i++) {
			var championId = championsList[i];

			var statsSource = getBasicCollection(dateSelection, queueSelection,
				championId)

			var numGames = statsSource.w + statsSource.l
			var timePlayed = statsSource.t

			var winrate = "0%"
			if (numGames > 0) {
				winrate = Math.round((statsSource.w/numGames) * 100) + "%"
			}

			var renderArgs = {
				imageName: championsDatabase[championId].image,
				championName: championsDatabase[championId].name,
				games: numGames.toString(),
				wins: winrate,
				minions: oneDecRound(statsSource.m/(timePlayed/600)),
				jungle: oneDecRound(statsSource.n/(timePlayed/600)),
				gold: getGoldAmount(statsSource.g/(timePlayed/600)),
				wards: oneDecRound(statsSource.wp/numGames),
				kills: oneDecRound(statsSource.k/numGames),
				deaths: oneDecRound(statsSource.d/numGames),
				assists: oneDecRound(statsSource.a/numGames)
			}

			stats.push({stats: tableRowTemplate(renderArgs)})
		}

		var footerStats = getDetailedCollection(dateSelection, queueSelection)

		var timePlayed = footerStats.t
		var numGames = (footerStats.w + footerStats.l)

		var winrate = "0%"
		if (numGames > 0) {
			winrate = Math.round((footerStats.w/numGames) * 100) + "%"
		}

		var footerArgs = {
			games: numGames.toString(),
			wins: winrate,
			minions: oneDecRound(footerStats.m/(timePlayed/600)),
			jungle: oneDecRound(footerStats.n/(timePlayed/600)),
			gold: getGoldAmount(footerStats.g/(timePlayed/600)),
			wards: oneDecRound(footerStats.wp/numGames),
			kills: oneDecRound(footerStats.k/numGames),
			deaths: oneDecRound(footerStats.d/numGames),
			assists: oneDecRound(footerStats.a/numGames)
		}

		tableFooter = tableFooterTemplate(footerArgs)
	}

	var renderArgs = {
		tableHeader: tableHeader,
		championStats: stats,
		tableFooter: tableFooter
	}

	var tableElement = $(tableStatsTemplate(renderArgs))

	$(".stats-table-container").append(tableElement)
	$(".stats-table-container").scroll(function() {
		$(window).trigger("resize.stickyTableHeaders")
	})

	tableElement.stickyTableHeaders()

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
		})
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
		})
	}
}

var regenerate = function() {
	selectCollection()

	var selection = getDetailedCollection(dateSelection, queueSelection)
	if (!selection) {
		$("#cards-area").hide()
		$("#table-area").hide()
		$("#warning-area").show()
		return
	} else {
		$("#warning-area").hide()
	}

	if (selectedDisplay == "cards") {
		$("#table-area").hide()
		$("#cards-area").show()
		indexChampions()
		generateGeneralArea()
		generateChampionCards()
	} else {
		$("#cards-area").hide()
		$("#table-area").show()

		generateTable()
	}
}

var onDatabaseLoaded = function() {
	console.log("Database loaded, loading player data...")

	generateFiltersArea()
	selectCollection()

	regenerate()

	$("#champion-input").on("input", function() {
		generateChampionCards()
	})

	$(".profile").append(profileTemplate({
		username: summonerName,
		imageName: summonerName.replace(" ", ""),
		region: regionCodes[playerData.r],
		regionCode: playerData.r.toUpperCase(),
	}))
}

$.tablesorter.addParser({
	// set a unique id
	id: "gold",
	is: function(s) {
		// return false so this parser is not auto detected
		return false
	},
	format: function(s) {
		// format your data for normalization
		if (s.indexOf("k") > 0) {
			return Math.round(parseFloat(s) * 1000)
		} else {
			return Math.round(parseFloat(s) * 1000000)
		}
	},
	// set type, either numeric or text
	type: "numeric"
})

$.tablesorter.defaults.sortInitialOrder = "desc"

