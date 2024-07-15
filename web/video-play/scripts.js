let  videoData = [
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
        const response = await fetch('http://127.0.0.1:8080/video/list'); // 替换为实际的 API 端点
        if (!response.ok) {
            throw new Error('Network response was not ok ' + response.statusText);
        }
        videoData = await response.json();
    } catch (error) {
        console.error('Fetch error:', error);
    }
}


let lastSelectedDiv = null;

async function init() {
    const selector = document.getElementById('selector');
    const playerContainer = document.getElementById('video-container');
    const playerName = document.getElementById('name');
    const playerAuthor = document.getElementById('author');

    await getVideoData()
    console.log(videoData)
    let player = new Player({
        id: 'video-container',
        autoplay: false,
        volume: 0.3,
        url: '//sf1-cdn-tos.huoshanstatic.com/obj/media-fe/xgplayer_doc_video/mp4/xgplayer-demo-360p.mp4',
        poster: "//lf9-cdn-tos.bytecdntp.com/cdn/expire-1-M/byted-player-videos/1.0.0/poster.jpg",
        playsinline: true,
        height: playerContainer.clientHeight,
        width: playerContainer.clientWidth,
    });


    player.emit('resourceReady', videoData[0].resource);

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
            player.src = videoElement.src;
            if (lastSelectedDiv) {
                lastSelectedDiv.style.backgroundColor = 'transparent';
            }
            videoDiv.style.backgroundColor = '#808080';
            playerName.textContent = data.title;
            playerAuthor.textContent = data.author;
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

        selector.appendChild(videoDiv);
    });
}


init()

let tab_list = document.querySelector('#tab');
let lis = tab_list.querySelectorAll('li');
let items = document.querySelectorAll('.content');

// for循环绑定事件
for (var i = 0; i < lis.length; i++) {
    // 开始给五个li设置索引号
    lis[i].setAttribute('index', i);
    lis[i].onclick = function () {
        // 干掉所有人 其余的li清除class这个类
        for (var i = 0; i < lis.length; i++) {
            lis[i].className = '';
        }
        // 留下我自己
        this.className = 'current';
        // 2.下面的显示内容模块
        var index = this.getAttribute('index');
        console.log(index);
        // 干掉所有人，让其余的item这些div隐藏
        for (var i = 0; i < items.length; i++) {
            items[i].style.display = 'none';
        }
        items[index].style.display = 'block';
    }
}


document.getElementById('submitFile').addEventListener('click', async () => {
    const fileInput = document.getElementById('fileInput');
    const submitResult = document.getElementById('submitResult');

    if (fileInput.files.length === 0) {
        submitResult.textContent = '请先选择一个文件';
        return;
    }

    const formData = new FormData();
    formData.append('file', fileInput.files[0]);

    try {
        const response = await fetch('http://127.0.0.1:8080/video/upload', { // 替换为实际的上传接口
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            throw new Error('网络响应错误: ' + response.statusText);
        }

        const result = await response.json();
        console.log(result)
        submitResult.textContent = '文件上传成功: ' + result.message;
    } catch (error) {
        console.error('文件上传错误:', error);
        submitResult.textContent = '文件上传失败: ' + error.message;
    }
});



