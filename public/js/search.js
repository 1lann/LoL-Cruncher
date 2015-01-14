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

serverDatabase = {
	"1lann": ["oce", "na", "kr"],
	"Sloganmaker": ["oce", "na", "euw"],
	"1lamb": ["br"],
}

searchDatabase = ["1lamb", "1lann", "Sloganmaker"];

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
		var searchElement = searchDatabase[start]
		searchElement = searchElement.substring(0, query.length);
		searchElement = searchElement.toLowerCase()
		if (start < searchDatabase.length && (searchElement == query)) {
			results.push(searchDatabase[start]);
			start++;
			if (start >= searchDatabase.length) {
				return results;
			}
		} else {
			return results;
		}
	}
}

var reverseAndFlood = function(start, query) {
	for (var reverseStart = start; reverseStart >= 0; reverseStart--) {
		var searchElement = searchDatabase[reverseStart];
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
	var test = searchDatabase[center];
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
	return binarySearch(query.toLowerCase(), 0, searchDatabase.length - 1)
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
	for (var i = 0; i < results.length; i++) {
		var serverResult = serverDatabase[results[i]];

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
	displayResults(realSearch.val())
})

$("#fake-search").on("input", function() {
	$(".landing").hide();
	$(".search").show();
	realSearch.focus();
	realSearch.val($(this).val());
	realSearch.triggerHandler("input");
})
