
/**
 * 常用函数方法
 */

layui.define(function(exports){
    var $ = layui.jquery,
        utils = {};

    /**
     * 判断变量是否为非定义或为空或为""
     * @param v
     * @returns {boolean}
     */
    utils.isEmptyOrNull = function(v){
        return  (v===undefined || v==="" || v===null);
    };

    /**
     * 判断内容是否是json格式字符串
     * @param data
     * @returns {boolean}
     */
    utils.isJson = function (data){
        try{
            $.parseJSON(data);
            return true;
        }catch(e){
            return false;
        }
    };

    /**
     * 文本编码
     * @param str
     * @returns {string}
     */
    utils.textEncode = function(str) {
        if(utils.isEmptyOrNull(str)) return "";
         str = str.replace(/&amp;/gi, '&');
         str = str.replace(/</g, '&lt;');
         str = str.replace(/>/g, '&gt;');
         return str;
    };

    /**
     * 文本解码
     * @param str
     * @returns {string}
     */
    utils.textDecode = function(str) {
        if(utils.isEmptyOrNull(str)) return "";
         str = str.replace(/&amp;/gi, '&');
         str = str.replace(/&lt;/gi, '<');
         str = str.replace(/&gt;/gi, '>');
         return str;
	};

    /**
     * 将多个URL参数串(格式x=0&y=1串)加到URL上
     * @param url
     * @param params
     */
    utils.setUrlParams = function(url, params){
        if(params instanceof Object){
            params = parseParam(params);
        }else if(utils.isJson(params)){
            params = parseParam(params);
        }

        if(url.indexOf('?')!==-1){
            str = url.substr(url.indexOf('?')+1)
            if(str==="" || str.substr(str.length+1)==='&'){
                return url + params;
            }else{
                return url + '&' + params;
            }
        }else{
            return url + '?' + params;
        }

    };

    /**
     * 将数组或json格式的转成URL参数串
     * @param param
     * @param key
     * @returns {string}
     */
    utils.parseParam = function(param, key){
        if(utils.isEmptyOrNull(param)) return "";
        var paramStr = "";
        if(param instanceof String||param instanceof Number||param instanceof Boolean){
            paramStr+="&"+key+"="+encodeURIComponent(param);
        }else{
            $.each(param,function(i){
                var k=utils.isEmptyOrNull(key) ? i : key+(param instanceof Array?"["+i+"]":"."+i);
                paramStr+='&'+parseParam(this, k);
            });
        }
        return paramStr.substr(1);
    };

    /**
     * ajax 提交
     * @param options
     *     type: 提交方式（get|post)
     *     data:提交的数据
     *     dataType：提交的格式
     *     showType: 对应layui.type
     *     width：窗口宽度
     *     height： 窗口高度
     *
     *
     */
    utils.ajax = function (options){
        options = options || {};
        options.type = options.type || "GET";
        options.contentType = options.contentType || 'application/x-www-form-urlencoded';
        options.data = options.data || {};
        options.data._t = Date.now();
        options.showType = options.showType || 0;
        if (options.showType === 1 || options.showType===2){
            options.width = (options.width || '600') + 'px';
            options.height = (options.height || '450') + 'px';
        }
        console.log(options);
        $.ajax({
            url: options.url,
            type: options.type,
            data: options.data,
            dateType: options.dateType,
            beforeSend: function(){
                layer.load(2); //layer.open({type: 3, content: '正在发送数据请求，请稍候...', icon: 16});
            },
            complete: function(){
                layer.closeAll('loading');
            },
            success: function(data){
                if (utils.isJson(data)) {
                    if (typeof options.successJson === "function") {
                        options.successJson(data);
                    } else {
                        data = $.parseJSON(data);
                        layer.open({title: '提示信息', content: data.msg, icon: data.code == 0 ? 6 : 5});
                    }
                }else{
                    if (typeof options.success === "function") {
                        options.success(data);
                    }else{
                        layer.open({type:options.showType, title:'提示信息', content: data, area:[options.width, options.height], maxmin:true, resize:true, moveOut:true});
                    }
                }
            },
            error: function(){
                if(typeof options.error === "function") {
                    options.error();
                }else{
                    layer.open({title: '错误信息', content: '发生未知的错误', icon: 5});
                }
            },
        });
    };


    exports('utils',utils);
});