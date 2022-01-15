document.onreadystatechange = function () {

	var location

	function getParameterByName(name, url) {
		if (!url) url = window.location.href;
		name = name.replace(/[\[\]]/g, '\\$&');
		var regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)'),
			results = regex.exec(url);
		if (!results) return null;
		if (!results[2]) return '';
		return decodeURIComponent(results[2].replace(/\+/g, ' '));
	}

	function open_item() {
		window.open(location, '_blank');
	}

	function showFile() {
		var id = getParameterByName('id');

		var request = new XMLHttpRequest();
		request.open('GET', 'http://localhost:5000/files/' + id);
		request.onreadystatechange = function () {
			if (request.readyState == 4) {
				if (request.status === 200) {
					var data = JSON.parse(request.responseText);

					name.innerHTML = `${data.name}`;
					path.innerHTML = `${data.path}`;
					date_modified.innerHTML = `${data.date_modified}`;
					size.innerHTML = `${data.size}`;
					ndirectory_id.innerHTML = `${data.ndirectory_id}`;
					nscan_id.innerHTML = `${data.nscan_id}`;
					location = data.path + "/" + data.name
					open_file.onclick = open_item

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
		var name = document.getElementById("name")
		var path = document.getElementById("path")
		var date_modified = document.getElementById("date_modified")
		var size = document.getElementById("size")
		var ndirectory_id = document.getElementById("ndirectory_id")
		var nscan_id = document.getElementById("nscan_id")
		var open_file = document.getElementById("open-file")

		showFile();
	}
}