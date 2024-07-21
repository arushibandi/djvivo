
async function requestSong() {
    console.log("yo")
    data = {
        song: document.getElementById("song").value,
        requester: document.getElementById("who").value,
    }
    console.log("sending", data)
    const resp = await fetch("/request", {
        method: "POST",
        body: JSON.stringify(data)
    })
    if (resp.status != 200) {
        console.error("got error from /request", resp)
    }
    await getSongRequests()
    document.getElementById("song").value = ""
    document.getElementById("who").value = ""
}

async function getSongRequests() {
    const resp = await fetch("/queue", {
        method: "GET",
    })
    if (resp.status != 200) {
        console.error("got error from /queue", resp)
        return
    }
    const body = await resp.json()
    console.log("got body", body)
    const ul = document.getElementById("queue")
    ul.innerHTML = ""
    for (let i = 0; i < body.songs.length; i++) {
        const song = body.songs[i]
        const item = document.createElement("li")
        item.style.listStyleType = "none"
        item.textContent = `${song.name} requested by ${song.requester}`
        ul.appendChild(item)
    }
}