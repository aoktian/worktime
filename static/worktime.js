$(document).ready(function() {
    // 设置页面最小宽度为1650px
    document.body.style.minWidth = "1650px";
});

function showSuccessToast(msg) {
    $.toast({
        text: msg,
        icon : 'success',
        hideAfter: 5000,
        position: 'top-center'
    })
}

function showErrorToast(msg) {
    $.toast({
        text: msg,
        icon : 'error',
        hideAfter: 5000,
        position: 'top-center'
    })
}

function handleResponse(r) {
    if (r.redirect_url) {
        window.location.href = r.redirect_url;
        return false;
    }

    if (r.assertToast) {
        showErrorToast(r.assertToast)
        return false;
    }

    if (r.errors) {
        $.alert({
            title: '提示信息',
            content: r.errors.map(error => `<p>${error}</p>`).join(''),
            type: 'red',
        });
        return false;
    }

    if (r.assertAlert) {
        $.alert({
            title: '提示信息',
            content: r.assertAlert,
        });
        return false;
    }

    if (r.assertDialog) {
        $('#dialog').html(r.assertDialog)
        return false;
    }

    if (r.toast) {
        showSuccessToast(r.toast)
    }

    if (r.alert) {
        $.alert({
            title: '提示信息',
            content: r.alert,
        });
    }

    for (var key in r) {
        // 判断 key 是否是字符串
        if (typeof key === 'string') {
            // 判断第一个字符是否是 '#'
            if (key.charAt(0) === '#') {
                $(key).html( r[key] )
            }
        }
    }

    if (r.dialog) {
        $('#dialog').html(r.dialog).modal('show');
    }

    if (r.script) {
        eval(r.script)
    }

    return true
}

function getFormData(formSelector) {
    const $form = $(formSelector);

    // 获取表单数据并转换为JSON格式
    const formElements = $form.find('[name]');
    const jsonData = {};

    formElements.each(function() {
        const $el = $(this);
        const name = $el.attr('name');
        const type = $el.attr('type');

        let value = $el.val();
        // 类型转换：如果是数字字符串则转换为数字
        if (type != "password") {
            if (value !== "" && !isNaN(value) && !isNaN(parseFloat(value))) {
                value = Number(value);
            }
        }

        // 处理不同类型的表单元素
        if (type === 'checkbox') {
            if ($el.is(':checked')) {
                // 如果是复选框且被选中
                if (jsonData[name]) {
                    // 如果已经存在该字段，转换为数组
                    if (!Array.isArray(jsonData[name])) {
                        jsonData[name] = [jsonData[name]];
                    }
                    jsonData[name].push(value);
                } else {
                    jsonData[name] = value;
                }
            }
        } else if (type === 'radio') {
            // 如果是单选框且被选中
            if ($el.is(':checked')) {
                jsonData[name] = value;
            }
        } else {
            // 其他类型的表单元素
            jsonData[name] = value;
        }
    });

    return jsonData;
}

/**
 * 通用 AJAX 表单提交函数
 * 1. 基本用法 - 提交所有带有 .ajax-form 类的表单
 * submitFormAjax('.ajax-form');
 */
function submitFormAjax(formSelector, options) {
    // 默认配置
    const defaults = {
        method: 'POST',
        success: handleResponse,
        error: function(xhr, status, error) {
            console.error('表单提交失败:', error);
            $.alert({
                title: '表单提交失败',
                content: error,
            });
        },
        beforeSend: function() {
            // 可以在这里添加加载动画
        },
        complete: function() {
            // 请求完成后执行（无论成功或失败）
        }
    };

    // 合并配置选项
    const settings = $.extend({}, defaults, options);

    const $form = $(formSelector);

    const url = $form.attr('action') || window.location.href;
    const method = $form.attr('method') || settings.method;

    var jsonData = getFormData(formSelector);

    console.log('提交数据:', jsonData);

    $.ajax({
        url: url,
        type: method,
        contentType: 'application/json', // 设置为JSON格式
        data: JSON.stringify(jsonData), // 将数据转换为JSON字符串
        dataType: 'json',
        beforeSend: function() {
            // 禁用提交按钮防止重复提交
            $form.find('button[type="submit"]').prop('disabled', true);
            if (settings.beforeSend) settings.beforeSend.call(this);
        },
        success: function(response) {
            settings.success.call(this, response);
        },
        error: function(xhr, status, error) {
            settings.error.call(this, xhr, status, error);
        },
        complete: function() {
            // 重新启用提交按钮
            $form.find('button[type="submit"]').prop('disabled', false);
            if (settings.complete) settings.complete.call(this);
        }
    });
}

function ajax(url, jsonData, callback) {
    $.ajax({
        url: url,
        type: "POST",
        contentType: 'application/json', // 设置为JSON格式
        data: JSON.stringify(jsonData), // 将数据转换为JSON字符串
        dataType: "json",
        cache: false,
        success: function( r ) {
            if (!handleResponse(r)) {
                return;
            }

            if (callback) {
                callback( r );
            }
        },
        error: function (XMLHttpRequest, textStatus, errorThrown) {
            $.alert({
                title: '提示信息',
                content: XMLHttpRequest.responseText,
            });
        },
        complete: function (jqXHR, textStatus) {
        }
    });
}

