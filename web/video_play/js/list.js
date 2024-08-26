let videoData = [
    {
        author: "xigua",
        title: "西瓜播放器宣传视频",
        src: "//sf1-cdn-tos.huoshanstatic.com/obj/media-fe/xgplayer_doc_video/mp4/xgplayer-demo-360p.mp4",
        resource: [
            {
                name: "1080p",
                url: "//sf1-cdn-tos.huoshanstatic.com/obj/media-fe/xgplayer_doc_video/mp4/xgplayer-demo-720p.mp4"
            },
            {
                name: "720p",
                url: "//sf1-cdn-tos.huoshanstatic.com/obj/media-fe/xgplayer_doc_video/mp4/xgplayer-demo-480p.mp4"
            },
            {
                name: "480p",
                url: "//sf1-cdn-tos.huoshanstatic.com/obj/media-fe/xgplayer_doc_video/mp4/xgplayer-demo-360p.mp4"
            }
        ]
    }
];

async function getVideoData() {
    try {
        const response = await fetch('http://127.0.0.1:8080/video/list'); // 替换为实际的 API
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        videoData = await response.json();
    } catch (error) {
        console.error('Fetch error:', error);
    }
}

async function deleteVideo(videoName) {
    try {
        const response = await fetch('http://127.0.0.1:8080/video/delete', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ name: videoName })
        });
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        console.log(`${videoName} have been deleted.`);
    } catch (error) {
        console.error('Delete error:', error);
    }
}


let lastSelectedDiv = null;

async function init() {
    const selector = document.getElementById('selector');
    const playerContainer = document.getElementsByClassName('fakevideo');

    await getVideoData()
    console.log("videoData:", videoData) // 打印videoData

    let config = {
        id: 'video-container',
        autoplay: false,
        volume: 0.3,
        url: '//sf1-cdn-tos.huoshanstatic.com/obj/media-fe/xgplayer_doc_video/mp4/xgplayer-demo-360p.mp4',
        poster: "//lf9-cdn-tos.bytecdntp.com/cdn/expire-1-M/byted-player-videos/1.0.0/poster.jpg",
        playsinline: true,
        height: playerContainer.clientHeight,
        width: playerContainer.clientWidth,
        download: true //设置download控件显示
    }

    let player = new Player(config);

    // player.emit('resourceReady', videoData[0].resource);

    videoData.forEach((data, index) => {
        const videoDiv = document.createElement('div');
        videoDiv.id = index + 1;

        const videoElement = document.createElement('video');
        videoElement.src = data.src;
        const videoNameDiv = document.createElement('div');
        videoNameDiv.className = 'video-name';
        videoNameDiv.textContent = data.title;

        videoDiv.appendChild(videoElement);
        videoDiv.appendChild(videoNameDiv);

        videoDiv.addEventListener('click', () => {
            player.destroy()

            config.url = videoElement.src
            config.poster = ""
            config.autoplay = true
            player = new Player(config);
            player.emit('resourceReady', data.resource)

            if (lastSelectedDiv) {
                lastSelectedDiv.style.backgroundColor = 'transparent';
            }
            videoDiv.style.backgroundColor = '#ab0';
            // playerName.textContent = data.title;
            // playerAuthor.textContent = data.author;
            lastSelectedDiv = videoDiv;
        });

        videoDiv.addEventListener('mouseover', () => {
            if (videoDiv !== lastSelectedDiv) {
                videoDiv.style.backgroundColor = 'cadetblue';
            }
        });

        videoDiv.addEventListener('mouseleave', () => {
            if (videoDiv !== lastSelectedDiv) {
                videoDiv.style.backgroundColor = 'transparent';
            }
        });

        // Handle Delete Key Press
        document.addEventListener('keydown', async (event) => {
            if (event.key === 'Backspace' && lastSelectedDiv === videoDiv) {
                await deleteVideo(data.title);
                selector.removeChild(videoDiv);
                lastSelectedDiv = null;
            }
        });

        selector.appendChild(videoDiv);
    });
}


init()
