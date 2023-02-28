const form = document.getElementById('upload-form');
const enableLogging = true;

const actionTypeUploadStart = 1
const actionTypeUploadEnd   = 2
const actionTypeUploadError = 3
const actionTypeUploadLoop  = 4

form.addEventListener('submit', function(event) {
	event.preventDefault();
	
	if (!this.datafile.value) return false;
	
	var f = document.getElementById('datafile');
	var k = document.getElementById('enckey');
	
	var file = f.files[0];
	var size = file.size;
	var key = (k.value && k.value.length) ? k.value : null;
	var sliceSize = 1024*1024;
	var start = 0;
	
	var result = uploadDatabase(f.value, key, file, size, sliceSize);
});
	
function uploadDatabase(path, key, file, size, chunkSize) {
	if (enableLogging) console.log("uploadDatabase: " + path + " key: " + key + " size: " + size + " chunk: " + chunkSize);
	
	if (uploadStart(path, key, size, chunkSize) == false) return false;
	if (uploadLoop(file, 0, size, chunkSize) == false) return false;
	return true;
}	

function uploadStart(file, key, size, chunkSize) {
	resetUI(actionTypeUploadStart);
	
	var name = file.split(/(\\|\/)/g).pop();
	if (enableLogging) console.log('uploadStart ' + name);
	
	var formdata = new FormData();
	formdata.append('action', 0);		// START
	formdata.append('name', name);
	if (key) formdata.append('key', key);
	
	var reqcount = Math.round((size / chunkSize) + 50);
	if (enableLogging) console.log("reqcount: " + reqcount);
	
	var xhr = new XMLHttpRequest();
	xhr.open('POST', '/upload_action.php', false);
	xhr.setRequestHeader("Connection", "keep-alive");
	xhr.setRequestHeader("Keep-Alive", "timeout=15, max=" + reqcount + "\"");
	xhr.send(formdata);
	
	return true;
}

function uploadEnd() {
	resetUI(actionTypeUploadEnd);
	
	if (enableLogging) console.log('uploadEnd');
	
	var formdata = new FormData();
	formdata.append('action', 2);		// END
	
	var xhr = new XMLHttpRequest();
	xhr.open('POST', '/upload_action.php', false);
	xhr.send(formdata);
	
	progressSet(100);
	displayMessage("Database succesfully uploaded.");
	
	setTimeout(() => {location.href = '/databases.php';}, 500);
	return true;
}

function uploadAbort() {
	resetUI(actionTypeUploadError);
	
	if (enableLogging) console.log('uploadAbort');
	
	var formdata = new FormData();
	formdata.append('action', 666);		// ABORT
	
	var xhr = new XMLHttpRequest();
	xhr.open('POST', '/upload_action.php', false);
	xhr.send(formdata);
	
	return true;
}

function uploadLoop (file, start, end, size) {
	// compute local values
	var islast = false;
	var len = size;
	if (start + len > end) {len = end - size; islast = true;}
	if (len < 0) len = end;
	
	// send next/final chunk ONLY after the previous one has been sent
	var xhr = new XMLHttpRequest();
	xhr.onreadystatechange = function () {
		if (xhr.readyState == XMLHttpRequest.DONE) {
			if (this.responseText != 0) {
				uploadAbort();
				displayError(this.responseText);
				return;
			}
			
			var value = Math.floor(( start / end) * 100);
			progressSet(value);
			(islast) ? uploadEnd() : uploadLoop(file, start + size, end, size);
		}
	}
	
	// compute chunk to send
	var chunk = slice(file, start, start + len);
	
	// chunk is now a Blob that can only be read async
	const reader = new FileReader();
	reader.onloadend = function () {
		// prepare parameters
		if (enableLogging) console.log('uploadChunk: ' + start + ' ' + len);
		   
		var formdata = new FormData();
		formdata.append('action', 1);					// UPLOAD
		formdata.append('start', start);
		formdata.append('len', len);
		formdata.append('end', end);
		formdata.append('chunk', reader.result);		// reader.result contains the contents of blob
		formdata.append('encoding', 1);					// 1 -> base64, 2 -> binary
		
		// post data
		xhr.open('POST', '/upload_action.php', true);
		xhr.send(formdata);
	}
	reader.readAsDataURL(chunk);
	
	return true;
}

function slice(file, start, end) {
	var slice = file.mozSlice ? file.mozSlice :
				file.webkitSlice ? file.webkitSlice :
  			  	file.slice ? file.slice : noop;
	return slice.bind(file)(start, end);
}

function noop() {
}

function displayError(msg) {
	var errorID = document.getElementById('upload_error');
	errorID.innerHTML = msg;
	errorID.style.display = "block";
}

function displayMessage(msg) {
	var messageID = document.getElementById('upload_message');
	messageID.innerHTML = '<strong>' + msg + '</strong>';
	messageID.style.display = "block";
}

function resetUI(actionType) {
	var errorID = document.getElementById('upload_error');
	var messageID = document.getElementById('upload_message');
	var progressID = document.getElementById('upload_progress');
	var buttonID = document.getElementById('upload_button');
	
	var showProgress = (actionType == actionTypeUploadStart);
	var hideProgrss = (actionType == actionTypeUploadEnd);
	var hideError = (actionType == actionTypeUploadStart);
	var enableButton = (actionType == actionTypeUploadEnd) || (actionType == actionTypeUploadError);
	var disableButton = (actionType == actionTypeUploadStart);
	var hideMessage = (actionType == actionTypeUploadStart);
	
	if (showProgress) progressID.style.display = "block";
	if (hideProgrss) progressID.style.display = "none";
	if (hideError) errorID.style.display = "none";
	if (hideMessage) messageID.style.display = "none";
	if (enableButton) buttonID.disabled = false;
	if (disableButton) buttonID.disabled = true;
}

function progressSet(value) {
	var progressID = document.getElementById('upload_progress');
	progressID.value = value;
}
