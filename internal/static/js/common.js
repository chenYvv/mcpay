layui.config({
    base: '/assets/layui'
})

layui.extend({
    treeTable: '/plugins/treeTable'
})

layui.use(['jquery', 'layer'], function () {
    let $ = layui.jquery;
    let layer = layui.layer;
    layer.config({
        move: false
    });
    $('body')
        .on('click', '[back]', function () {
            window.location.href = "javascript:history.go(-1)";
        })
        .on('click', '[reset]', function () {
            $(this).parents('form')[0].reset();
            $(this).parents('form').find('[search]').click();
        })
        .on('click', '[refresh]', function () {
            location.reload();
            return false;
        })
        .on('click', '[go]', function () {
            window.location.href = $(this).attr('go');
        })
        .on('click', '[noselect]', function () {
            $(this).parent().removeClass('layui-this');
        })
        .on('click', '[window]', function () {
            if ($(this).attr('window') && $(this).attr('title')) {
                iframe($(this).attr('window'), $(this).attr('title'))
                return false;
            }
        })

    showMsg(layer);
});

function iframe(url, name) {
    let id = new Date().getTime();
    parent.tools_element.tabAdd('layui-body-tab', {
        id: id,
        title: name,
        content: '<iframe width="100%" height="100%" frameborder="0" src="' + url + '"></iframe>'
    });
    parent.tools_element.tabChange('layui-body-tab', id);
}

function inArray(value, array) {
    let i = array.length;
    while (i--) {
        if (array[i] === value) {
            return true;
        }
    }
    return false;
}

function timeFormat(col) {
    return timeFormat_(col[col.LAY_COL.field])
}

function timeFormat_(time) {
    if (time === '0001-01-01T00:00:00Z' || time === null) {
        return '';
    }
    return new Date(+new Date(time) + 8 * 3600 * 1000).toISOString().replace(/T/g, ' ').replace(/\.[\d]{3}Z/, '');
}

function unixFormat(col) {
    return unixFormat_(col[col.LAY_COL.field])
}

function unixFormat_(unix_time) {
    unix_time = parseInt(unix_time) * 1000
    if (isNaN(unix_time) || unix_time <= 0) {
        return '';
    }
    let time = new Date(unix_time);
    let y = time.getFullYear();
    let m = time.getMonth() + 1;
    let d = time.getDate();
    let h = time.getHours();
    let mm = time.getMinutes();
    let s = time.getSeconds();
    return y + '-' + fillZero(m) + '-' + fillZero(d) + ' ' + fillZero(h) + ':' + fillZero(mm) + ':' + fillZero(s);
}

function fillZero(num, length) {
    length = length || 2
    return (Array(length).join("0") + num).slice(-length);
}

function cd100(col) {
    if (col.LAY_COL) {
        return d100(col[col.LAY_COL.field])
    }
    return 0
}

function d100(num) {
    return num / 100
}

function cj100(col) {
    if (col.LAY_COL) {
        return j100(col[col.LAY_COL.field])
    }
    return 0
}

function j100(num) {
    return `${num}%`
}

/**
 * 显示消息
 * @param layer
 */
function showMsg(layer) {
    let msg = getCookie('MSG');
    if (msg) {
        msg = decodeURIComponent(msg)
        layer.msg(msg);
        delCookie('MSG');
    }
}

/**
 * 获取cookie
 * @param name
 * @returns {string|null}
 */
function getCookie(name) {
    let arr, reg = new RegExp("(^| )" + name + "=([^;]*)(;|$)");
    if (arr = document.cookie.match(reg)) {
        return decodeURI(arr[2]);
    }
    return null;
}

/**
 * 删除cookie
 * @param name
 */
function delCookie(name) {
    let exp = new Date();
    exp.setTime(exp.getTime() - 1);
    document.cookie = name + "='';expires=" + exp.toGMTString() + ";Path=/";
}

/**
 * 设置cookie
 * @param name
 * @param value
 * @param exp 30 * 24 * 60 * 60
 */
function setCookie(name, value, exp) {
    let expo = new Date();
    expo.setTime(expo.getTime() + (exp || 30 * 24 * 60 * 60));
    document.cookie = name + "=" + encodeURI(value) + ";expires=" + expo.toGMTString() + ";Path=/";
}

function transUserId(d) {
    return `<span window="/users/detail?user=${d.UserId}" title="用户详情" style="color: #1E9FFF">${d.UserId}</span>`
}