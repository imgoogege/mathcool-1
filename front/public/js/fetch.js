var getDataFetch = (value,url)=> {
    fetch(url,{
        "Access-Control-Allow-Origin":"https://nudao.xyz",
    })
        .then(res => res.json())
        .then(function (b) {
            value = b
        })
        .catch(e => console.log("获取fetch数据出现错误：",e))
}



var postDataFetch = (url,objectData) => {
    fetch(url, {
        credentials: 'same-origin',
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(objectData)
    })
}
// js 获取query
function GetUrlParam(paraName) {
    var url = document.location.toString();
    var arrObj = url.split("?");

    if (arrObj.length > 1) {
        var arrPara = arrObj[1].split("&");
        var arr;

        for (var i = 0; i < arrPara.length; i++) {
            arr = arrPara[i].split("=");

            if (arr != null && arr[0] == paraName) {
                return arr[1];
            }
        }
        return "";
    }
    else {
        return "";
    }
}

// 模拟a的post请求
function doPost(to, p) { // to:提交动作（action）,p:参数
    var myForm = document.createElement("form");
    myForm.method = "post";
    myForm.action = to;
    for (var i in p){
        var myInput = document.createElement("input");
        myInput.setAttribute("name", i); // 为input对象设置name
        myInput.setAttribute("value", p[i]); // 为input对象设置value
        myForm.appendChild(myInput);
    }
    document.body.appendChild(myForm);
    myForm.submit();
    document.body.removeChild(myForm); // 提交后移除创建的form
}

