document.onreadystatechange = function () {

	function getParameterByName(name, url) {
		if (!url) url = window.location.href;
		name = name.replace(/[\[\]]/g, '\\$&');
		var regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)'),
			results = regex.exec(url);
		if (!results) return null;
		if (!results[2]) return '';
		return decodeURIComponent(results[2].replace(/\+/g, ' '));
	}

	function showDir() {
		result.innerHTML = ''
		var id = getParameterByName('id');

		var request = new XMLHttpRequest();
		request.open('GET', 'http://localhost:5000/directories/' + id);
		request.onreadystatechange = function () {
			if (request.readyState == 4) {
				if (request.status === 200) {
					var data = JSON.parse(request.responseText);

					dname.innerHTML = data.name;
					parent.innerHTML = "<a href='dir.html?id=" + data.parent_id + "'>Parent directory</a>"

					content = '<table style="border-collapse: collapse;border-spacing: 0;">'
					content += '<thead>'
					content += '<tr>'
					content += '<th>Name</th>'
					content += '<th>Modified</th>'
					content += '<th>Size</th>'
					content += '</tr>'
					content += '</thead>'
					content += '<tbody>'

					if (data.directories != null) {
						for (var i = 0; i < data.directories.length; i++) {
							content += '<tr>'
							content += `<td>[<a href="dir.html?id=${data.directories[i].id}">${data.directories[i].name}</a>]</td>`
							content += `<td></td>`
							content += `<td>${data.directories[i].size}</td>`
							content += '</tr>'
						}
					}

					if (data.files != null) {
						for (var i = 0; i < data.files.length; i++) {
							content += '<tr>'
							content += `<td><a href="vip.html?id=${data.files[i].id}">${data.files[i].name}</a></td>`
							content += `<td>${data.files[i].date_modified}</td>`
							content += `<td>${data.files[i].size}</td>`
							content += '</tr>'
						}
					}

					content += '</tbody>'
					content += '</table>'
					result.innerHTML = content;
					result.style.display = "block"

					if (data.files == null && data.directories == null) {
						result.style.display = "none";
						message.style.display = "block";
						message.innerHTML = `no results for '${data.name}'`;
					}

					if (data.files != null && data.directories != null) {
						message.style.display = "block";
						message.innerHTML = 'total directories: ' + data.directories.length + ', total files: ' + data.files.length;
					} else if (data.files != null) {
						message.style.display = "block";
						message.innerHTML = 'total files: ' + data.files.length;					
					} else if (data.directories != null) {
						message.style.display = "block";
						message.innerHTML = 'total directories: ' + data.directories.length;					
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
		var result = document.querySelector('#result');
		var dname = document.querySelector('#name');
		var parent = document.querySelector('#parent');

		showDir()
	}
}