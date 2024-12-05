document.getElementById("send").addEventListener('click', async function() {
    const username = document.getElementById('username').value;
    const password = document.getElementById('password').value;

    try {
        const response = await fetch('http://127.0.0.1:6969/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ username,password }),
            credentials: 'same-origin'
        })

        if (response.ok) {
            alert("LOGIN OK");
            console.log(response.json());
        } else {
            alert("LOGIN NOK!");
        }
    } catch (error) {
        console.error('Error:', error);
        document.getElementById('message').textContent = 'An error occurred during login.';
    }
});
