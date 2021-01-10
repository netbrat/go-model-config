/**
 * 自定义表格
 */


layui.define(['form','table'],function(exports) {
    var table = layui.table,
        admin = layui.admin,
        $ = layui.jquery,
        $win = $(window),

        ctable = {
            toolbarEvents: {}, //工具栏事件集
            options: { //表格选择及数据
                elem: '#main_table',
                method: 'post',
                url: 'http://' + window.location.host + window.location.pathname,
                toolbar: '#main_table_toolbar',
                searchToolbar: '#main_table_search_toolbar',
                defaultToolbar: ['filter','print','exports','refresh','searchFlexible'],
                cellMinWidth: 60,
                page: true,
                limit: 50,
                cols: [],
            },
            defaultToolbar: {
                refresh: {
                    title: '刷新',
                    layEvent: 'refresh',
                    icon: 'layui-icon-refresh',
                },
                searchFlexible:{
                    title: '折叠查询栏',
                    layEvent: 'searchFlexible',
                    icon: 'layui-icon-up',
                }
            }
        };


    //---------方法定义------------------------------------------------------------

    /**
     * 全屏表格自动调整高度
     */
    ctable.resizeFullTable = function(){

        var winH = $(window).height(), //窗口高度
            fullTtables = $('.admin-table-full'); //窗体内所有全屏表格

        fullTtables.each(function(){
            othis = $(this);
            //获取各项高度
            searchH = othis.find('.admin-search-form').outerHeight(true) || 0,//查询框容器高度
            toolH = othis.find('.layui-table-tool').outerHeight(true) || 0, //表格工具栏高度
            headerH = othis.find('.layui-table-header').outerHeight(true) || 0 ,//表头高度
            pageH = othis.find('.layui-table-page').outerHeight(true) || 0, //分页栏高度
            newH = parseInt(winH-searchH-toolH-headerH-pageH-2);
            othis.find('.layui-table-body').height(newH); //设置表格主体高度
        });

    };

    ctable.getDefaultToolbar = function(options){
        options = options || [];

        for(key in ctable.defaultToolbar){
            index = options.indexOf(key);
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

        console.log(options)

        options.defaultToolbar = ctable.getDefaultToolbar(options.defaultToolbar);

        table.render(options); //渲染表格

        //插入查询栏
        if(typeof options.searchToolbar!='undefined' && options.searchToolbar){
            layid = options.elem.substring(1,options.elem.length);
            $(".layui-table-view[lay-id=" + layid + "]").prepend($(options.searchToolbar).html());
            $(options.searchToolbar).html('');
        }



        ctable.resizeFullTable();
    };


    //---------事件定义------------------------------------------------------------

    /**
     * 伸缩查询框事件
     */
    ctable.toolbarEvents.searchFlexible = function(obj){
        var CSS_ICON_UP = 'layui-icon-up',
            CSS_ICON_DOWN = 'layui-icon-down',
            iconElem = $(this).children('i'),
            layid = obj.config.id;

        if(iconElem.hasClass(CSS_ICON_UP)){ //展开-> 收缩
            iconElem.removeClass(CSS_ICON_UP).addClass(CSS_ICON_DOWN);
            $('[lay-id=' + layid + ']').find('.admin-search-form').addClass('layui-hide');
        }else{ //收缩->展开
            iconElem.removeClass(CSS_ICON_DOWN).addClass(CSS_ICON_UP);
            $('[lay-id=' + layid + ']').find('.admin-search-form').removeClass('layui-hide');
        }
        ctable.resizeFullTable();
    };



    //---------事件监听------------------------------------------------------------

    //监听表头工具栏
    table.on('toolbar', function(obj){
        //使用连接方式
        if (obj.event.substring(0,4) === "mtt_") {
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


    $(function(){
        // //渲染自定义数据表格

    });


    exports('ctable', ctable);
});