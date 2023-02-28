const form = document.getElementById('login-form');

form.addEventListener('submit', function(event) {
	event.preventDefault();
	
	// sanity check
	if (!this.hostname.value) return false;
	if (!this.port.value) return false;
	if (!this.username.value) return false;
	if (!this.password.value) return false;
	
	const data = {
		hostname: this.hostname.value,
		port: this.port.value,
		username: this.username.value,
		password: this.password.value
	};
	
	postData(data).then (reply => {
	  if (reply['result'] == 0) {
		  var div = document.getElementById('message');
		  div.innerHTML = reply['msg'];
		  div.style.display = "block";
		  return;
	  }
	  location.href = '/index.php';
	});
});

async function postData(data) {
	try {
		const response = await fetch ('/login_action.php', {
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