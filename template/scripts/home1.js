$(document).ready(function(){

    DEFAULT_COOKIE_EXPIRE_TIME = 30;

    uname = '';
    session = '';
    uid = 0;
    currentVideo = null;
    listedVideos = null;

    session = getCookie('session'); //从cookie中读取session
    uname = getCookie('username'); //从cookie中读取username



    listAllVideos(function(res,err){
        if (err != null) {
            //window.alert('encounter an error, pls check your username or pwd');
            popupErrorMsg('encounter an error, pls check your username or pwd');
            return;
        }
        var obj = JSON.parse(res);
        listedVideos = obj['videos'];
        obj['videos'].forEach(function (item, index) {
            console.log(item['id'], item['name'], item['display_ctime']);
            var ele = listAllElement(item['id'], item['name'], item['display_ctime']);  //为每个视频创建一个div
            $("#whole").append(ele);
        });

    });
})


function setCookie(cname, cvalue, exmin) {
    var d = new Date();
    d.setTime(d.getTime() + (exmin * 60 * 1000));
    var expires = "expires=" + d.toUTCString();
    document.cookie = cname + "=" + cvalue + ";" + expires + ";path=/";
}

function getCookie(cname) {
    var name = cname + "="; //cname  sessionId or username
    var ca = document.cookie.split(';');
    for (var i = 0; i < ca.length; i++) {
        var c = ca[i];
        while (c.charAt(0) == ' ') {
            c = c.substring(1);
        }
        if (c.indexOf(name) == 0) {
            return c.substring(name.length, c.length);
        }
    }
    return "";
}

function listAllVideos(callback) {
    var dat = {
        'url': 'http://' + window.location.hostname + ':8000/user/' + uname + '/videos',
        'method': 'GET',
        'req_body': ''
    };

    $.ajax({
        url: 'http://' + window.location.hostname + ':8080/api',
        type: 'post',
        data: JSON.stringify(dat),
        headers: {
            'X-Session-Id': session
        },
        statusCode: {
            500: function () {
                callback(null, "Internal error");
            }
        },
        complete: function (xhr, textStatus) {
            if (xhr.status >= 400) {
                callback(null, "Error of Signin");
                return;
            }
        }
    }).done(function (data, statusText, xhr) {
        if (xhr.status >= 400) {
            callback(null, "Error of Signin");
            return;
        }
        callback(data, null);
    });
}

function listAllElement(vid, name, ctime) {
    var myLink = $('<a/>', {
        href: "http://47.106.252.71:8080/userhome"
    })

    var myVideoItem = $('<div/>', {
        class: "video-item",
        id: vid
    });

    myVideoItem.append($('<img/>', {
        src: "../template/img/preloader.jpg"
    }));
    myVideoItem.append($('<div/>', {
        text: name,
        class: "displayName"
    }));
    myVideoItem.append($('<div/>', {
        text: ctime,
        class: "displayCtime"
    }));

    myLink.append(myVideoItem);
    return myLink;
}
