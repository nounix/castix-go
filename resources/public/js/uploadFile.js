function uploadFile() {
	var url = window.location;
	var fileName = $("#fileUpload").val().split('\\')[2];

    $.ajax({
    	type : "POST",
		url : url + "/store-file-async",
		data : "fileName=" + fileName,
    	success: function(){
    		var form = new FormData();
    		var file = $("#fileUpload")[0].files[0];
    		form.append('file', file);
	        $.ajax({
		    	type : "POST",
				url : url + "/upload-file",
				dataType: 'json',
		        data: form,
		        processData: false,
		        contentType: false
		    });
	    }
    });
    return false;	//to prevent submit
}
