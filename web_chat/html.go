package main

import "html/template"

var (
	rootTemplate = template.Must(template.New("root").Parse(`
<!DOCTYPE html>
<html>
<head>
<ul id='chat'>
</ul>
<form>
<textarea rows="8" cols="80" id="message"></textarea>
  <br>
  <button type="submit">Send</button>
</form>
<meta charset="utf-8" />
<script>
    websocket = new WebSocket("ws://{{.}}/socket");
    websocket.onmessage = function(m) {
		console.log("Received:", m.data);
		let li = document.createElement('li');
		li.innerText = m.data;
		document.querySelector('#chat').append(li);
}

    
    document.querySelector('form').addEventListener('submit', (event) => {
        event.preventDefault();
        let message = document.querySelector('#message').value;
	    websocket.send(message);
        document.querySelector('#message').value = '';
    });
</script>
</html>
`))
)
