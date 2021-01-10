/**
 主入口
 */


layui.config({
    base: '/static/js/'
}).use([
    'jquery', 'element', 'layer', 'form',
    'ajaxform',
    'utils','admin'
], function(){
    //扩展 lib 目录下的其它模块
    var extend = ['ctable'];
    layui.each(extend, function(index, item){
        var mods = {};
        mods[item] = '{/}' + '/static/js/extend/' + item;
        layui.extend(mods);
    });

    //将主目录定位在modules下
    layui.config({
        base: '/static/js/modules/'
    });

    if (typeof pageCallBack === "function") pageCallBack();//页面回调
});

