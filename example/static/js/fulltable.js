/**
 * 自定义表格
 */

layui.define(['table','utils','admin'],function(exports) {
    var table = layui.table,
        admin = layui.admin,
        utils = layui.utils,
        $ = layui.jquery,
        $win = $(window),
        defaultOptions = { //表格选择及数据
            id: 'main_table',
            elem: '#main_table',
            method: 'post',
            url: window.location.href, //'http://' + window.location.host + window.location.pathname,
            toolbar: '#main_table_toolbar',
            searchToolbar: '#main_table_search_toolbar',
            defaultToolbar: ['filter', 'print', 'exports', 'refresh', 'searchFlexible'],
            cellMinWidth: 80,
            autoSort: false,
            page: true,
            limit: 50,
            cols: []
        },
        extraDefaultToolbar = { //扩展默认工具栏定义
            refresh: {
                title: '刷新',
                layEvent: 'refresh',
                icon: 'layui-icon-refresh'
            },
            searchFlexible: {
                title: '折叠查询栏',
                layEvent: 'searchFlexible',
                icon: 'layui-icon-up'
            }
        },
        fulltable = {
            toolbarEvents: {}, //工具栏事件集
            options:null,
            selectType: "",
            rowDoubleClick: function(obj){}
        };


    //---------方法定义------------------------------------------------------------

    /**
     * 将扩展的默认工具栏合并到表格工具栏中
     * @param options
     * @returns {*|Array}
     */
    fulltable.margeDefaultToolbar = function(options){
        options = options || [];
        for(var key in extraDefaultToolbar){
            var index = options.indexOf(key);
            if(index > 0){
                options[index] = extraDefaultToolbar[key];
            }
        }
        return options;
    };

    /**
     * 渲染表格
     */
    fulltable.render = function(options){
        if(typeof options=='undefined') {
            fulltable.options = defaultOptions;//如果没有传参，侧使用fulltable.options的值
        } else{
            fulltable.options = $.extend(true, defaultOptions, options)
        }
        fulltable.options.defaultToolbar = fulltable.margeDefaultToolbar(fulltable.options.defaultToolbar); //将扩展默认工具栏合并
        var selectType = options.cols[0][0].type || "";
        if (selectType==='checkbox' || selectType == 'radio'){
            fulltable.selectType = selectType;
        }

        table.render(fulltable.options); //渲染表格
        //初始化查询栏
        fulltable.initSearchToolbar($(fulltable.options.searchToolbar).html());

        fulltable.resize();
    };

    /**
     * 初始化工具栏方法
     */
    fulltable.initSearchToolbar = function(contentObj){
        if (!utils.isEmptyOrNull(fulltable)) {
            $('.layui-table-view[lay-id=' + fulltable.options.id + ']').prepend('<div class="main_table_search_toolbar"></div>');
            $('.main_table_search_toolbar').append(contentObj);
        }
    };

    /**
     * 表格重设大小
     */
    fulltable.resize = function(){
        var winH = $(window).height(), //窗口高度
            fullTtables = $('.admin-table-full'); //窗体内所有全屏表格

        fullTtables.each(function(){
            var othis = $(this),
                searchH = othis.find('.main_table_search_toolbar').outerHeight(true) || 0,//查询框容器高度
                toolH = othis.find('.layui-table-tool').outerHeight(true) || 0, //表格工具栏高度
                headerH = othis.find('.layui-table-header').outerHeight(true) || 0 ,//表头高度
                pageH = othis.find('.layui-table-page').outerHeight(true) || 0, //分页栏高度
                totalH = othis.find('.layui-table-total').outerHeight(true) || 0, //汇总栏高度
                newH = parseInt(winH - searchH - toolH - headerH - pageH -totalH - 2);
            othis.find('.layui-table-body').height(newH); //设置表格主体高度
        });
        //回调
        if (fulltable.onResize instanceof Function){
            fulltable.onResize();
        }
    };

    /**
     * 刷新方法
     */

    fulltable.refresh = function (options) {
        var searchObj = $('.main_table_search_toolbar').children();
        table.reload(fulltable.options.id, options);
        fulltable.initSearchToolbar(searchObj);
        fulltable.resize();
    };

    /**
     * 查询方法
     */
    fulltable.search = function(){
        searchForm = utils.jqueryId(admin.global.searchFormId);
        var where = utils.urlParamsToJSON($(searchForm).formSerialize());
        layui.fulltable.refresh({
            where: where,
            page:1
        });
        return false;
    };

    /**
     * 只选中当前行 （checkbox类型)
     * @param obj
     */
    fulltable.onlyOneSelected = function(obj){
        if (fulltable.selectType!=="checkbox") {return}
        var start = obj.tr.selector.indexOf('data-index=') + 12,
            index = obj.tr.selector.substring(start, obj.tr.selector.length - 2);

        //去除所有的行
        $('.layui-table-body tr').removeClass('layui-table-click').find('.laytable-cell-checkbox .layui-form-checkbox').removeClass('layui-form-checked');
        for (var i = 0; i < table.cache[fulltable.options.id].length; i++) {
            table.cache[fulltable.options.id][i].LAY_CHECKED = (i.toString() === index);
        }
        //单前行
        obj.tr.addClass('layui-table-click').find('.laytable-cell-checkbox .layui-form-checkbox').addClass('layui-form-checked');
    };

    /**
     * 行单击事件回调方法
     */

    fulltable.onRow = function(obj){};

    /**
     * 行双击事件回调方法
     */
    fulltable.onRowDouble = function(obj){};

    /**
     * 工具栏事件回调方法
     */
    fulltable.onToolbar = function(obj){};


    /**
     * 行工具事件回调方法
     */
    fulltable.onTool = function(obj){};

    /**
     * 排序事件回调方法
     */
    fulltable.onSort = function(obj){
        var order = '';
        if (obj.type){
            order = obj.field + ' ' + obj.type;
        }
        fulltable.refresh({
            initSort: obj,
            where: {'order': order}
        })
    };

    /**
     * 复选框事件回调方法
     */
    fulltable.onCheckbox = function(obj){};

    /**
     * 单选框事件回调方法
     */
    fulltable.onRadio = function(obj){};

    /**
     * 重设大小事件回调方法
     */
    fulltable.onResize = function(){};

    //---------事件定义------------------------------------------------------------

    /**
     * 伸缩查询框事件
     */
    fulltable.toolbarEvents.searchFlexible = function(obj){
        var CSS_ICON_UP = 'layui-icon-up',
            CSS_ICON_DOWN = 'layui-icon-down',
            iconElem = $(this).children('i');

        if(iconElem.hasClass(CSS_ICON_UP)){ //展开-> 收缩
            iconElem.removeClass(CSS_ICON_UP).addClass(CSS_ICON_DOWN);
            $('[lay-id=' + fulltable.options.id + ']').find('.main_table_search_toolbar').addClass('layui-hide');
        }else{ //收缩->展开
            iconElem.removeClass(CSS_ICON_DOWN).addClass(CSS_ICON_UP);
            $('[lay-id=' + fulltable.options.id + ']').find('.main_table_search_toolbar').removeClass('layui-hide');
        }
        fulltable.resize();
    };

    /**
     * 刷新
     * @param obj
     */
    fulltable.toolbarEvents.refresh = function(obj){
        fulltable.refresh();
    };

    //---------事件监听------------------------------------------------------------


    /**
     * 监听表头工具栏
     */
    table.on('toolbar', function(obj){
        //使用连接方式
        if (obj.event.substring(0,9) === 'openLink_') {
            admin.openLink('#'+ obj.event);
        }
        //使用事件方式
        fulltable.toolbarEvents[obj.event] && fulltable.toolbarEvents[obj.event].call(this, obj);
        //回调
        if (fulltable.onToolbar instanceof Function){
            return fulltable.onToolbar(obj);
        }
    });

    /**
     * 监听行工具事件
     */

    table.on('tool', function(obj){
        //回调
        if (fulltable.onTool instanceof Function){
            return fulltable.onTool(obj);
        }
    });

    /**
     * 监听行单击事件
     */
    table.on('row', function(obj){
        if (fulltable.selectType !== '') {
            var start = obj.tr.selector.indexOf('data-index=') + 12,
                index = obj.tr.selector.substring(start, obj.tr.selector.length - 2);

            if (fulltable.selectType === "radio") {
                //去除所有的行
                $('.layui-table-body tr').removeClass('layui-table-click').find('.laytable-cell-radio .layui-form-radio').removeClass('layui-form-radioed').find("i").removeClass("layui-anim-scaleSpring").html('&#xe63f;');
                for (var i = 0; i < table.cache[fulltable.options.id].length; i++) {
                    table.cache[fulltable.options.id][i].LAY_CHECKED = (i.toString() === index);
                }
                //当前行
                obj.tr.addClass('layui-table-click').find('.laytable-cell-radio .layui-form-radio').addClass('layui-form-radioed').find("i").addClass("layui-anim-scaleSpring").html('&#xe643;');

            } else {
                var checked = table.cache[fulltable.options.id][index].LAY_CHECKED || false; // 当前选中状态
                if (checked) {
                    obj.tr.removeClass('layui-table-click').find('.laytable-cell-checkbox .layui-form-checkbox').removeClass('layui-form-checked');
                } else {
                    obj.tr.addClass('layui-table-click').find('.laytable-cell-checkbox .layui-form-checkbox').addClass('layui-form-checked');
                }
                //更新缓存
                table.cache[fulltable.options.id][index].LAY_CHECKED = !checked;
            }
        }
        //回调
        if (fulltable.onRow instanceof Function){
            return fulltable.onRow(obj);
        }
    });


    /**
     * 监听行双击事件
     */
    table.on('rowDouble', function(obj){
        if (fulltable.onRowDouble instanceof Function) {
            return fulltable.onRowDouble(obj);
        }
    });


    /**
     * 监听复选框
     */
    table.on('checkbox', function (obj) {
        var trs = obj.type === 'all' ? $('.layui-table-body').find('tr') : obj.tr;
        if (obj.checked){
            trs.addClass('layui-table-click');
        }else{
            trs.removeClass('layui-table-click');
        }
        //回调
        if (fulltable.onCheckbox instanceof Function){
            return fulltable.onCheckbox(obj);
        }
    });

    /**
     * 监听单选框
     */
    table.on('radio', function (obj) {
        //回调
        if (fulltable.onCheckbox instanceof Function){
            return fulltable.onCheckbox(obj);
        }
    });

    /**
     * 监听排序
     */
    table.on('sort', function(obj){
        //回调
        if (fulltable.onSort instanceof Function){
            return fulltable.onSort(obj);
        }
    });

    /**
     * 监听窗口大小变化
     */
    $win.on('resize',function(){
        fulltable.resize();

    });


    exports('fulltable', fulltable);
});