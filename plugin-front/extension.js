(async function () {
    if (window.location.href.includes('/file/')) {
        return 
    }

    const API = 'http://localhost:8080'
    const rid = Math.random().toString(36).substring(2, 10)

    const button = document.createElement('button')
    button.textContent = '点击复制分享连接'
    button.style.display = 'none';
    button.id = "cwButton"
    button.style.backgroundColor = 'green';
    button.style.color = 'black';

    button.addEventListener('click', () => {
        const textToCopy = API + "/file/" + rid
        navigator.clipboard.writeText(textToCopy)
    })

    document.body.insertBefore(button, document.body.firstChild);

    window.addEventListener('message', async e => {
        if (e.data.type == "set-file") {
            console.log('file', e.data.data)
            fetch(API + "/file", {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    rid: rid,
                    ...e.data.data,
                }),
            })
                .then(response => response.json())
                .then()
                .catch(error => console.error('出错了:', error));

            button.style.display = 'block'
        }
    })


    setInterval(() => {
        if (button.style.display === 'block') {
            const video = document.querySelector('video')
            current = video.currentTime
            fetch(API + "/progress/" + rid, {
                method: 'PUT',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    currentTime: String(current),
                    unixMill: String(Date.now()),
                }),
            })
                .then(response => response.json())
                .then()
                .catch(error => console.error('出错了:', error));
        }
    }, 2000)
})()