document.onreadystatechange = function () {

	function checkConnection() {
		var request = new XMLHttpRequest();
		request.open('GET', 'http://localhost:5000/ping');
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
		request.open('GET', 'http://localhost:5000/filescount');
		request.onreadystatechange = function () {
			if (request.readyState == 4) {
				if (request.status === 200) {
					var data = JSON.parse(request.responseText);

					if (data.message !== 'Not Found') {
						footer.innerHTML = `searching over ${data.count} files`
					}
				} else {
					footer.innerHTML = 'not found'
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
		var request = new XMLHttpRequest();

		request.open('GET', 'http://localhost:5000/search?q=' + query);
		request.onreadystatechange = function () {
			if (request.readyState == 4) {
				if (request.status === 200) {
					var data = JSON.parse(request.responseText);

					if (data.files != null) {
						message.style.display = "none"
						result.style.display = "block"

						content = `<p>${data.files.length} results for '${data.query}'</p>`;
						content += '<table style="border-collapse: collapse;border-spacing: 0;">'
						content += '<thead>'
						content += '<tr>'
						content += '<th>Name</th>'
						content += '<th>Modified</th>'
						content += '<th>Size</th>'
						content += '</tr>'
						content += '</thead>'
						content += '<tbody>'

						for (var i = 0; i < data.files.length; i++) {
							content += '<tr>'
							content += `<td><a href="vip.html?id=${data.files[i].id}">${data.files[i].name}</a></td>`
							content += `<td>${data.files[i].date_modified}</td>`
							content += `<td>${data.files[i].size}</td>`
							content += '</tr>'
						}

						content += '</tbody>'
						content += '</table>'

						result.innerHTML = content;
					} else {
						result.style.display = "none";
						message.style.display = "block";
						message.innerHTML = `no results for '${data.query}'`;
					}
				} else {
					result.style.display = "none";
					message.style.display = "block";
					message.innerHTML = 'there was an error';
				}
			}
		}

		request.setRequestHeader("Access-Control-Allow-Origin", "*")
		request.setRequestHeader("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers")
		request.send();
	}

	if (document.readyState === 'complete') {
		checkConnection();

		var result = document.querySelector('#result');
		var searchForm = document.querySelector('#search-form');

		welcomeMessage()

		searchForm.addEventListener('submit', function (e) {
			e.preventDefault()
			search()
		});
	}
}