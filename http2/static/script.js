document.addEventListener("DOMContentLoaded", function () {
    console.log("JavaScript is loaded!")
    const message = document.createElement("p")
    message.textContent = "This is a message from the script.js file"
    document.body.appendChild(message)
})