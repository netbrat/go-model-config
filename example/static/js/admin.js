
let admin = {
    events: {}, //事件集
    tabsPage: { //主tabs
        index: 0, //当前选中项索引
    },
    side: { //侧边
        auto: true, //是否自动侧边伸缩
        spread: true,  //当前侧边伸缩状态（默认展开)
        isMinWin: false, //窗口大小状态
    }
};

const FILTER_TAB_TABS = 'admin-layout-tabs', //主标签页 LAY-FILTER
    ELEM_MAIN_TABS_LIS = '#admin_tabsheader>li'; //主标签项


/**
 * 侧边伸缩方法
 * @param spread  展开状态 true(展开)|false(收缩)
 * 实现方式：
 * 1、在导航菜单项上加 admin-logo-min-width（logo收缩宽度） 和 admin-logo-max-width（logo展开宽度） 属性
 * 2、在顶层页面和栏目主页面在外层加上id为 admin_app 的 div，在伸缩缩时，会自动在此div添加或移动 side-shrink 样式
 * 3、在样式表中添加 side-shrink（收缩后）的样式
 * 4、会同步向栏目主页调用 syncFlexible（spread)方法
 */
admin.sideFlexible = function(spread){
    if (spread!==admin.side.spread) { //状态变更时处理
        const ID_ADMIN_APP = "#admin_app",
            ID_FLEXIBLE = "#admin_flexible",
            CSS_ICON_SPREAD = 'layui-icon-spread-left', //展开图标样式
            CSS_ICON_SHRINK = 'layui-icon-shrink-right', //收缩图标样式
            CSS_SIDE_SHRINK = 'side-shrink'; //主窗口收缩样式

        console.log($(ID_FLEXIBLE));
        if (spread) {   //从收缩到展开
            $(ID_FLEXIBLE).removeClass(CSS_ICON_SPREAD).addClass(CSS_ICON_SHRINK); //伸缩按钮切换到收缩图标
            $(ID_ADMIN_APP).removeClass(CSS_SIDE_SHRINK); //主窗口容器移除收缩样式
        } else { //从展开到收缩
            $(ID_FLEXIBLE).removeClass(CSS_ICON_SHRINK).addClass(CSS_ICON_SPREAD);//伸缩按钮切换到展开图标
            $(ID_ADMIN_APP).addClass(CSS_SIDE_SHRINK); //主窗口容器添加收缩样式
        }
        admin.side.spread = !admin.side.spread; //记录伸缩状态
    }

};

/**
 *根据窗口大小自动侧边伸缩
 */
admin.autoSideFlexible = function(){
    if(window.top !== window.self) return; //如果不是顶层窗口则退出
    if (!admin.side.auto) return; //当前面是手动收缩侧边时，则不会自动展开
    let isMinWin = ($(window).width() <=992);
    if (isMinWin === admin.side.isMinWin) return; //当前状态未发生变化时退出
    admin.side.isMinWin = isMinWin;
    admin.sideFlexible(!isMinWin);
};


/**
 * 打开页面标签方法
 * @param title  标签文本
 * @param url  url
 * @param params url参数
 * @param noClose  是否是不允许关闭的标签 true|false
 */
admin.openTabsPage = function(title, url, params, noClose){
    url = utils.setUrlParams(url, params);
    let matchTo,
        tabs = $(ELEM_MAIN_TABS_LIS);

    tabs.each(function(index){
        let li = $(this),
            layId = li.attr('lay-id');

        if(layId === url){
            matchTo = true;
            admin.tabsPage.index = index;
        }
    });
    title = title || '新标签页';
    if(!matchTo) { //添加一个标签
        admin.tabsPage.index = tabs.length;
        layui.element.tabAdd(FILTER_TAB_TABS,{
            title: '<span>' + title + '</span>',
            id: url,
            content: '<iframe src="' + url + '" class="admin-iframe"></iframe>'
        });
        if(noClose) {
            $(ELEM_MAIN_TABS_LIS).eq(admin.tabsPage.index).addClass('admin-no-close');
        }
    }
    //定位到当前标签
    layui.element.tabChange(FILTER_TAB_TABS, url);
};


/**
 * 在新窗口打开页面方法
 * @param url
 * @param params
 * @param target 打开页面的目标
 */
admin.openWin = function(url, params, target){
    url = utils.setUrlParams(url, params);
    if(target==="self"){
        self.location.href = url;
    }else{
        window.open(url);
    }
};

/**
 * 无窗口POST提交方式
 * @param url
 * @param params
 */
