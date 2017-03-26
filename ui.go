package main

const html = `
<!DOCTYPE html>
<!--suppress ALL -->
<html>
<head>
    <meta charset="utf-8">
    <meta http-equiv="X-UA-Compatible" content="IE=edge">
    <meta name="viewport" content="width=device-width, initial-scale=1">
    <title>录直播 - WebUI</title>
    <link href="http://cdn.bootcss.com/bootstrap/3.3.7/css/bootstrap.min.css" rel="stylesheet">
    <script src="http://cdn.bootcss.com/jquery/3.1.1/jquery.min.js"></script>
    <script src="http://cdn.bootcss.com/bootstrap/3.3.7/js/bootstrap.min.js"></script>
    <script src="http://cdn.bootcss.com/flv.js/1.1.0/flv.min.js"></script>
    <script language='javascript'>
        function addTaskAction() {
            $("#processing_ui").modal('show');
            $("#addTaskDialog").modal('hide');
            if ($("#addTask_pathg").attr("hidden") == "hidden") {
                addTaskCheck();
            } else
                doAddTask();
            $("#processing_ui").modal('hide');
            return false;
        }

        function addTasksAction() {
            var urls = $("#addTask_urls").val();
            urls = encodeURIComponent(urls);
            if (urls == "") {
                alert("至少输入一行地址.")
            } else {
                $("#processing_ui").modal('show');
                $("#addTasksDialog").modal('hide');
                var aj = $.ajax({url: "/ajax?act=addex&urls=" + urls, async: false});
                var ret = aj.responseText;
                if (ret > 0) {
                    alert("成功添加" + ret + "条任务.");
                } else {
                    alert("没有成功添加任务.");
                }
                location.reload();
            }
            return false;
        }

        var theTasks = null;
        function showTasks() {
            $("#tasklist").html("");
            $("#processing_ui").modal('show');
            var aj = $.ajax({url: "/ajax?act=tasks", async: false});
            var ret = aj.responseText;
            theTasks = JSON.parse(ret).Tasks;
            var rows = "";
            var row = "";
            for (var i = 0; i < theTasks.length; i++) {
                var v = theTasks[i];
                var inf = v.LiveInfo;
                t = "[无数据]";
                if (inf != null)
                    var t = inf.RoomTitle;
                s = v.Run ? "<div class=\"progress progress-striped active\"><div class=\"progress-bar progress-success\" style=\"width: 100%;\"></div></div>" : "<div class=\"progress progress-striped\"><div class=\"progress-bar progress-bar-danger\" role=\"progressbar\" style=\"width: 100%;\"></div></div>";
                var m = v.M ? "循环" : "普通";
                var l = !v.Run ? "未运行" : v.TimeLong;
                var ls = v.Live ? "<span class=\"glyphicon glyphicon-ok-circle\" />" : "<span class=\"glyphicon glyphicon-remove-circle\" />";
                var buttons = "<th>";
                if (v.Run)
                    buttons += "<button onclick=\"stopBtnEvt(" + (i + 1) + ")\" class=\"btn btn-warning\" type=\"button\"><span class=\"glyphicon glyphicon-stop\" /> 停止</button>\n";
                else {
                    buttons += "<button onclick=\"startBtnEvt(" + (i + 1) + ")\" class=\"btn btn-success\" type=\"button\"><span class=\"glyphicon glyphicon-play\" /> 开始</button>\n";
                    buttons += "<button onclick=\"delBtnEvt(" + (i + 1) + ")\" class=\"btn btn-danger\" type=\"button\"><span class=\"glyphicon glyphicon-remove\" /> 删除</button>\n"
                }
                if (v.Files != null) {
                    if (v.Files.length == 1) {
                        buttons += "<button onclick=\"down(" + (i + 1) + ",0)\" class=\"btn btn-primary\" type=\"button\"><span class=\"glyphicon glyphicon-download\" /> 下载</button>\n";
                        buttons += "<button onclick=\"play(" + (i + 1) + ",0)\" class=\"btn btn-primary\" type=\"button\"><span class=\"glyphicon glyphicon-play-circle\" />  播放</button>\n"
                    } else {
                        var dlBtn = "<div class=\"btn-group\"><button type=\"button\" class=\"btn btn-primary dropdown-toggle\" data-toggle=\"dropdown\"><span class=\"glyphicon glyphicon-download\" /> 下载 <span class=\"caret\"></span></button><ul class=\"dropdown-menu\" role=\"menu\">";
                        var pyBtn = "<div class=\"btn-group\"><button type=\"button\" class=\"btn btn-primary dropdown-toggle\" data-toggle=\"dropdown\"><span class=\"glyphicon glyphicon-play-circle\" /> 播放 <span class=\"caret\"></span></button><ul class=\"dropdown-menu\" role=\"menu\">";
                        for (var x = 0; x < v.Files.length; x++) {
                            dlBtn += "<li onclick=\"down(" + (i + 1) + "," + (x + 1) + ")\"><button class=\"btn btn-link\"><span class=\"glyphicon glyphicon-download-alt\" /> 分段" + (x + 1) + "</a></li>";
                            pyBtn += "<li onclick=\"play(" + (i + 1) + "," + (x + 1) + ")\"><button class=\"btn btn-link\"><span class=\"glyphicon glyphicon-play\" /> 分段" + (x + 1) + "</a></li>"
                        }
                        dlBtn += "</ul></div>\n";
                        pyBtn += "</ul></div>\n";
                        buttons += dlBtn;
                        buttons += pyBtn;
                    }
                }
                buttons += "<button onclick=\"infoBtnEvt(" + (i + 1) + ")\" class=\"btn btn-info\" type=\"button\"><span class=\"glyphicon glyphicon-option-horizontal\" /> 详情</button>\n";
                buttons += "</th>";
                row = "<tr>";
                row += "<td>" + (i + 1) + "</td>";
                row += "<td><a target=\"_blank\" href=\"" + v.SiteURL + "\"  title=\"" + v.Site + "\"><img height=\"16\" width=\"16\" src=" + v.SiteIcon + " /></a></td>";
                row += "<td>" + ls + "</td>";
                row += "<td>" + m + "</td>";
                row += "<td>" + s + "</td>";
                row += "<td>" + l + "</td>";
                row += "<td>" + t + "</td>";
                row += buttons;
                row += "</tr>";
                rows += row;
            }
            $("#tasklist").html(rows);
            $("#processing_ui").modal('hide');
        }

        $(document).ready(function () {
            var aj = $.ajax({url: "/ajax?act=ver", async: false});
            var ret = aj.responseText;
            var arr = ret.split("|");
            $("#tver").text("(Ver " + arr[0] + ")");
            if (arr[1] != "null") {
                $("#uver").text("下载最新版本预编译包(Ver " + arr[1] + ")");
                $("#uver").removeAttr("hidden");
            }
            showTasks();
        });

        function doAddTask() {
            var url = $("#addTask_url").val();
            var path = $("#addTask_path").val();
            if (checkPathExist(path)) {
                alert("文件(路径)已存在,请更换.");
                $("#addTaskDialog").modal('show');
                return;
            }
            var m = $("#addTask_m").is(':checked');
            var r = $("#addTask_run").is(':checked');
            if (!m)
                path += ".flv";
            var aj = $.ajax({
                url: "/ajax?act=add&url=" + url + "&path=" + path + "&m=" + m + "&run=" + r,
                async: false
            });
            var ret = aj.responseText;
            if (ret != "ok")
                alert("添加任务失败.");
            location.reload();
        }

        function checkPathExist(path) {
            var aj = $.ajax({url: "/ajax?act=exist&path=" + path, async: false});
            var ret = aj.responseText;
            return ret == "exist";
        }

        function addTaskCheck() {
            var url = $("#addTask_url").val();
            var aj = $.ajax({url: "/ajax?act=check&url=" + url, async: false});
            var ret = aj.responseText;
            var j = JSON.parse(ret);
            if (!j.Pass)
                alert("不支持的地址.");
            else if (j.Has) {
                if (!j.Live) {
                    $("#addTask_m").attr("checked", "checked");
                    $("#addTask_m").attr("disabled", "disabled");
                    $("#addTask_mg").attr("class", "checkbox disabled");
                }
                $("#addTask_path").val(j.Path);
                $("#addTask_url").attr("readonly", 'readonly');
                $("#addTask_pathg").removeAttr("hidden");
                $("#addTaskDialog").modal('show');
                return;
            } else
                alert("不存在的房间.");
            $("#addTaskDialog").modal('show');
        }

        function startBtnEvt(o) {
            if (theTasks[o - 1].Files != null && !confirm("文件(路径)已存在,是否覆盖并继续?"))
                return;
            var aj = $.ajax({url: "/ajax?act=start&id=" + o, async: false});
            var ret = aj.responseText;
            if (ret != "ok")
                alert("开始任务失败.");
            else
                showTasks();
        }

        function stopBtnEvt(o) {
            if (confirm("确定要停止此任务?")) {
                var aj = $.ajax({url: "/ajax?act=stop&id=" + o, async: false});
                var ret = aj.responseText;
                if (ret != "ok")
                    alert("停止任务失败.");
                else
                    showTasks();
            }
        }

        function delBtnEvt(o) {
            if (confirm("确定要删除此任务?")) {
                var f = confirm("删除文件(路径)?");
                var aj = $.ajax({url: "/ajax?act=del&id=" + o + "&f=" + f, async: false});
                var ret = aj.responseText;
                if (ret != "ok")
                    alert("删除任务失败.");
                else
                    showTasks();
            }
        }

        function infoBtnEvt(o) {
            var v = theTasks[o - 1];
            var i = v.LiveInfo;
            $("#info_url").val(v.URL);
            $("#info_start").val(v.Run ? v.StartTime : "未开始");
            $("#info_index").val(v.Index);
            $("#info_path").val(v.Path);
            $("#info_live").attr('hidden', "hidden");
            if (v.Live) {
                $("#info_live").removeAttr('hidden');
                $("#info_nick").val(i.LiveNick);
                $("#info_d").val(i.RoomDetails);
                $("#info_i").attr("src", i.LivingIMG);
            }
            $("#info_ui").modal('show').on("");
        }

        function down(o, s) {
            var u = "/ajax?act=get&id=" + o;
            if (s != 0)
                u += "&sub=" + s;
            window.location.href = u;
        }

        function play(o, s) {
            var u = "/ajax?act=get&id=" + o;
            if (s != 0)
                u += "&sub=" + s;
            var flvPlayer = flvjs.createPlayer({
                type: "flv",
                url: u
            });
            $("#player_ui").modal('show').on("hide.bs.modal", function () {
                flvPlayer.unload();
            });
            if (flvjs.isSupported()) {
                var videoElement = document.getElementById("videoElement");
                flvPlayer.attachMediaElement(videoElement);
                flvPlayer.load();
                flvPlayer.play();
            }
        }
    </script>
</head>

<body>
<div class="container-fluid ">
    <div class="row-fluid ">
        <div class="span12 ">
            <div class="page-header ">
                <h1><span class="glyphicon glyphicon-facetime-video"></span> 录直播
                    <small id="tver"></small>
                </h1>
            </div>
            <div style="text-align:center">
                <h3>微信打赏</h3>
                <img src="data:image/jpeg;base64,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"/>
                <h4>反馈,建议QQ群:<a target="_noblank"
                                href="https://shang.qq.com/wpa/qunwpa?idkey=c5cb06c2e524eb7872f80317621bf51f9c0b1d3451451232dbe6aa3fbdbb15ff">23153821</a>
                </h4>
            </div>
            <div class="btn-group">
                <button class="btn btn-primary" data-toggle="modal" data-target="#addTaskDialog"><span
                        class="glyphicon glyphicon-plus"></span> 添加任务...
                </button>
                <button data-toggle="dropdown" class="btn btn-primary dropdown-toggle"><span class="caret"></span>
                </button>
                <ul class="dropdown-menu">
                    <li>
                        <button class="btn btn-link" data-toggle="modal" data-target="#addTasksDialog"><span
                                class="glyphicon glyphicon-plus"></span> 批量添加...
                        </button>
                    </li>
                </ul>
            </div>
            <button class="btn btn-default" onclick="showTasks()"><span
                    class="glyphicon glyphicon-refresh"></span> 刷新列表
            </button>
            <button class="btn btn-link" id="uver" type="button" hidden="hidden"
                    onclick="window.open('https://github.com/Baozisoftware/Luzhibo-go/releases/download/latest/releases.7z')"></button>
            <table class="table ">
                <thead>
                <tr>
                    <th><span class="glyphicon glyphicon-minus"></span> 任务编号</th>
                    <th><span class="glyphicon glyphicon-hdd"></span> 直播平台</th>
                    <th><span class="glyphicon glyphicon-play-circle"></span> 开播状态</th>
                    <th><span class="glyphicon glyphicon-flag"></span> 任务模式</th>
                    <th><span class="glyphicon glyphicon-stats"></span> 运行状态</th>
                    <th><span class="glyphicon glyphicon-time"></span> 运行时长</th>
                    <th><span class="glyphicon glyphicon-book"></span> 房间标题</th>
                    <th><span class="glyphicon glyphicon-hand-right"></span> 其他操作</th>
                </tr>
                </thead>
                <tbody id="tasklist">
                </tbody>
            </table>
        </div>
    </div>
</div>

<!--添加任务遮罩层-->
<div class="modal fade" id="addTaskDialog" tabindex="-1" role="dialog" data-backdrop="static" data-keyboard="false">
    <div class="modal-dialog">
        <div class="modal-content">
            <form onsubmit="return addTaskAction()">
                <div class="modal-header">
                    <h4 class="modal-title"><span class="glyphicon glyphicon-plus"></span> 添加任务</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label><span class="glyphicon glyphicon-film"></span> 直播地址:</label>
                        <label for="addTask_url"></label><input type="url" class="form-control" id="addTask_url"
                                                                required="required"/>
                    </div>
                    <div id="addTask_pathg" hidden="hidden">
                        <div class="form-group">
                            <label><span class="glyphicon glyphicon-folder-open"></span>
                                保存路径(文件名,自动添加后缀”.flv“):</label>
                            <label for="addTask_path"></label><input type="text" class="form-control"
                                                                     id="addTask_path"
                                                                     required="required" value="#"/>
                        </div>
                        <div class="form-group">
                            <div class="checkbox" id="addTask_mg">
                                <label><input type="checkbox" id="addTask_m"/><span
                                        class="glyphicon glyphicon-flash"></span> 循环模式</label>
                            </div>
                            <div class="checkbox">
                                <label><input type="checkbox" id="addTask_run" checked="checked"/><span
                                        class="glyphicon glyphicon-play"></span> 立即开始</label>
                            </div>
                        </div>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-danger" onclick="location.reload()"><span
                            class="glyphicon glyphicon-remove"></span> 关闭
                    </button>
                    <button type="submit" class="btn btn-primary"><span class="glyphicon glyphicon-ok"></span> 添加
                    </button>
                </div>
            </form>
        </div>
    </div>
</div>
<!--进度条遮罩层-->
<div class="modal fade" id="processing_ui" tabindex="-1" role="dialog" data-backdrop="static" data-keyboard="false">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <h4 class="modal-title"><span class="glyphicon glyphicon-tasks"></span> 正在处理中...</h4>
            </div>
            <div class="modal-body">
                <div class="progress progress-striped active">
                    <div class="progress-bar progress-success" style="width: 100%;"></div>
                </div>
            </div>
        </div>
    </div>
</div>
<!--详情遮罩层-->
<div class="modal fade" id="info_ui" tabindex="-1" role="dialog">
    <div class="modal-dialog">
        <div class="modal-content">
            <div class="modal-header">
                <h4 class="modal-title"><span class="glyphicon glyphicon-tasks"></span> 任务详情</h4>
            </div>
            <div class="modal-body">
                <div class="form-group">
                    <label><span class="glyphicon glyphicon-film"></span> 直播地址:</label>
                    <label for="info_url"></label><input class="form-control" id="info_url" readonly="readonly"/>
                </div>
                <div class="form-group">
                    <label><span class="glyphicon glyphicon-time"></span> 开始时间:</label>
                    <label for="info_start"></label><input type="text" class="form-control" id="info_start"
                                                           readonly="readonly"/>
                </div>
                <div class="form-group">
                    <label><span class="glyphicon glyphicon-file"></span> 已有文件:</label>
                    <label for="info_index"></label><input type="text" class="form-control" id="info_index"
                                                           readonly="readonly"/>
                </div>
                <div class="form-group">
                    <label><span class="glyphicon glyphicon-folder-open"></span> 保存路径:</label>
                    <label for="info_path"></label><input type="text" class="form-control" id="info_path"
                                                          readonly="readonly"/>
                </div>
                <br/>
                <div id="info_live" hidden="hidden">
                    <div class="form-group">
                        <label><span class="glyphicon glyphicon-user"></span> 主播昵称:</label>
                        <label for="info_nick"></label><input type="text" class="form-control" id="info_nick"
                                                              readonly="readonly"/>
                    </div>
                    <div class="form-group">
                        <label><span class="glyphicon glyphicon-picture"></span> 直播截图:</label>
                        <br/>
                        <img id="info_i" width="320" height="180"/>
                    </div>
                    <div class="form-group">
                        <label><span class="glyphicon glyphicon-pencil"></span> 房间说明:</label>
                        <textarea class="form-control" id="info_d" readonly="readonly" rows="5" title=""></textarea>
                    </div>
                </div>
            </div>
        </div>
    </div>
</div>
<!--批量添加遮罩层-->
<div class="modal fade" id="addTasksDialog" tabindex="-1" role="dialog" data-backdrop="static" data-keyboard="false">
    <div class="modal-dialog">
        <div class="modal-content">
            <form onsubmit="return addTasksAction()">
                <div class="modal-header">
                    <h4 class="modal-title"><span class="glyphicon glyphicon-plus"></span> 批量添加</h4>
                </div>
                <div class="modal-body">
                    <div class="form-group">
                        <label><span class="glyphicon glyphicon-film"></span> 直播地址(一行一个):</label>
                        <label for="addTask_url"></label><textarea class="form-control" id="addTask_urls"
                                                                   rows="5"></textarea>
                    </div>
                </div>
                <div class="modal-footer">
                    <button type="button" class="btn btn-danger" onclick="location.reload()"><span
                            class="glyphicon glyphicon-remove"></span> 关闭
                    </button>
                    <button type="submit" class="btn btn-primary"><span class="glyphicon glyphicon-ok"></span> 添加
                    </button>
                </div>
            </form>
        </div>
    </div>
</div>
</body>
</html>
`
