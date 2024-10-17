var navLinks = document.querySelectorAll("nav a");
for (var i = 0; i < navLinks.length; i++) {
	var link = navLinks[i]
	if (link.getAttribute('href') == window.location.pathname) {
		link.classList.add("live");
		break;
	}
}

var messageHeader = document.getElementById("messages")

const socket = new WebSocket('ws://localhost:4000/ws')
socket.onmessage = (event) => {
    messageHeader.innerHTML = event.data;
}
