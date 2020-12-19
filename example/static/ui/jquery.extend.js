

/**
 * 判断变量是否为非定义或为空或为""
 * @param v
 * @returns {boolean}
 */
$.isEmptyOrNull = function(v){
    return  (v==undefined || v=="" || v==null);
};

/**
 * 判断内容是否是json格式字符串
 * @param data
 * @returns {boolean}
 */
$.isJson = function (data){
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
$.textEncode = function(str) {
    if($.isEmptyOrNull(str)) return "";
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
$.textDecode = function(str) {
    if($.isEmptyOrNull(str)) return "";
     str = str.replace(/&amp;/gi, '&');
     str = str.replace(/&lt;/gi, '<');
     str = str.replace(/&gt;/gi, '>');
     return str;
};



/**
 *
 * 将form表单元素的值序列化成对象
 *
 * @returns object
 */
$.serializeObject = function(form) {
	var o = {};
	$.each(form.serializeArray(), function(index) {
		if (o[this['name']]) {
			o[this['name']] = o[this['name']] + "," + this['value'];
		} else {
			o[this['name']] = this['value'];
		}
	});
	return o;
};


/**
 * 将多个URL参数串(格式x=0&y=1串)加到URL上
 * @param url
 * @param params
 */
$.setUrlParams = function(url, params){
    if(params instanceof Object){
        params = $.parseParam(params);
    }else if($.isJson(params)){
        params = $.parseParam(params);
    }

    if(url.indexOf('?')!==-1){
        str = url.substr(url.indexOf('?')+1)
        if(str=="" || str.substr(str.length+1)=='&'){
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
$.parseParam = function(param, key){
    if($.isEmptyOrNull(param)) return "";
    var paramStr = "";
    if(param instanceof String||param instanceof Number||param instanceof Boolean){
        paramStr+="&"+key+"="+encodeURIComponent(param);
    }else{
        $.each(param,function(i){
            var k=$.isEmptyOrNull(key) ? i : key+(param instanceof Array?"["+i+"]":"."+i);
            paramStr+='&'+parseParam(this, k);
        });
    }
    return paramStr.substr(1);
};


/**
 * 获得URL参数
 *
 * @returns 对应名称的值
 */
$.getUrlParam = function(name) {
	var reg = new RegExp("(^|&)" + name + "=([^&]*)(&|$)");
	var r = window.location.search.substr(1).match(reg);
	if (r != null)
		return unescape(r[2]);
	return null;
};

/**
 * 接收一个以逗号分割的字符串，返回List，list里每一项都是一个字符串
 *
 * @returns list
 */
$.getList = function(value) {
	if (value != undefined && value != '') {
		var values = [];
		var t = value.split(',');
		for ( var i = 0; i < t.length; i++) {
			values.push('' + t[i]);/* 避免他将ID当成数字 */
		}
		return values;
	} else {
		return [];
	}
};

/**
 * json字符串转换为Object对象.
 * @param json json字符串
 * @returns Object
 */
$.jsonToObj = function(json){
    return eval("("+json+")");
};

/**
 * json对象转换为String字符串对象.
 * @param o Json Object
 * @returns   Object对象
 */
$.jsonToString = function(o) {
	var r = [];
	if (typeof o == "string")
		return "\"" + o.replace(/([\'\"\\])/g, "\\$1").replace(/(\n)/g, "\\n").replace(/(\r)/g, "\\r").replace(/(\t)/g, "\\t") + "\"";
	if (typeof o == "object") {
		if (!o.sort) {
			for ( var i in o)
				r.push(i + ":" + obj2str(o[i]));
			if (!!document.all && !/^\n?function\s*toString\(\)\s*\{\n?\s*\[native code\]\n?\s*\}\n?\s*$/.test(o.toString)) {
				r.push("toString:" + o.toString.toString());
			}
			r = "{" + r.join() + "}";
		} else {
			for ( var i = 0; i < o.length; i++)
				r.push(obj2str(o[i]));
			r = "[" + r.join() + "]";
		}
		return r;
	}
	return o.toString();
};


/**
 * 根据长度截取先使用字符串，超长部分追加...
 * @param str 对象字符串
 * @param len 目标字节长度
 * @return 处理结果字符串
 */
$.cutString = function(str, len) {
	//length属性读出来的汉字长度为1
	if(str.length*2 <= len) {
		return str;
	}
	var strlen = 0;
	var s = "";
	for(var i = 0;i < str.length; i++) {
		s = s + str.charAt(i);
		if (str.charCodeAt(i) > 128) {
			strlen = strlen + 2;
			if(strlen >= len){
				return s.substring(0,s.length-1) + "...";
			}
		} else {
			strlen = strlen + 1;
			if(strlen >= len){
				return s.substring(0,s.length-2) + "...";
			}
		}
	}
	return s;
};

/**
 *
 * 增加formatString功能
 *
 * 使用方法：$.formatString('字符串{0}字符串{1}字符串','第一个变量','第二个变量');
 *
 * @returns 格式化后的字符串
 */
$.formatString = function(str) {
	for ( var i = 0; i < arguments.length - 1; i++) {
		str = str.replace("{" + i + "}", arguments[i + 1]);
	}
	return str;
};


/**
 * 日期格式化.
 * @param value 日期
 * @param format 格式化字符串 例如："yyyy-MM-dd"、"yyyy-MM-dd HH:mm:ss"
 * @returns 格式化后的字符串
 */
 $.formatDate = function(value,format) {
     if($.isEmptyOrNull(value)) return "";
	var dt;
	if (value instanceof Date) {
		dt = value;
	} else {
		dt = new Date(value);
		if (isNaN(dt)) {
			//将那个长字符串的日期值转换成正常的JS日期格式
			value = value.replace(/\/Date\((-?\d+)\)\//, '$1');
			dt = new Date();
			dt.setTime(value);
		}
	}
	return dt.format(format);
};

 /**
 * 扩展日期格式化 例：new Date().format("yyyy-MM-dd hh:mm:ss")
 *
 * "M+" :月份
 * "d+" : 天
 * "h+" : 小时
 * "m+" : 分钟
 * "s+" : 秒钟
 * "q+" : 季度
 * "S" : 毫秒数
 * "X": 星期 如星期一
 * "Z": 返回周 如周二
 * "F":英文星期全称，返回如 Saturday
 * "L": 三位英文星期，返回如 Sat
 * @param format 格式化字符串
 * @returns {*}
 */
Date.prototype.format = function(format) {
    var week = ['Sunday', 'Monday', 'Tuesday', 'Wednesday', 'Thursday', 'Friday', 'Saturday', '日', '一', '二', '三', '四', '五', '六'];
	var o = {
		"M+" : this.getMonth() + 1, //月份
		"d+" : this.getDate(), //天
		"h+" : this.getHours(), //小时
		"m+" : this.getMinutes(), //分钟
		"s+" : this.getSeconds(), //秒钟
		"q+" : Math.floor((this.getMonth() + 3) / 3), //季度
		"S" : this.getMilliseconds(),//毫秒数
        "X": "星期" + week[this.getDay() + 7], //星期
        "Z": "周" + week[this.getDay() + 7],  //返回如 周二
        "F": week[this.getDay()],  //英文星期全称，返回如 Saturday
        "L": week[this.getDay()].slice(0, 3)//三位英文星期，返回如 Sat
	};
	if (/(y+)/.test(format))
		format = format.replace(RegExp.$1, (this.getFullYear() + "")
				.substr(4 - RegExp.$1.length));
	for ( var k in o)
		if (new RegExp("(" + k + ")").test(format))
			format = format.replace(RegExp.$1, RegExp.$1.length == 1 ? o[k]
					: ("00" + o[k]).substr(("" + o[k]).length));
	return format;
};