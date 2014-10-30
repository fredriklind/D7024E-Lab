$(function() {

	$('.section h2').click(function() {
		showSection($(this).parent('.section'));
	});

	$('input[type=text]').focus(function(e) {
		var par = $(this).parents('.section');

		if (par.hasClass('hidden')) {
			showSection(par);
		}
	});

	var apiURL = "http://localhost:"+ getAPIport() +"/nodes";
	$.support.cors = true

	$.ajax({
			url: apiURL + "/",
			type: "GET",
			dataType: "application/json",
			success: function(data, textStatus, jqXHR) {
				console.log(data)
			},
			error: function(jqXHR, status, error) {
				console.log(error)
			}
		});

	// Set section
	$('#set').submit(function(e) {
		var keyField = $(this).find(".key"),
			valueField = $(this).find(".value"),
			section = keyField.parents('.section'),
			key = keyField.val();
		section.removeClass("set-successful");
		dataString = JSON.stringify({
			Key: keyField.val(),
			Value: valueField.val()
		})

		$.ajax({
			url: apiURL + "storage/",
			type: "POST",
			dataType: "text",
			data: dataString,
			success: function(data, textStatus, jqXHR) {
				keyField.val("");
				valueField.val("");
				showStatus(section, "success", "Set successful");
			},
			error: function(jqXHR, status, error) {
				showStatus(section, "error", error);
			}
		});

		e.preventDefault();
		return false;
	});


	// Update section
	$('#update').submit(function(e) {
		var keyField = $(this).find(".key"),
			valueField = $(this).find(".value"),
			section = keyField.parents('.section'),
			key = keyField.val();

		$.ajax({
			url: apiURL + "storage/" + key,
			type: "PUT",
			dataType: "text",
			data: valueField.val(),
			success: function(data, textStatus, jqXHR) {
				keyField.val("");
				valueField.val("");
				showStatus(section, "success", "Update successful");
			},
			error: function(jqXHR, status, error) {
				keyField.focus()
				keyField.val("");
				showStatus(section, "error", error);
			}
		});

		e.preventDefault();
		return false;
	});

	// Delete section
	$('#delete').submit(function(e) {
		var keyField = $(this).find(".key"),
			section = keyField.parents('.section'),
			key = keyField.val();

		$.ajax({
			url: apiURL + "storage/" + key,
			type: "DELETE",
			success: function(data, textStatus, jqXHR) {
				keyField.focus();
				keyField.val("");
				showStatus(section, "success", "Delete successful");
			},
			error: function(jqXHR, status, error) {
				keyField.focus()
				keyField.val("");
				showStatus(section, "error", error);
			}
		});

		e.preventDefault();
		return false;
	});

});

function showSection(t) {
	$('.section').addClass('hidden');
	t.removeClass('hidden');
	t.find('.key').focus();
}

function showStatus($section, status, text) {
	var statusElement = $section.find("h2 > span");
	statusElement.attr("class", "");
	statusElement.addClass(status);
	statusElement.html(text);
	statusElement.addClass("show");

	setTimeout(function() {
		statusElement.removeClass("show");
	}, 1200);
}

function getAPIport() {
	var port;
	$.ajax({
		url: "port.txt",
		success: function(data) { port = data; },
		async: false
	});
	return port;
}