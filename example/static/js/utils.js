
/**
 * 常用函数方法
 */

layui.define(['jquery'], function(exports){
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
     * 判断内容是否是json格式
     * @param data
     * @returns {boolean}
     */
    utils.isJSON = function (data){
        return typeof(data) == "object" && Object.prototype.toString.call(data).toLowerCase() == "[object object]" && !data.length;
    };

    utils.parseJSON = function (data) {
        if (!utils.isJSON(data)) {
            data = JSON.parse(str)
        }
        return data
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
            params = utils.parseParam(params);
        }else if(utils.isJSON(params)){
            params = utils.parseParam(params);
        }

        if(url.indexOf('?')!==-1){
            var str = url.substr(url.indexOf('?')+1);
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
                paramStr+='&'+utils.parseParam(this, k);
            });
        }
        return paramStr.substr(1);
    };

    /**
     * 将url参数转成json     *
     * @param params name1=value1&amp;name2=value2
     * @return { name1:value1, name2:value2 }
     */
    utils.urlParamsToJSON = function(params){
        var hash,
            myJson = {},
            hashes = params.split('&');

        for(var i=0; i<hashes.length; i++){
            hash = hashes[i].split('=');
            var key = hash[0];
            if (myJson.hasOwnProperty(key)){
                if (myJson[key] instanceof Array){
                    myJson[key].push(hash[1]);
                }else{
                    myJson[key] = [myJson[key], hash[1]];
                }
            }else {
                myJson[key] = hash[1];
            }
        }
        return myJson;
    };

    /**
     * 将id字符串转成jquery的id
     * id -> #id
     * @param id
     * @returns {string}
     */
    utils.jqueryId = function(id){
        if (typeof(id)== 'string' && id.substring(0,1) !== "#") {
            id = "#" + id
        }
        return id;
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
     */
    utils.ajax = function (options){
        options = options || {};
        options.type = options.type || "get";
        options.data = options.data || {};
        options.data._t = Date.now();
        options.showType = options.showType || 0;
        if (options.showType === 1 || options.showType===2){
            options.width = (options.width || '600') + 'px';
            options.height = (options.height || '450') + 'px';
        }
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
                try {
                    data = utils.parseJSON(data);
                    if (typeof options.successJson === "function") { //如果设置了成功json回调
                        options.successJson(data);
                    }else{ //未设置回调
                        if (data.code.toString() === "0") { //返回成功
                            layer.open({
                                title: '提示信息', content: data.msg, icon: 6,
                                yes: function () {
                                    layer.closeAll();
                                    try{layui.fulltable.refresh();}catch (e) {} //表格存在时，刷新表格
                                }
                            });
                        } else { //返回失败
                            layer.open({title: '出错啦', content: data.msg, icon: 5});
                        }
                    }
                }catch(e) { //非JSON格式
                    if (typeof options.success === "function") { //回调
                        options.success(data);
                    }else {
                        layer.open({
                            type: options.showType,
                            title: '提示信息',
                            content: data,
                            area: [options.width, options.height],
                            maxmin: true,
                            resize: true,
                            moveOut: true
                        });
                    }
                }
            },
            error: function(){
                if(typeof options.error === "function") { //回调
                    options.error();
                }else{
                    layer.open({title: '错误信息', content: '发生未知的错误', icon: 5});
                }
            },
        });
    };


    exports('utils',utils);
});