
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

// Database management: Performs back-end AJAX queries for data

var gameVersion = 6162
var championsURL = "//ddragon.leagueoflegends.com/cdn/6.20.1/data/en_US/champion.json"
var imageURL = "//ddragon.leagueoflegends.com/cdn/6.20.1/img/champion/"

var championsDatabase = {}
// Database indexed by stringed key, and {name: "Aatrox", image: "Aatrox.png"}

var regionsDatabase = {};
// Array indexed by player names and their regions
var playersDatabase = [];
// Array full of player names, sorted.

var storedVersion = localStorage["version"]
var lastPlayerUpdate = localStorage["playersUpdate"]

if (typeof(onDatabaseLoaded) == "undefined" || !onDatabaseLoaded) {
	onDatabaseLoaded = function() {}
}

var parseWrapper = function(content) {
	try {
		var ret = JSON.parse(content);
		return ret;
	} catch (err) {
		return false;
	}
}

var processDatabase = function(database) {
	for (key in database.data) {
		var champion = database.data[key]
		championsDatabase[champion.key] = {
			name: champion.name,
			image: champion.image.full,
		}
	}
	console.log("Champion database built!")

	localStorage["version"] = gameVersion
	localStorage["champions"] = JSON.stringify(championsDatabase)
	onDatabaseLoaded();
}

var buildChampionsDatabase = function() {
	$.get(championsURL, null, null, "json")
	.done(function(data) {
		try {
			processDatabase(data)
		} catch (err) {
			console.error(err);
			alert("Failed to parse champions database! " +
				"Try refreshing the page, although you may need to report " +
				"this to me (1lanncontact@gmail.com or /u/1lann on reddit). " +
				"Error message: " + err.message)
		}
	})
	.fail(function() {
		alert("Failed to download champions database! Try refreshing the page.")
	})
}

var checkChampionsDatabase = function() {
	if (typeof(storedVersion) == "undefined") {
		console.log("Welcome! Building champion database...");
		buildChampionsDatabase();
	} else if (parseInt(storedVersion) < gameVersion) {
		console.log("Downloading champion database update...");
		buildChampionsDatabase();
	} else {
		var rawDatabase = localStorage["champions"];
		if (typeof(rawDatabase) == "undefined") {
			console.log("Corrupted champion database? Rebuiliding...");
			buildChampionsDatabase();
		} else {
			var output = parseWrapper(rawDatabase)
			if (typeof(output) != "undefined" &&
				output && Object.keys(output).length > 100) {
				console.log("Using champion database from local storage");
				championsDatabase = output;
				onDatabaseLoaded();
			} else {
				console.log("Tampered champion database? Rebuiliding...");
				buildChampionsDatabase();
			}
		}
	}
}

var loadPlayersDatabase = function(database) {
	if (typeof(database) == "undefined" || !database) {
		var rawDatabase = localStorage["players"];
		if (typeof(rawDatabase) == "undefined") {
			localStorage.removeItem("players");
			localStorage.removeItem("playersUpdate");
			alert("Failed to load players database! [1] " +
				"Try refreshing the page.");
			return;
		}
		expandedDatabase = LZString.decompress(rawDatabase)
		readDatabase = parseWrapper(expandedDatabase)
		if (typeof(readDatabase) == "undefined" || !readDatabase) {
			localStorage.removeItem("players");
			localStorage.removeItem("playersUpdate");
			alert("Failed to parse players database! [2] " +
				"Try refreshing the page.");
			return;
		}
		if (typeof(readDatabase.players) == "undefined" ||
				!readDatabase.players ||
				typeof(readDatabase.regions) == "undefined" ||
				!readDatabase.regions) {
			localStorage.removeItem("players");
			localStorage.removeItem("playersUpdate");
			alert("Failed to parse players database! [3] " +
				"Try refreshing the page.");
			return;
		}

		console.log("Loaded players database from local storage")
		playersDatabase = readDatabase.players;
		regionsDatabase = readDatabase.regions;
	} else {
		// The array is already sorted on the server. I think.
		playersDatabase = [];
		regionsDatabase = {};
		var parsedDatabase = parseWrapper(database)
		if (typeof(parsedDatabase) == "undefined" || !parsedDatabase ||
				typeof(parsedDatabase.Players) == "undefined" ||
				!parsedDatabase.Players ||
				typeof(parsedDatabase.Time) == "undefined" ||
				!parsedDatabase.Time) {
			alert("Failed to parse players database! [4] " +
				"Try refreshing the page, although you may need to report " +
				"this to me (1lanncontact@gmail.com or /u/1lann on reddit). ");
			return;
		}

		for (var i = 0; i < parsedDatabase.Players.length; i++) {
			var name = parsedDatabase.Players[i].summonerName
			var region = parsedDatabase.Players[i].region
			if (name in regionsDatabase) {
				regionsDatabase[name].push(region)
			} else {
				regionsDatabase[name] = [region];
				playersDatabase.push(name);
			}
		}

		var storeDatabase = {
			players: playersDatabase,
			regions: regionsDatabase
		}

		var textDatabase = JSON.stringify(storeDatabase);

		console.log("Stored updated players database")
		localStorage["playersUpdate"] = parsedDatabase.Time.toString();
		localStorage["players"] = LZString.compress(textDatabase);
	}
}

var checkPlayersDatabase = function() {
	if (typeof(lastPlayerUpdate) == "undefined" ||
			typeof(localStorage["players"]) == "undefined") {
		lastPlayerUpdate = 0;
	} else {
		// Check for legacy non-compression
		var compressionTest = localStorage["players"];
		if (compressionTest.indexOf('{"players":') == 0) {
			localStorage.removeItem("players");
			lastPlayerUpdate = 0;
		} else {
			// Preload the database from localStorage
			loadPlayersDatabase();
		}
	}
	$.post("/players", {"lastupdate": lastPlayerUpdate}, null, "text")
	.done(function(resp){
		if (resp == "error") {
			alert("Server error while requesting for players database!");
		} else if (resp != "false") {
			loadPlayersDatabase(resp);
		}
	})
}

checkChampionsDatabase();
checkPlayersDatabase();
