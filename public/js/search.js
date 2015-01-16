
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


var realSearch = $("#real-search")

regionCodes = {
	"na": "North America",
	"euw": "Europe West",
	"eune": "Europe North East",
	"lan": "Latin America North",
	"las": "Latin America South",
	"oce": "Oceania",
	"br": "Brazil",
	"ru": "Russia",
	"kr": "Korea",
	"tr": "Turkey",
}

regionOrder = {
	"na": 1,
	"euw": 2,
	"eune": 3,
	"lan": 4,
	"las": 5,
	"oce": 6,
	"br": 7,
	"ru": 8,
	"kr": 9,
	"tr": 10,
}

var noResults = '\
<div class="no-results">\
	<div class="container">\
		<p class="exclamation">Select a Region</p>\
		<div class="button-container col-md-3 col-md-offset-3 col-xs-6">\
			<a id="na">North America</a>\
		</div>\
		<div class="button-container col-md-3 col-xs-6">\
			<a id="euw">Europe West</a>\
		</div>\
		<div class="button-container col-md-3 col-md-offset-3 col-xs-6">\
			<a id="eune">Europe North East</a>\
		</div>\
		<div class="button-container col-md-3 col-xs-6">\
			<a id="lan">Latin America North</a>\
		</div>\
		<div class="button-container col-md-3 col-md-offset-3 col-xs-6">\
			<a id="las">Latin America South</a>\
		</div>\
		<div class="button-container col-md-3 col-xs-6">\
			<a id="oce">Oceania</a>\
		</div>\
		<div class="button-container col-md-3 col-md-offset-3 col-xs-6">\
			<a id="br">Brazil</a>\
		</div>\
		<div class="button-container col-md-3 col-xs-6">\
			<a id="ru">Russia</a>\
		</div>\
		<div class="button-container col-md-3 col-md-offset-3 col-xs-6">\
			<a id="kr">Korea</a>\
		</div>\
		<div class="button-container col-md-3 col-xs-6">\
			<a id="tr">Turkey</a>\
		</div>\
	</div>\
</div>'

var resultSource = '\
<div class="result">\
	<div class="container">\
		<img src="http://avatar.leagueoflegends.com/{{regionCode}}/{{name}}.png">\
		<p class="username">{{name}}</p>\
		<div class="region-wrapper">\
			<p class="region">{{region}}</p>\
		</div>\
	</div>\
</div>'

var notHere = '\
<div class="result">\
	<div class="container">\
		<p class="not-here">Not here?</p>\
	</div>\
</div>'

var enterText = '\
<div class="enter-text">\
	<div class="container">\
		<p class="exclamation">Enter a Summoner Name</p>\
	</div>\
</div>'


var resultTemplate = Handlebars.compile(resultSource)
var displayedResults = [];

var floodResults = function(query, start) {
	var results = [];

	while (true) {
		var searchElement = playersDatabase[start]
		searchElement = searchElement.substring(0, query.length);
		searchElement = searchElement.toLowerCase()
		if (start < playersDatabase.length && (searchElement == query)) {
			results.push(playersDatabase[start]);
			start++;
			if (start >= playersDatabase.length) {
				return results;
			}
		} else {
			return results;
		}
	}
}

var reverseAndFlood = function(start, query) {
	for (var reverseStart = start; reverseStart >= 0; reverseStart--) {
		var searchElement = playersDatabase[reverseStart];
		searchElement = searchElement.substring(0, query.length);
		searchElement = searchElement.toLowerCase();

		if (searchElement != query) {
			reverseStart++;
			break;
		}
	}

	if (reverseStart < 0) {
		reverseStart = 0;
	}

	return floodResults(query, reverseStart);
}

var binarySearch = function(query, start, end) {
	if (start > end) {
		return [];
	}

	var center = Math.ceil((start + end) / 2);
	var test = playersDatabase[center];
	if (test == center) {
		return floodResults(query, center)
	} else if (test.substring(0, query.length).toLowerCase() == query) {
		return reverseAndFlood(center, query)
	} else if (test.substring(0, query.length).toLowerCase() > query) {
		return binarySearch(query, start, center - 1);
	} else {
		return binarySearch(query, center + 1, end);
	}
}

var search = function(query) {
	var lowerQuery = query.toLowerCase();
	var results = binarySearch(lowerQuery, 0, playersDatabase.length - 1);
	return results;
}

var sortRegions = function(regions) {
	return regions.sort(function(a, b) {
		return regionOrder[a] - regionOrder[b]
	})
}

var displayResults = function(query) {
	$(".search-results").empty();
	displayedResults = [];
	if (query == "") {
		$(".search-results").append(enterText)
		return;
	}
	var results = search(query);
	if (results.length <= 0) {
		$(".search-results").append(noResults)
		return;
	}
	var resultNumber = 0;
	for (var i = 0; i < Math.min(results.length, 5); i++) {
		var serverResult = sortRegions(regionsDatabase[results[i]]);

		for (var r = 0; r < serverResult.length; r++) {
			var arguments = {
				"name": results[i],
				"region": regionCodes[serverResult[r]],
				"regionCode": serverResult[r],
			}

			var resultElement = resultTemplate(arguments);
			$(".search-results").append(resultElement);

			displayedResults.push(resultElement);
			resultNumber++;
		}
	}
	var notHereElement = $(notHere);
	notHereElement.click(function() {
		$(".search-results").empty();
		$(".search-results").append(noResults);
	})
	$(".search-results").append(notHereElement);

}

realSearch.on("input", function() {
	if (realSearch.val().trim() == "") {
		realSearch.val("")
	}
	$("#main").hide();
	displayResults(realSearch.val())
})

$("#fake-search").on("input", function() {
	$(".landing").hide();
	$(".search").show();
	realSearch.focus();
	realSearch.val($(this).val());
	realSearch.triggerHandler("input");
})