admin.openPost = function(url, params){
    utils.ajax({url: url, type:'post', data: params});
};

/**
 * 普通弹窗方式
 * @param title
 * @param url
 * @param params
 * @param width
 * @param height
 */
admin.openDialog = function (title, url, params, width, height){
    utils.ajax({url: url, type:'get', data: params,showType:1, width: width, height:height});
};

/**
 * 编辑弹窗方式
 * @param title
 * @param url
 * @param params
 * @param editFormId
 * @param width
 * @param height
 */
admin.openEditDialog = function(title, url, params,  width, height, editFormId){
    editFormId = editFormId || "edit_form";
    let theForm = null;
    utils.ajax({
        url: url,
        type:'get',
        data: params,
        success: function(data){
            layer.open({
                type:1,
                title:title,
                content: data,
                area: [width, height],
                maxMin:true,
                resize:true,
                moveOut:true,
                btn:['确定','取消'],
                yes: function(){
                    theForm.submit();
                },
                success:function(){
                    theForm = $('#' + editFormId).ajaxForm({
                        url: url,
                        data: params,
                        beforeSubmit: function(){
                            layer.load(2);  //open({type:3, content:'数据提交中，请稍候...',icon:16});
                        },
                        complete: function(){
                            layer.closeAll('loading');
                        },
                        success:function(data){
                            if(utils.isJson(data)){
                                data = $.parseJSON(data);
                                if(data.code!==200){ //非正常
                                    layer.open({title:'提示信息1', content:data.msg, icon:6,
                                        yes: function(){
                                            layer.closeAll();
                                        }
                                    });
                                }else { //正常
                                    layer.open({title: '提示信息', content: data.msg, icon: 5});
                                }
                            }else{ // 不是json，直接显示data
                                layer.open({title:'提示信息', content:data})
                            }
                        },
                        error:function(){ // 未知错误
                            layer.open({title:'错误信息', content:'发生未知的错误', icon:5})
                        }
                    });
                }
            });
        }
    });
};


/**
 * 打开连接
 * 属性：
 * link-type, 连接类型，（1-js,0-连接(默认）
 * open-type,打开类型（0-标签页（默认）,1-新页面,2-普通弹窗,3-编辑弹窗,4-无窗口)
 * param-type, 参数获取类型（0-无参（默认），1-表单,2-单行列表,3-多行列表)
 * width, 弹窗宽度
 * height, 弹窗高度
 * title, 弹窗或tab标题
 * no-close, tab时是否不允许关闭 (只对标签页有效）
 * confirm, 操作前的确认提示信息
 * edit-form-id, 编辑时对应的编辑formid (只对编辑窗类型有效）
 * param-obj-id 参数获取对象id
 * logo-max-width logo最大宽度 (只对导航页有效）
 * logo-min-width logo最小宽度 (只对导航页有效）
 */
admin.openLink = function(obj){
    let oThis = $(obj),
        title = oThis.attr('title') || oThis.text(),
        url = oThis.attr('admin-href'),
        confirm = oThis.attr("confirm"),
        width = (oThis.attr('width') || '600') + 'px',
        height = (oThis.attr('height') || '450') + 'px',
        isCancel = false;
    if (utils.isEmptyOrNull(url)) {
        layer.open({title: '提示信息', content: "无法执行该操作!<br />该操作还未定义操作路径 [admin-href] 属性"});
    }
    if(!utils.isEmptyOrNull(confirm)){
        layer.confirm(confirm, {
            btn: ['确定', '取消']
        },function(){
            isCancel=false;
        },function() {
            isCancel = true;
        });
    }
    if (isCancel) return;

    if(oThis.attr('link-type')==="1"){ //js
        alert("OK");
        eval(url);
    }else{ //链接
        let params = '',
            paramType = oThis.attr('param-type'),
            openType = oThis.attr('open-type');

        if(!utils.isEmptyOrNull(paramType)){
            params = utils.getParam(paramType); //获取参数
            if(params===false) return;
        }
        //openType:打开方式(0-标签页, 1-新窗口, 2-本页面, 3-普通弹窗, 4-编辑弹窗, 5-无窗口)
        switch (parseInt(openType)) {
            case 1: //新窗口
                admin.openWin(url,params);
                break;
            case 2: //本页面
                admin.openWin(url,params,"self");
                break;
            case 3: //普通弹窗
                admin.openDialog(title, url, params, width, height);
                break;
            case 4: //编辑弹窗
                admin.openEditDialog(title, url, params, width, height, oThis.attr('edit-form-id'));
                break;
            case 5: //无窗口
                admin.openPost(url, params);
                break;
            default: //标签页
                admin.openTabsPage(title, url, params, oThis.attr('no-close'));
                break;
        }
    }
};

