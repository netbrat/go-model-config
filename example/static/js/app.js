/**
 主入口
 */


layui.config({
    base: '/static/js/'
}).use([
    'jquery', 'element', 'layer', 'laydate', 'form','table',
    'ajaxform','utils','admin','fulltable'
], function(){
    //扩展 lib 目录下的其它模块
    // var extend = ['fulltable'];
    // layui.each(extend, function(index, item){
    //     var mods = {};
    //     mods[item] = '{/}' + '/static/js/extend/' + item;
    //     layui.extend(mods);
    // });

    //将主目录定位在modules下
    layui.config({
        base: '/static/js/modules/'
    });

    if (typeof pageCallBack === "function") pageCallBack();//页面回调
});

