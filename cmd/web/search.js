function checkConnection() {
    var request = new XMLHttpRequest();
    request.open('GET', config.api + '/ping');
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

function getParameterByName(name, url) {
    if (!url) url = window.location.href;
    name = name.replace(/[\[\]]/g, '\\$&');
    var regex = new RegExp('[?&]' + name + '(=([^&#]*)|&|#|$)'),
        results = regex.exec(url);
    if (!results) return null;
    if (!results[2]) return '';
    return decodeURIComponent(results[2].replace(/\+/g, ' '));
}

document.onreadystatechange = function () {
    var showAdvancedSearch = 0;

    function showHideAdvancedSearch() {
        if (showAdvancedSearch == 0) {
            advancedSeachPanel.style.display = "none";
        } else {
            advancedSeachPanel.style.display = "inline";
        }
    }

    function advancedSearchClick() {
        showAdvancedSearch = !showAdvancedSearch;
        showHideAdvancedSearch();
    }

    function searchText() {
        result.innerHTML = ''
        var query = document.querySelector('#search').value;

        if (showAdvancedSearch) {
            type = document.querySelector('#type').value;
            dateFrom = document.querySelector('#date-from').value;
            dateTo = document.querySelector('#date-to').value;
            includeDirectories = false;
            if (document.querySelector('#include-directories').value == 'on') {
                includeDirectories = true;
            }

            location = 'search.html?q=' + query + '&type=' + type + '&from=' + dateFrom + '&to=' + dateTo + '&id=' + includeDirectories;
        } else {
            location = 'search.html?q=' + query;
        }
    }

    function search(query, type, from, to, id) {
        result.innerHTML = ''
        var request = new XMLHttpRequest();

        if (showAdvancedSearch) {
            request.open('GET', config.api + '/search?q=' + query + '&type=' + type + '&after=' + from + '&before=' + to + '&id=' + id);
        } else {
            request.open('GET', config.api + '/search?q=' + query);
        }

        request.onreadystatechange = function () {
            if (request.readyState == 4) {
                if (request.status === 200) {
                    var data = JSON.parse(request.responseText);

                    content = '<table style="border-collapse: collapse;border-spacing: 0;">'
                    content += '<thead>'
                    content += '<tr>'
                    content += '<th>Name</th>'
                    content += '<th>Modified</th>'
                    content += '<th>Size</th>'
                    content += '<th>Location</th>'
                    content += '</tr>'
                    content += '</thead>'
                    content += '<tbody>'

                    // directories
                    if (data.directories != null) {
                        for (var i = 0; i < data.directories.length; i++) {
                            content += '<tr>'
                            content += `<td>[<a href="dir.html?id=${data.directories[i].id}">${data.directories[i].name}</a>]</td>`
                            content += `<td>${data.directories[i].date_modified}</td>`
                            content += `<td>${data.directories[i].size}</td>`
                            content += `<td>hello</td>`
                            content += '</tr>'
                        }
                    }

                    // files
                    if (data.files != null) {
                        for (var i = 0; i < data.files.length; i++) {
                            path = data.files[i].path + "/" + data.files[i].name;
                            path = path.replace(/'/g, "\\'");
                            content += '<tr>'
                            content += `<td><a href="javascript:window.open('${path}', '_blank');">${data.files[i].name}</a></td>`
                            content += `<td>${data.files[i].date_modified}</td>`
                            content += `<td>${data.files[i].size}</td>`
                            content += `<td><a href='dir.html?id=${data.files[i].ndirectory_id}'>${data.files[i].ndirectory_name}</a></td>`
                            content += '</tr>'
                        }
                    }

                    content += '</tbody>'
                    content += '</table>'

                    result.innerHTML = content;

                    var rc = 0
                    if (data.files != null) {
                        rc += data.files.length
                    }
                    if (data.directories != null) {
                        rc += data.directories.length
                    }

                    message.style.display = "block";
                    message.innerHTML = `<p>${rc} results for '${data.query}'</p>`;
                } else if (request.status === 404) {
                    message.style.display = "block";
                    message.innerHTML = "no results for '" + query + "'";
                } else {
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

        var message = document.querySelector('#message');
        var result = document.querySelector('#result');
        var searchForm = document.querySelector('#search-form');

        searchForm.addEventListener('submit', function (e) {
            e.preventDefault()
            searchText()
        });

        var advancedSeachPanel = document.querySelector('#advanced-search-panel');
        document.getElementById("advanced-search").addEventListener("click", advancedSearchClick, false);

        var q = getParameterByName('q');
        if (q != '') {
            document.querySelector('#search').value = q;
            var type = getParameterByName('type');
            var from = getParameterByName('from');
            var to = getParameterByName('to');
            var id = getParameterByName('id');
            if (type != null || from != null || to != null || id != null) {
                showAdvancedSearch = 1;

                document.querySelector('#type').value = type;
                document.querySelector('#date-from').value = from;
                document.querySelector('#date-to').value = to;
                if (id == "true") {
                    document.querySelector('#include-directories').checked = true;
                } else {
                    document.querySelector('#include-directories').checked = false;
                }

                includeDirectories = document.querySelector('#include-directories').value;

                showHideAdvancedSearch();
            }

            search(q, type, from, to, id);
        }
    }
}