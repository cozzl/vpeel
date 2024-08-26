const templates = {}; // 用来存储所有的模板数据

document.getElementById('addParam').addEventListener('click', function () {
    const form = document.getElementById('transParam');
    const formData = new FormData(form);
    const transParamsJson = {};

    formData.forEach((value, key) => {
        transParamsJson[key] = value;

        // 根据具体的需求解析不同的数据类型
        if (key === "bitrate" || key === "width" || key === "height" || key === "fps" || key === "gop" || key === "bframes" || key === "thread") {
            transParamsJson[key] = parseInt(value, 10); // 解析为整数
        } else {
            transParamsJson[key] = value; // 默认处理为字符串  
        }
    });

    // 获取模板名
    const templateName = document.getElementById('templateName').value;
    if (!templateName) {
        alert('请输入模板名');
        return;
    }

    if (templates[templateName] != null) {
        alert('模板已存在');
        return;
    }

    // 将 transParamsJson 以 templateName 为键存入 templates 对象
    templates[templateName] = transParamsJson;

    // 将JSON数据添加到表格中
    const tableBody = document.getElementById('paramTable').getElementsByTagName('tbody')[0];

    const row = document.createElement('tr');

    row.setAttribute('data-template-name', templateName); // 设置一个自定义属性以便删除

    const templateNameCell = document.createElement('td');
    templateNameCell.textContent = templateName;

    row.appendChild(templateNameCell);

    const keys = ["vcodec", "acodec", "width", "height", "bitrate", "fps", "gop", "bframes", "filter", "thread", "codec_param", "profile", "preset"];
    keys.forEach(key => {
        const cell = document.createElement('td');
        cell.textContent = transParamsJson[key] || ''; // 如果表单中没有该字段，显示为空
        row.appendChild(cell);
    });

    tableBody.appendChild(row);

    // 打印JSON对象到控制台
    console.log('Form Data as JSON:', JSON.stringify(transParamsJson));

    // 动态设置 overflow 样式
    const tableContainer = document.getElementById('tableContainer');
});

document.getElementById('delParam').addEventListener('click', function () {
    const templateName = document.getElementById('delTemplateName').value;
    if (!templateName) {
        alert('请输入模板名');
        return;
    }

    const row = document.querySelector(`[data-template-name="${templateName}"]`);
    if (row) {
        deleteTemplate(templateName, row);
    } else {
        alert('模板不存在');
    }
});

function deleteTemplate(templateName, row) {
    delete templates[templateName]; // 从模板对象中删除
    row.remove(); // 从表格中删除行
    console.log(`模板 ${templateName} 已删除`);
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

    // 将JSON字符串添加到FormData
    formData.append('params', JSON.stringify(templates));

    try {
        const response = await fetch('http://127.0.0.1:8080/video/uploadAndTrans', { // 替换为实际的上传接口
            method: 'POST',
            body: formData
        });

        if (!response.ok) {
            throw new Error('网络响应错误: ' + response.statusText);
        }

        const result = await response.json();
        submitResult.textContent = '文件上传成功!!!';
    } catch (error) {
        submitResult.textContent = '文件上传失败: ' + error.message;
    }
});