function get_form_values(id, forme ) {
    var rtn = {};

    if (!forme) {
        forme = 'val';
    }

    var els = $("#" + id + " [itag='" + forme + "']");

    for (var i = 0; i < els.length; i++) {
        var el = $(els[i]);
        var value = el.val();
        // 类型转换：如果是数字字符串则转换为数字
        if (value !== "" && !isNaN(value) && !isNaN(parseFloat(value))) {
            value = Number(value);
        }
        if ( el.prop('type') != "checkbox" || el.prop("checked")) {
            rtn[el.prop("name")] = value
        }
    }
    return rtn;
}

function checkall( id, name, b ) {
    var els = $("#" + id + " :checkbox");
    for (var i = 0; i < els.length; i++) {
        var el = $(els[i]);
        if (name == el.prop("name")) {
            if ("undefined" == typeof(b) ) {
                if (el.prop("checked")) {
                    el.prop("checked", false);
                } else {
                    el.prop("checked", true);
                }
            } else{
                el.prop("checked", b);
            }
        }
    }
}


function getUsers( id ) {
    var a = [];
    for (var i in users) {
        if (users[i].department == id) {
            a.push(users[i]);
        }
    }

    a.sort(function(x, y){return x.team - y.team});

    var options = "";
    for (var i = 0; i < a.length; i++) {
        options += "<option value='" + a[i].id + "'>" + a[i].name + "</option>";
    }
    return options;
}

function onChangeDepartment( id, tar, first ) {
    var options = "";
    if (first) {
        options += "<option value='0'>"+first+"</option>";
    }

    $(tar).html( options + getUsers( id ) );
}

function getTags( id ) {
    var options = "";
    for (var i in tags) {
        if (tags[i].project_id == id) {
        options += "<option value='" + tags[i].id + "'>" + tags[i].name + "</option>";
        }
    }
    return options;
}
function onChangePro( id, tar, first ) {
    var options = "";
    if (first) {
        options += "<option value='0'>"+first+"</option>";
    }

    $(tar).html( options + getTags( id ) );
}

function batchModify( cb ) {
    // 收集所有带 itag="val" 属性的元素的值
    let updateData = get_form_values("batchModify");
    let ids = [];

    // 获取选中的任务ID（假设有一个复选框机制）
    $('#tasklist input[name="task-checkbox"]:checked').each(function() {
        ids.push(Number($(this).val()));
    });

    if (ids.length === 0) {
        alert('请选择至少一个任务');
        return;
    }

    //将updateData.level转为数字
    updateData.level = Number(updateData.level);

    // 构造请求数据
    let requestData = {
        ids: ids,
        update: updateData
    };

    // 以 JSON 格式发送 POST 请求
    $.ajax({
        url: '/task/batchupdate',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify(requestData),
        dataType: 'json',
        success: function(response) {
            if (handleResponse(response)) {
                cb();
            }
        },
        error: function(xhr, status, error) {
            console.error('批量修改失败:', error);
            alert('操作失败: ' + error);
        }
    });
}

var _page = 1;
function filter_result( page ) {
    if (page) {
      _page = page;
    }
    let requestData = get_form_values( "taskfilter" );
    requestData.page = _page;

    _getlist( requestData );
}

function _getlist( requestData ) {
    $.ajax({
        url: '/task/list',
        type: 'POST',
        contentType: 'application/json',
        data: JSON.stringify(requestData),
        dataType: 'json',
        success: handleResponse,
        error: function(xhr, status, error) {
            alert('操作失败: ' + error);
        }
    });
}


function initEditor( id, focus, height ) {
    if (typeof(height) == "undefined") {
        height = 300
    }
    if (typeof(focus) == "undefined") {
        focus = false
    }
    var $summernote = $(id);
    $summernote.summernote({
        focus : focus,
        toolbar : [
            ['style', ['style', 'bold', 'italic', 'underline', 'clear']],
            ['font', ['color', 'fontsize']],
            ['para', ['ul', 'ol', 'paragraph']],
            ['insert', ['hr', 'link', 'table', 'picture', 'video']],
            ['misc', ['codeview', 'fullscreen', 'help']]
        ],
        height: height,             // set minimum height of editor
        lang : "zh-CN",
        callbacks: {
            // onPaste: function(e) {
            //     var content = e.originalEvent.clipboardData.getData('text/plain');
            //     e.preventDefault();
            //     document.execCommand("insertHTML", false, "<pre>" + content + "</pre>");
            // },
            onImageUpload: function(files) {
                var data = new FormData();
                data.append("upload", files[0]);
                $.ajax({
                    data: data,
                    type: "POST",
                    url: "/task/uploadImage",
                    cache: false,
                    dataType: "json",
                    contentType: false,
                    processData: false,
                    success: function( res ) {
                        if (!res.path) {
                            $.alert({
                                title: '上传失败图片失败',
                                content: res.err,
                                type: 'red',
                            });
                            return;
                        }

                        console.log( res.path );
                        // var img = document.createElement("img");
                        // img.src = url;
                        // $summernote.summernote('insertNode', img);
                        $summernote.summernote('insertImage', res.path);
                    }
                });
            }
        }
    });
}