function httpRequest(mode, url, data, fnc) {

    var jsonObj = {
        "Data": data
    };

    data = mode == 'GET' ? data : JSON.stringify(jsonObj);
    dataType = mode == 'GET' ? 'text' : 'json';
    contentType = mode == 'GET' ? 'text/plain' : 'application/json; charset=utf-8';

    $.ajax({
        type: mode,
        url: window.location.origin + url,
        data: data,
        dataType: dataType,
        contentType: contentType,
        success: fnc
    });
}
