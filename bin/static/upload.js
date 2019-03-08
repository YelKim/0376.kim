//上传对象
var uploader = WebUploader.create({
    accept: {
        extensions: 'gif,jpg,jpeg,bmp,png,mp3,wav,xls,xlsx,wbmp'
    },
    fileVal: 'upfile',
    swf: '/lib/webuploader/0.1.5/Uploader.swf',
    // 其他配置项
    compress: {
        quality: 90,
        allowMagnify: false,
        preserveHeaders: true,
        noCompressIfLarger: false,
        compressSize: 1024 * 1024
    },
    auto: true,
    server: '/admin/simple/upload.html',
    compress: {
        // 图片质量，只有type为`image/jpeg`的时候才有效。
        quality: 90,
        // 是否允许放大，如果想要生成小图的时候不失真，此选项应该设置为false.
        allowMagnify: false,
        // 是否允许裁剪。
        crop: false,
        // 是否保留头部meta信息。
        preserveHeaders: true,
        // 如果发现压缩后文件大小比原来还大，则使用原来图片
        // 此属性可能会影响图片自动纠正功能
        noCompressIfLarger: false,
        // 单位字节，如果图片大小小于此值，不会采用压缩。
        compressSize: 100
    }
});
// 添加鼠标事件
function filelistMouseEvent ($li, $btns) {
    // 鼠标移动到图片上面的时候
    $li.on('mouseenter', function () {
        $btns.stop().animate({height: 30});
    });
    //鼠标移出图片上面时候
    $li.on('mouseleave', function () {
        $btns.stop().animate({height: 0});
    });

    //弹层上面按钮操作删除、设置为主图
    $btns.on('click', 'span', function () {
        var index = $(this).index(), deg;
        var fileid = $(this).parent().attr("fileid");
        var path = $(this).parent().attr("path");
        switch (index) {
            case 0: //删除
                var imgLength = $queue.find("li").length
                // 如果删除的是主图
                if ($("#mainurl").val() == path) {
                    if (imgLength > 1) {
                        // 循环查找不是主图的第一张作为主图
                        var i = 0;
                        do {
                            if ($('#WU_FILE_' + i).length) {
                                var imgPanel = $('#WU_FILE_' + i).find(".file-panel");
                                if (imgPanel.attr("fileid") != fileid) {
                                    $("#mainurl").val(imgPanel.attr("path"));
                                    break;
                                }
                            }
                            i++;
                        } while (true);
                    } else {
                        //全部删除
                        $("#mainurl").val("");
                    }
                    layer.msg("删除成功");
                }
                $('#' + fileid).remove();
                if (imgLength == 1) {
                    $placeHolder.removeClass('element-invisible');
                    $statusBar.hide();
                }
                var imgurl = $("#imgurl").val();
                $("#imgurl").val(imgurl.replace("," + path, ""));
                // 显示上传按钮
                $statusBar.show();
                break;
            case 1:
                //设置主图
                $("#mainurl").val(response.data.path);
                layer.msg("设置主图成功");
                break;
        }
        if (upload.supportTransition) {
            deg = 'rotate(' + file.rotation + 'deg)';
            $wrap.css({
                '-webkit-transform': deg,
                '-mos-transform': deg,
                '-o-transform': deg,
                'transform': deg
            });
        } else {
            $wrap.css('filter', 'progid:DXImageTransform.Microsoft.BasicImage(rotation=' + (~~((file.rotation / 90) % 4 + 4) % 4) + ')');
        }
        uploader.refresh();
    });
}
var $wrap = $('.uploader-list-container'),
    // 状态栏，包括进度和控制按钮
    $statusBar = $wrap.find('.statusBar'),
    // 优化retina, 在retina下这个值是2
    $ratio = window.devicePixelRatio || 1,
    $placeHolder = $wrap.find( '.placeholder' ),
    // 文件总体选择信息。
    $info = $statusBar.find('.info');

// 图片容器
$queue = $(".filelist");
if ($(".filelist").length == 0) {
    $queue = $('<ul class="filelist"></ul>').appendTo($wrap.find('.queueList'));
}

