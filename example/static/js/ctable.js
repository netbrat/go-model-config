/**
 * 自定义表格
 */

layui.define(['form','table'],function(exports) {
    var table = layui.table,
        admin = layui.admin,
        utils = layui.utils,
        $ = layui.jquery,
        $win = $(window),

        ctable = {
            toolbarEvents: {}, //工具栏事件集
            options: { //表格选择及数据
                id:"main_table",
                elem: '#main_table',
                method: 'post',
                url: window.location.href, //'http://' + window.location.host + window.location.pathname,
                toolbar: '#main_table_toolbar',
                searchToolbar: '#main_table_search_toolbar',
                defaultToolbar: ['filter', 'print', 'exports', 'refresh', 'searchFlexible'],
                cellMinWidth: 60,
                //totalRow: true,
                //totalRowText: "合计",
                page: true,
                limit: 50,
                cols: []
            },
            defaultToolbar: {
                refresh: {
                    title: '刷新',
                    layEvent: 'refresh',
                    icon: 'layui-icon-refresh'
                },
                searchFlexible:{
                    title: '折叠查询栏',
                    layEvent: 'searchFlexible',
                    icon: 'layui-icon-up'
                }
            },
            id: "",
            searchToolbar: ''
        };


    //---------方法定义------------------------------------------------------------

    /**
     * 全屏表格自动调整高度
     */
    ctable.resizeFullTable = function(){

        var winH = $(window).height(), //窗口高度
            fullTtables = $('.admin-table-full'); //窗体内所有全屏表格

        fullTtables.each(function(){
           var othis = $(this),
            searchH = othis.find('.admin-search-form').outerHeight(true) || 0,//查询框容器高度
            toolH = othis.find('.layui-table-tool').outerHeight(true) || 0, //表格工具栏高度
            headerH = othis.find('.layui-table-header').outerHeight(true) || 0 ,//表头高度
            pageH = othis.find('.layui-table-page').outerHeight(true) || 0, //分页栏高度
            totalH = othis.find(".layui-table-total").outerHeight(true) || 0, //汇总栏高度
            newH = parseInt(winH - searchH - toolH - headerH - pageH -totalH - 2);
            othis.find('.layui-table-body').height(newH); //设置表格主体高度
        });

    };

    /**
     * 将自定义的工具栏合并到表格工具栏中
     * @param options
     * @returns {*|Array}
     */
    ctable.getDefaultToolbar = function(options){
        options = options || [];
        for(var key in ctable.defaultToolbar){
            var index = options.indexOf(key);
            if(index > 0){
                options[index] = ctable.defaultToolbar[key];
            }
        }
        return options;
    };


    /**
     * 渲染表格
     */
    ctable.render = function(options){
        if(typeof options=='undefined') {
            options = ctable.options;//如果没有传参，侧使用ctable.options的值
        } else{
            options = $.extend(true, ctable.options, options)
        }

        options.defaultToolbar = ctable.getDefaultToolbar(options.defaultToolbar); //合并工具栏
        ctable.id = options.id;  //表格id
        table.render(options); //渲染表格
        //初始化查询栏
        ctable.searchToolbar = options.searchToolbar;
        ctable.initSearchToolbar();

        ctable.resizeFullTable();
    };

    //初始化工具栏方法
    ctable.initSearchToolbar = function(){
        if (!utils.isEmptyOrNull(ctable)) {
            $(".layui-table-view[lay-id=" + ctable.id + "]").prepend($(ctable.searchToolbar).html())
        }
    };

    //刷新方法(options可不传)
    ctable.refresh = function (options) {
        table.reload(ctable.id, options);
        ctable.initSearchToolbar();
        ctable.resizeFullTable();
    };

    //---------事件定义------------------------------------------------------------

    /**
     * 伸缩查询框事件
     */
    ctable.toolbarEvents.searchFlexible = function(obj){
        var CSS_ICON_UP = 'layui-icon-up',
            CSS_ICON_DOWN = 'layui-icon-down',
            iconElem = $(this).children('i');

        if(iconElem.hasClass(CSS_ICON_UP)){ //展开-> 收缩
            iconElem.removeClass(CSS_ICON_UP).addClass(CSS_ICON_DOWN);
            $('[lay-id=' + ctable.id + ']').find('.admin-search-form').addClass('layui-hide');
        }else{ //收缩->展开
            iconElem.removeClass(CSS_ICON_DOWN).addClass(CSS_ICON_UP);
            $('[lay-id=' + ctable.id + ']').find('.admin-search-form').removeClass('layui-hide');
        }
        ctable.resizeFullTable();
    };

    /**
     * 刷新
     * @param obj
     */
    ctable.toolbarEvents.refresh = function(obj){
        ctable.refresh();
    };

    //---------事件监听------------------------------------------------------------

    //监听表头工具栏
    table.on('toolbar', function(obj){
        //使用连接方式
        if (obj.event.substring(0,9) === "openLink_") {
            admin.openLink('#'+ obj.event);
        }
        //使用事件方式
        ctable.toolbarEvents[obj.event] && ctable.toolbarEvents[obj.event].call(this, obj);
    });

    /**
     * 监听窗口大小变化
     */
    $win.on('resize',function(){
        ctable.resizeFullTable(); //自动调整全屏表格高度
    });



    exports('ctable', ctable);
});