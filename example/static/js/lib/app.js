/**
主入口
 */


// layui.config({
//     common: '/static/js/lib/'
// }).use('app', function () {
//     if (typeof pageCallBack === "function") pageCallBack();//页面回调
// });


layui.config({
    base: '/static/js/lib/'
}).use([
    'jquery', 'element', 'layer', 'form',
    'ajaxform',
    'utils','admin'
],function(){

    //扩展 lib 目录下的其它模块
    var extend = ['customtable'];
    layui.each(extend, function(index, item){
        var mods = {};
        mods[item] = '{/}' + "/static/js/lib/extend/" + item;
        layui.extend(mods);
    });

    //将主目录定位在modules下
    layui.config({
       base: '/static/js/modules/'
    });

    if (typeof pageCallBack === "function") pageCallBack();//页面回调

    // exports('app',{});
});

