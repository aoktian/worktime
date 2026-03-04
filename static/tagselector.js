// 当项目改变时，清空版本并重置
function onChangePro(o, tagSelect) {
    var $tagSelect = $(tagSelect);
    $tagSelect.attr("project_id", $(o).val())

    // 清空现有选项，只保留默认的“版本”占位符
    $tagSelect.empty().append('<option value="0" selected>全部</option>');
    $tagSelect.val(0).trigger('change');
}


function initTagSelect2(url, selector, refreshTaskList) {
    var $tagSelect = $(selector);
    $tagSelect.select2({
        theme: "bootstrap-5",
        width: '120px',
        ajax: {
            url: url,
            dataType: 'json',
            contentType: 'application/json',
            type: 'POST',
            delay: 250, // 防抖处理，减少服务器压力
            data: function (params) {
                // 获取当前选中的项目ID
                var projectId = $tagSelect.attr("project_id")
                var searchData = {
                    keyword: params.term, // 用户输入的搜索词
                    page: params.page || 1,
                    project_id :parseInt(projectId)
                };

                return JSON.stringify(searchData);
            },
            processData: false, // 阻止 jQuery 处理数据
            processResults:function (data) {
                $.map(data.results, function (obj) {
                    obj.text = obj.text || obj.name; // replace name with the property used for the text
                    return obj;
                });
                data.results.unshift({id:0, text: "全部"})
                return data;
            }
        },
        //自定义显示：如果是已归档，加个灰色标签提示
        templateResult: function (repo) {
            if (repo.loading) return repo.text;
            var $container = $("<span>" + repo.text + "</span>");
            if (repo.is_archived) {
                $container.append(" <small class='text-muted'>(已归档)</small>");
            }
            return $container;
        },
        templateSelection: function (repo) {
            return repo.text;
        }
    })
    // 关键：绑定 Select2 专用选择事件
    .on('select2:select', function (e) {
        if (refreshTaskList) {
            _taskList(1);
        }
    })
    // 可选：绑定清除事件（如果点了那个小叉叉）
    .on('select2:unselect', function (e) {
        $(this).val(0).trigger('change');
        if (refreshTaskList) {
            _taskList(1);
        }
    });
}

function initProjectSelect2(selector, tagSelect, refreshTaskList) {
    var $select = $(selector);
    $select.select2({
        theme: "bootstrap-5",
        width: '120px',
        ajax: {
            url: "/project/selector",
            dataType: 'json',
            contentType: 'application/json',
            type: 'POST',
            delay: 250, // 防抖处理，减少服务器压力
            data: function (params) {
                // 获取当前选中的项目ID
                var searchData = {
                    keyword: params.term, // 用户输入的搜索词
                    page: params.page || 1
                };

                return JSON.stringify(searchData);
            },
            processData: false, // 阻止 jQuery 处理数据
            processResults:function (data) {
                $.map(data.results, function (obj) {
                    obj.text = obj.text || obj.name; // replace name with the property used for the text
                    return obj;
                });
                data.results.unshift({id:0, text: "全部"})
                return data;
            }
        },
        //自定义显示：如果是已归档，加个灰色标签提示
        templateResult: function (repo) {
            if (repo.loading) return repo.text;
            var $container = $("<span>" + repo.text + "</span>");
            if (repo.is_archived) {
                $container.append(" <small class='text-muted'>(已归档)</small>");
            }
            return $container;
        },
        templateSelection: function (repo) {
            return repo.text;
        }
    })
    // 关键：绑定 Select2 专用选择事件
    .on('select2:select', function (e) {
        if (tagSelect) {
            var $tagSelect = $(tagSelect);
            $tagSelect.attr("project_id", $(this).val())

            // 清空现有选项，只保留默认的“版本”占位符
            $tagSelect.empty().append('<option value="0" selected>全部</option>');
            $tagSelect.val(0).trigger('change');
        }
        if (refreshTaskList) {
            refreshTaskList(1)
        }
    })
    // 可选：绑定清除事件（如果点了那个小叉叉）
    .on('select2:unselect', function (e) {
        $(this).val(0).trigger('change');

        if (tagSelect) {
            var $tagSelect = $(tagSelect);
            $tagSelect.attr("project_id", 0)

            // 清空现有选项，只保留默认的“版本”占位符
            $tagSelect.empty().append('<option value="0" selected>全部</option>');
            $tagSelect.val(0).trigger('change');
        }

        if (refreshTaskList) {
            refreshTaskList(1);
        }
    });
}