//---------事件定义--------------------------------------------------------------------------------------------------------

/**
 * 全屏事件
 * @param oThis
 */
admin.events.fullScreen = function(oThis){
    let CSS_SCREEN_FULL = 'layui-icon-screen-full',
        CSS_SCREEN_REST = 'layui-icon-screen-restore',
        iconElem = oThis.children("i");

    if(iconElem.hasClass(CSS_SCREEN_FULL)){
        let elem = document.body;
        if(elem.webkitRequestFullScreen){
            elem.webkitRequestFullScreen();
        } else if(elem.mozRequestFullScreen) {
            elem.mozRequestFullScreen();
        } else if(elem.requestFullScreen) {
            elem.requestFullscreen();
        }

        iconElem.addClass(CSS_SCREEN_REST).removeClass(CSS_SCREEN_FULL);
    } else {
        let elem = document;
        if(elem.webkitCancelFullScreen){
            elem.webkitCancelFullScreen();
        } else if(elem.mozCancelFullScreen) {
            elem.mozCancelFullScreen();
        } else if(elem.cancelFullScreen) {
            elem.cancelFullScreen();
        } else if(elem.exitFullscreen) {
            elem.exitFullscreen();
        }
        iconElem.addClass(CSS_SCREEN_FULL).removeClass(CSS_SCREEN_REST);
    }
};

/**
 * 侧边伸缩事件
 */
admin.events.flexible = function(oThis){
    admin.side.auto = !admin.side.spread; //当手动收缩时，则不进行自动伸缩
    admin.sideFlexible(!admin.side.spread);
};

/**
 * 左右滚动页面标签事件
 * @param type  滚动方式(auto|left|right)
 * @param index 当前标签页索引
 */
admin.events.rollPage = function(type, index){
    let tabsHeader = $('#admin_tabsheader'),
        liItem = tabsHeader.children('li'),
        scrollWidth = tabsHeader.prop('scrollwidth'),
        outerWidth = tabsHeader.outerWidth(),
        tabsLeft = parseFloat(tabsHeader.css('left'));
    if(type === 'left'){ //从左住右
        if(!tabsLeft && tabsLeft <= 0) return;
        //当前的left减去可视宽度，用于与上一轮的页标比较
        let prefLeft = -tabsLeft - outerWidth;
        liItem.each(function(index, item){
            var li = $(item),
                left = li.position().left;
            if(left >= prefLeft){
                tabsHeader.css('left', -left);
                return false;
            }
        });
    } else if(type=== 'auto'){ //自动滚动
        (function(){
            let thisLi = liItem.eq(index),
                thisLeft;

            if(!thisLi[0])  return;

            thisLeft = thisLi.position().left;

            //当目标标签在可视区域左侧时
            if(thisLeft < -tabsLeft){
                return tabsHeader.css('left', -thisLeft);
            }

            //当目标标签在可视区域右侧时
            if(thisLeft + thisLi.outerWidth() >= outerWidth - tabsLeft){
                //alert("OK");
                var subLeft = thisLeft + thisLi.outerWidth() - (outerWidth - tabsLeft);
                liItem.each(function(i, item){
                    let li = $(item),
                        left = li.position().left;

                    //从当前可视区域的最左第二个节点遍历，如果减去最左节点的差 > 目标在右侧不可见的宽度，则将该节点放置可视区域最左
                    if(left + tabsLeft > 0){
                        if(left - tabsLeft > subLeft){
                            tabsHeader.css('left', -left);
                            return false;
                        }
                    }
                });
            }
        }());
    } else {
        //默认向左滚动
        liItem.each(function(i, item){
            let li = $(item),
                left = li.position().left;

            if(left + li.outerWidth() >= outerWidth - tabsLeft){
                tabsHeader.css('left', -left);
                return false;
            }
        });
    }
};

/**
 * 向右滚动页面标签事件
 */

admin.events.leftPage = function(){
    admin.events.rollPage('left');
};

/**
 * 向左滚动页面标签事件
 */
admin.events.rightPage = function(){
    admin.events.rollPage();
};

/**
 * 关闭当前标签页事件
 */
admin.events.closeThisTabs = function(){
    if(!admin.tabsPage.index) return;
    let oThis = $(ELEM_MAIN_TABS_LIS).eq(admin.tabsPage.index);
    if(oThis.hasClass('admin-no-close')) return; //如果当前标签不允许关闭，则退出
    oThis.find(".layui-tab-close").trigger('click');
};

