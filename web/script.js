const prompts = [
    "what's your jam? 🍓",
    "please no country",
    "if you're reading this, it's beyonce time",
    "sf summer songz only",
    "kamala IS brat",
    "disclaimer: djvivo will likely not play your song",
    "oh hey, it's time for another shot!",
    "who is djvivo really?",
    "hey babe, you're looking cute today ;)",
    "no, this is definitely not a sign to request chappell roan",
    "i'm feeling ... THROWBACKS",
    "despacito? despacito.",
    "hit me with your best song request",
    "it's a beautiful day to hang in the park and listen to _______",
    "if you were about to request calvin harris, you came to the wrong place",
    "for the girls songs only"
]

const emojis = ['🫐', '🌞', '😮‍💨', '✨', '🌼', '🐳', '🎾', '🌺', '🥭', '🍭', '🍯', '🍸', '🧊', '🧃', '🥥', '🌻', '🌸', '🌍', '🎾', '🥥', '🌝', '🍌', '🐣', '🐙', '🌀', '🥥', '🍹', '🌴', '♻️', '🎈', '🔫', '🚬', '💓', '🍻', '💞', '🥥', '💙', '💜', '💦', '🗯️', '🦁', '🐸', '🦜', '🐬', '🌱', '⚠️', '💰', '🦋', '🎾', '🍐', '🍒', '🍑', '🍉', '🍧', '🎾']


function songPrompt() {
    const picked = prompts[Math.floor(Math.random() * prompts.length)]
    document.getElementById("songprompt").innerText = picked
}

function randomEmoji() {
    return Math.floor(Math.random() * emojis.length)
}


async function requestSong() {
    const song = document.getElementById("song").value
    const req = document.getElementById("who").value
    if (!song || !req) {
        alert("please fill out both fields!")
        return
    }
    data = {
        emoji: randomEmoji(),
        song: song,
        requester: req,
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
        alert(resp)
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
        let emoji = emojis[3]
        if (song.emoji !== 0) {
            emoji = emojis[song.emoji]
        }
        item.innerHTML = `${emoji} <b>${song.name}</b> requested by ${song.requester}`
        ul.appendChild(item)
        ul.appendChild(document.createElement("br"))
    }
}

async function loadUp() {
    songPrompt()
    await getSongRequests()
}