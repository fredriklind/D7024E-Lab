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

	var apiURL = "http://localhost:"+ getAPIport() +"/api/";
	$.support.cors = true

	// Get section
	$('#get').submit(function(e) {
		var keyField = $(this).find(".key"),
			section = keyField.parents('.section'),
			key = keyField.val();
		section.find("label[for=getvalue]").remove();
		keyField.parents('label').after('\
			<label for="getvalue">\
				<span class="text">Value:</span>\
				<span class="field">\
					<input tabindex="1" id="getvalue" class="value loading" type="text" name="Value">\
				</span>\
			</label>');

		$.ajax({
			url: apiURL + "storage/" + key,
			type: "GET",
			dataType: "text",
			success: function(data, textStatus, jqXHR) {
				section.find('#getvalue').val(data).removeClass('loading');
				keyField.bind("keydown", function() {
					keyField.unbind("keydown");
					section.find("label[for=getvalue]").remove();
				});
			},
			error: function(jqXHR, status, error) {
				section.find("label[for=getvalue]").remove();
				keyField.val("");
				keyField.focus();
				showStatus(section, "error", error);
			}
		});

		e.preventDefault();
		return false;
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