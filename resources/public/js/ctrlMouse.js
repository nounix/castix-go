var isSending = false;
var x1 = 0, x2 = 0, y1 = 0, y2 = 0;

function sendPoint(x1, x2, y1, y2) {
	var points = {
		"X1": parseInt(x1),
		"X2": parseInt(x2),
		"Y1": parseInt(y1),
		"Y2": parseInt(y2)
    };

	isSending = true;

    $.ajax({
    	type : "POST",
		url : window.location.origin + "/chrome/xdotool/mousemove",
		data : JSON.stringify(points),
		dataType: "json",
        contentType: 'application/json; charset=utf-8',
		success: function(){
	        isSending = false;
	    }
    });
    return false;
}

$.fn.ctrlMouse = function() {
	var start = function(e) {
		e.preventDefault();
        e = e.originalEvent.changedTouches ? e.originalEvent.changedTouches[0] : e;
		x1 = e.pageX;
		y1 = e.pageY;
	};
	var stop = function(e) {
		e = e.originalEvent.changedTouches ? e.originalEvent.changedTouches[0] : e;
		x2 = e.pageX;
		y2 = e.pageY;
		!isSending && sendPoint(x1, x2, y1, y2);
	};
	$(this).on("touchstart", start);
	$(this).on("touchend", stop);
	$(this).on("mousedown", start);
	$(this).on("mouseup", stop);
};

$(document).ready(function () {
// setup a new canvas for drawing wait for device init
    setTimeout(function(){
		$("#mousePad").ctrlMouse();
    }, 1000);
});