/**
 * 关闭其它标签页事件
 * @param type 类型： all(所有) | （其它）
 */
admin.events.closeOtherTabs = function(type){
    let currIndex = admin.tabsPage.index,
        TABS_REMOVE = 'admin-pagetabs-remove',
        tabsBodys = $('#admin_app_main_body>div.layui-tab-item'),
        tabs = $(ELEM_MAIN_TABS_LIS),
        noCloseNum = 0;

    tabs.each(function(index,item){
        if($(item).hasClass('admin-no-close')) { //记录不能关闭的标签数量
            if (index < currIndex) noCloseNum ++;
            if (index === currIndex && type === 'all') noCloseNum ++;
        }else{
            if(!(type!=='all' && index === currIndex)) { //给要关闭的标签全部加上移动临时样式
                $(item).addClass(TABS_REMOVE);
                tabsBodys.eq(index).addClass(TABS_REMOVE);
            }
        }
    });
    $('.' + TABS_REMOVE).remove(); //将所有要移除的标签全部移除
    $('#admin_app_tabsheader').css('left',0); //滚动到最左侧
    //选择一个未关闭的标签
    if(type === 'all'){
        admin.tabsPage.index = noCloseNum -1;
        $(ELEM_MAIN_TABS_LIS).eq(admin.tabsPage.index).trigger('click');
    }else{
        admin.tabsPage.index = noCloseNum;
    }
};

/**
 * 关闭全部标签页事件
 */
admin.events.closeAllTabs = function(){
    admin.events.closeOtherTabs('all');
};



//---------事件监听-------------------------------------------------------------

$(function(){
    /**
     * 监听点击事件
     * 范围（所有具有 admin-event 属性的对象
     */
    $('body').on('click', '*[admin-event]',function(){
        let oThis = $(this),
            attrEvent = oThis.attr('admin-event');
        console.log(attrEvent);
        admin.events[attrEvent] && admin.events[attrEvent].call(this, oThis);
    });


    /**
     * 监听页面连接点击事件
     * 范围 （所有具有 admin-href 属性的对象）
     * 其它属性：参考admin.openLink()说明
     */
    $('body').on('click', '*[admin-href]', function(){
        admin.openLink(this);
    });


    /**
     * 监听 主 tab 组件切换，同步 admin.tabsPage.index 及菜单选择
     */
    layui.element.on('tab('+ FILTER_TAB_TABS +')', function(data){
        admin.tabsPage.index = data.index;
        admin.events.rollPage('auto', data.index);
        let layId = $(ELEM_MAIN_TABS_LIS).eq(admin.tabsPage.index).attr('lay-id');
        $('[admin-href]').parent().removeClass('layui-this');
        $("[admin-href='" + layId + "']").parent().addClass('layui-this')
    });

    /**
     * 监听 tab 组件删除，同步 admin.tabsPage.index
     * layui tabs 在删除时有个bug，当删除的标签不是排在最后，则切换后的index会比实际的大1
     */
    layui.element.on('tabDelete('+ FILTER_TAB_TABS +')', function(data){
        if (data.index < admin.tabsPage.index){
            admin.tabsPage.index -= 1;
        }
    });

    /**
     * 监听窗口大小变化
     */
    $(window).on('resize',function(){
        admin.autoSideFlexible(); //自动侧边伸缩
    });

    /**
     * 监听tips
     */
    $('body').on('mouseenter','*[lay-tips]',function(){
        let oThis = $(this),
            tips = oThis.attr('lay-tips');
        layer.tips(tips,oThis);
    }).on('mouseleave', '*[lay-tips]', function(){
        layer.close(layer.index);
    });

    /**
     * 侧边分组tabs点击时
     */
    $("#admin_menu_group_tabs>.layui-tab-title>li").on('click',function () {
        //展开侧边
        admin.sideFlexible(true);
    }).on('mouseenter',function () {
        //显示tips
        if (!$("#admin_app").hasClass('side-shrink')) return;
        let text = $(this).find("span").html();
        layui.layer.tips(text, this);
    }).on('mouseleave',function () {
        //隐藏tips
        layui.layer.close(layer.index);
    });

    /**
     * 自动打开
     * 之所以要加延时处理，是因为layui动态加载模块时，并不随页面加载完成而完成
     */
    setTimeout(function(){
        $('body').find('[admin-default-open]').click(); //自动打开页面
    },200);

});