var upload = {
    multiple: false,
    loading: new Object(),
    // 添加的文件数量
    thumbWidth: 110 * $ratio,
    thumHeight: 110 * $ratio,
    isSupportBase64: (function () {
        var data = new Image();
        var support = true;
        data.onload = data.onerror = function () {
            if (this.width != 1 || this.height != 1) {
                support = false;
            }
        }
        data.src = "data:image/gif;base64,R0lGODlhAQABAIAAAAAAAP///ywAAAAAAQABAAACAUwAOw==";
        return support;
    })(),
    supportTransition: (function () {
        var s = document.createElement('p').style,
            r = 'transition' in s ||
                'WebkitTransition' in s ||
                'MozTransition' in s ||
                'msTransition' in s ||
                'OTransition' in s;
        s = null;
        return r;
    })(),
    flashVersion: (function () {
        var version;
        try {
            version = navigator.plugins['Shockwave Flash'];
            version = version.description;
        } catch (ex) {
            try {
                version = new ActiveXObject('ShockwaveFlash.ShockwaveFlash')
                    .GetVariable('$version');
            } catch (ex2) {
                version = '0.0';
            }
        }
        version = version.match(/\d+/g);
        return parseFloat(version[0] + '.' + version[1], 10);
    })(),
    // 多个文件文件添加进来时执行，负责view的创建
    multipleFile: function (file, response) {
        var $li = $('<li id="' + file.id + '">'
            + '<p class="title">' + file.name + '</p>'
            + '<p class="imgWrap"></p>'
            + '</li>'
            ),
            // 定义图片按钮删除、设置为主图
            $btns = $('<div class="file-panel" fileid="'+ file.id +'" path="'+ response.data.path +'">'
                + '<span class="cancel">删除</span>'
                + '<span class="cancel">设置主图</span>'
                + '</div>'
            ).appendTo($li);
        $wrap = $li.find('p.imgWrap');
        $info = $('<p class="error"></p>');
        if (file.getStatus() === 'invalid') {
            $info.text(file.statusText).appendTo($li);
        }
        // 预览图片
        uploader.makeThumb(file, function (error, src) {
            var img;
            if (error) {
                $wrap.text('不能预览');
                return;
            }
            img = $('<img src="' + src + '">');
            $wrap.empty().append(img);
        }, upload.thumbWidth, upload.thumHeight);
        if ($queue.find("li").length == 0) {
            $placeHolder.addClass('element-invisible');
            $statusBar.show();
        }
        if ($queue.find("li").length == 4) {
            $statusBar.hide();
        }
        // 添加鼠标事件
        filelistMouseEvent($li, $btns);

        //当上传第一张图片时候
        if ($queue.find("li").length == 0) {
            $("#mainurl").val(response.data.path);
        }
        var path = $("#imgurl").val();
        $("#imgurl").val(path + "," + response.data.path)
        uploader.refresh();
        $li.appendTo($queue);
    },
    // 单文件文件添加进来时执行，负责view的创建
    simpleFile: function (file, response) {
        var ext = "jpg|jpeg|png|bmp|wbmp|gif";
        var show_span = $("#rt_" + file.source.ruid).parent().next();
        //显示图片
        if (ext.indexOf(file.ext) != -1) {
            show_span.html('<img style="margin: 0 auto;" width="120" src="' + response.data.path + '" />');
        } else {
            show_span.html('<a htef="' + response.data.path + '" download="' + response.data.path + '">' + response.data.path + '</a>');
        }
        $(show_span).show();
        //保存文件地址
        if (show_span.next()) {
            show_span.next().val(response.data.path);
        }
    }
}

$(function () {
// 配置对象


    // 判断浏览器是否下载flash
    if (!WebUploader.Uploader.support('flash') && WebUploader.browser.ie) {
        // flash 安装了但是版本过低。
        if (upload.flashVersion) {
            (function (container) {
                window['expressinstallcallback'] = function (state) {
                    switch (state) {
                        case 'Download.Cancelled':
                            layer.msg('您取消了更新！')
                            break;
                        case 'Download.Failed':
                            layer.msg('安装失败')
                            break;
                        default:
                            layer.msg('安装已成功，请刷新！');
                            break;
                    }
                    delete window['expressinstallcallback'];
                };
                var swf = 'expressInstall.swf';
                var html = '<object type="application/'
                    + 'x-shockwave-flash" data="'
                    + swf + '" ';
                if (WebUploader.browser.ie) {
                    html += 'classid="clsid:d27cdb6e-ae6d-11cf-96b8-444553540000" ';
                }
                html += 'width="100%" height="100%" style="outline:0">' +
                    '<param name="movie" value="' + swf + '" />' +
                    '<param name="wmode" value="transparent" />' +
                    '<param name="allowscriptaccess" value="always" />' +
                    '</object>';

                container.html(html);
            })($wrap);
        } else {
            //没有装flash
            $wrap.html('<a href="http://www.adobe.com/go/getflashplayer" target="_blank" border="0"><img alt="get flash player" src="http://www.adobe.com/macromedia/style_guide/images/160x41_Get_Flash_Player.jpg" /></a>');
            return;
        }
    }
    if (!WebUploader.Uploader.support()) {
        layer.msg('Web Uploader 不支持您的浏览器！');
        return;
    }

    // 开始上传
    uploader.on('uploadStart', function (file) {
        upload.loading = layer.load(0, {shade: false});
    });

    // 拖拽时不接受 js, txt 文件。
    uploader.on('dndAccept', function (items) {
        var denied = false,
            len = items.length,
            i = 0,
            // 修改js类型
            unAllowed = 'text/plain;application/javascript ';
        for (; i < len; i++) {
            // 如果在列表里面
            if (~unAllowed.indexOf(items[i].type)) {
                denied = true;
                break;
            }
        }
        return !denied;
    });

    // 上传出错时候
    uploader.on('uploadError', function (file) {
        layer.msg(file);
    });

    // 上传成功
    uploader.on('uploadSuccess', function (file, response) {
        layer.close(upload.loading);
        layer.msg(response.message);
        if (response.code == 0) {
            if (upload.multiple) {
                // 多文件上传
                upload.multipleFile(file, response);
            } else {
                // 单文件上传
                upload.simpleFile(file, response);
            }
        }
    });

    // 查询filelist是否存在 增加鼠标事件
    filelistMouseEvent($(".filelist").find( 'li' ), $(".file-panel"));
});