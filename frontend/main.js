function shorten(){
    fetch("http://localhost:6969/shorten", {
        method: "POST",
        body: JSON.stringify({
            "longurl": "google.com"
        }),
        headers: {
            "Content-Type": "application/json"
        }
    })
    .then((response) => response.json())
    .then((json) => console.log(json));
} 

