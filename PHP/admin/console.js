const form = document.getElementById('console-form');

form.addEventListener('submit', function(event) {
	event.preventDefault();
	
	// sanity check
	if (!this.database.value) return false;
	if (!this.sql.value) return false;
	
	const data = {
		database: this.database.value,
		sql: this.sql.value
	};
	
	console.log(this.database.value);
	console.log(this.sql.value);
	
	postData(data).then (reply => {
		console.log(reply);
		
		var errorID = document.getElementById('console-error');
		var messageID = document.getElementById('console-message');
		var tableID = document.getElementById('console-table');
		
		// hide all
		errorID.style.display = "none";
		messageID.style.display = "none";
		tableID.style.display = "none";
		
		if (reply['result'] == 0) {
			errorID.innerHTML = reply['msg'];
			errorID.style.display = "block";
			return;
		}
		
		if (reply['result'] == 1) {
			messageID.innerHTML = reply['msg'];
			messageID.style.display = "block";
			return;
		}
		
		// else display table
		if (reply['result'] == 2) {
			tableID.innerHTML = reply['msg'];
			tableID.style.display = "block";
		}
		
	});
});

async function postData(data) {
	try {
		const response = await fetch ('/console_action.php', {
				method: 'POST',
				body: JSON.stringify(data)
			}
		);
		
		const reply = await response.json();
		return reply;
		
	} catch (error) {
		const reply = {
			result: 0,
			msg: error	
		};
		return JSON.stringify(reply);
	}
}