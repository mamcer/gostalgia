document.onreadystatechange = function () {

	var showAdvancedSearch = 0;

	function checkConnection() {
		var request = new XMLHttpRequest();
		request.open('GET', config.api+'/ping');
		request.onreadystatechange = function () {
			if (request.readyState == 4) {
				if (request.status != 200) {
					alert('There is no connection with the back end api');
				}
			}
		}

		request.setRequestHeader("Access-Control-Allow-Origin", "*")
		request.setRequestHeader("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
		request.send();
	}

	function welcomeMessage() {
		var request = new XMLHttpRequest();
		request.open('GET', config.api+'/filescount');
		request.onreadystatechange = function () {
			if (request.readyState == 4) {
				if (request.status === 200) {
					var data = JSON.parse(request.responseText);

					if (data.message !== 'Not Found') {
						message.innerHTML = `searching over ${data.count} files`
					}
				} else {
					message.innerHTML = 'not found'
				}
			}
		}

		request.setRequestHeader("Access-Control-Allow-Origin", "*")
		request.setRequestHeader("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
		request.send();
	}

	function search() {
		result.innerHTML = ''
		var query = document.querySelector('#search').value;
		location = 'search.html?q=' + query
	}

	function advancedSearch() {
		showAdvancedSearch = !showAdvancedSearch;
		if (showAdvancedSearch == 0) {
			advancedSeachPanel.style.display = "none";
		} else {
			advancedSeachPanel.style.display = "inline";
		}
			
	}

	if (document.readyState === 'complete') {
		checkConnection();

		var message = document.querySelector('#message');
		var result = document.querySelector('#result');
		var searchForm = document.querySelector('#search-form');
		var advancedSeachPanel = document.querySelector('#advanced-search-panel');
		document.getElementById("advanced-search").addEventListener("click", advancedSearch, false);

		welcomeMessage()

		searchForm.addEventListener('submit', function (e) {
			e.preventDefault()
			search()
		});
	}
}