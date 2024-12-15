document.getElementById("send").addEventListener("click", async function () {
    const username = document.getElementById("username").value;
    const password = document.getElementById("password").value;

    try {
        await fetch("/api/v1/login", {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({ username, password }),
            credentials: "same-origin",
        })
            .then((response) => {
                console.log(response.status);
                return response.text();
            })
            .then((data) => console.log(data));
    } catch (error) {
        console.error("Error:", error);
        document.getElementById("message").textContent =
            "An error occurred during login.";
    }
});
